package mls

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

const (
	version10               uint16 = 1
	cipherSuiteID           uint16 = 2
	leafNodeSourceKeyPackage uint8  = 1
	credentialTypeBasic      uint8  = 1
)

type KeyPackageBundle struct {
	Serialized     []byte
	initPriv       []byte
	sigKey         *signatureKeyPair
	encryptionPriv []byte
}

type WelcomeResult struct {
	ExporterSecret []byte
	Epoch          uint64
	GroupID        []byte
}

func GenerateKeyPackage(identity []byte) (*KeyPackageBundle, error) {
	initKP, err := generateHPKEKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generating init key: %w", err)
	}
	encKP, err := generateHPKEKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generating encryption key: %w", err)
	}
	sigKP, err := generateSignatureKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generating signature key: %w", err)
	}

	leafContent := buildLeafNodeContent(encKP.pub, sigKP.pub, identity)

	leafSig, err := signWithLabel(sigKP, "LeafNodeTBS", leafContent)
	if err != nil {
		return nil, fmt.Errorf("signing leaf node: %w", err)
	}

	leafNode := &tlsWriter{}
	leafNode.writeRaw(leafContent)
	leafNode.writeVec(leafSig)

	kpContent := buildKeyPackageContent(initKP.pub, leafNode.bytes())

	kpSig, err := signWithLabel(sigKP, "KeyPackageTBS", kpContent)
	if err != nil {
		return nil, fmt.Errorf("signing key package: %w", err)
	}

	kpFull := &tlsWriter{}
	kpFull.writeRaw(kpContent)
	kpFull.writeVec(kpSig)

	return &KeyPackageBundle{
		Serialized:     kpFull.bytes(),
		initPriv:       initKP.priv,
		sigKey:         sigKP,
		encryptionPriv: encKP.priv,
	}, nil
}

func ProcessWelcome(data []byte, bundle *KeyPackageBundle) (*WelcomeResult, error) {
	r := &tlsReader{data: data}

	cs := r.readUint16()
	if r.err != nil {
		return nil, fmt.Errorf("reading cipher suite: %w", r.err)
	}
	if cs != cipherSuiteID {
		return nil, fmt.Errorf("unexpected cipher suite: %d", cs)
	}

	secretsData := r.readVec()
	if r.err != nil {
		return nil, fmt.Errorf("reading secrets: %w", r.err)
	}

	encryptedGroupInfo := r.readVec()
	if r.err != nil {
		return nil, fmt.Errorf("reading encrypted group info: %w", r.err)
	}

	ourRef := keyPackageRef(bundle.Serialized)

	sr := &tlsReader{data: secretsData}
	var kemOutput, encryptedSecrets []byte
	found := false

	for sr.remaining() > 0 && sr.err == nil {
		newMember := sr.readVec()
		kemOut := sr.readVec()
		ct := sr.readVec()
		if sr.err != nil {
			break
		}

		if bytesEqual(newMember, ourRef) {
			kemOutput = kemOut
			encryptedSecrets = ct
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no matching EncryptedGroupSecrets for our KeyPackageRef")
	}

	infoW := &tlsWriter{}
	infoW.writeVec([]byte("MLS 1.0 Welcome"))
	infoW.writeVec(encryptedGroupInfo)

	groupSecretsPlain, err := hpkeDecrypt(bundle.initPriv, kemOutput, infoW.bytes(), nil, encryptedSecrets)
	if err != nil {
		return nil, fmt.Errorf("HPKE decrypting group secrets: %w", err)
	}

	gsr := &tlsReader{data: groupSecretsPlain}
	joinerSecret := gsr.readVec()
	if gsr.err != nil {
		return nil, fmt.Errorf("reading joiner secret: %w", gsr.err)
	}
	hasPathSecret := gsr.readUint8()
	if hasPathSecret == 1 {
		_ = gsr.readVec()
	}

	pskSecret := make([]byte, 32)
	memberSecret := hkdfExtract(joinerSecret, pskSecret)

	welcomeSecret, err := ExpandWithLabel(memberSecret, "welcome", nil, 32)
	if err != nil {
		return nil, fmt.Errorf("deriving welcome secret: %w", err)
	}

	welcomeKey, err := ExpandWithLabel(welcomeSecret, "key", nil, 16)
	if err != nil {
		return nil, fmt.Errorf("deriving welcome key: %w", err)
	}

	welcomeNonce, err := ExpandWithLabel(welcomeSecret, "nonce", nil, 12)
	if err != nil {
		return nil, fmt.Errorf("deriving welcome nonce: %w", err)
	}

	block, err := aes.NewCipher(welcomeKey)
	if err != nil {
		return nil, fmt.Errorf("creating AES cipher for GroupInfo: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM for GroupInfo: %w", err)
	}
	groupInfoPlain, err := gcm.Open(nil, welcomeNonce, encryptedGroupInfo, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypting GroupInfo: %w", err)
	}

	groupContext, epoch, groupID, err := parseGroupInfoForContext(groupInfoPlain)
	if err != nil {
		return nil, fmt.Errorf("parsing GroupInfo: %w", err)
	}

	epochSecret, err := ExpandWithLabel(memberSecret, "epoch", groupContext, 32)
	if err != nil {
		return nil, fmt.Errorf("deriving epoch secret: %w", err)
	}

	exporterSecret, err := deriveSecret(epochSecret, "exporter")
	if err != nil {
		return nil, fmt.Errorf("deriving exporter secret: %w", err)
	}

	return &WelcomeResult{
		ExporterSecret: exporterSecret,
		Epoch:          epoch,
		GroupID:        groupID,
	}, nil
}

func Export(exporterSecret []byte, label string, context []byte, length int) ([]byte, error) {
	derivedSecret, err := deriveSecret(exporterSecret, label)
	if err != nil {
		return nil, err
	}
	contextHash := sha256.Sum256(context)
	return ExpandWithLabel(derivedSecret, "exported", contextHash[:], length)
}

func ExpandWithLabel(secret []byte, label string, context []byte, length int) ([]byte, error) {
	mlsLabel := []byte("MLS 1.0 " + label)

	w := &tlsWriter{}
	w.writeUint16(uint16(length))
	w.writeVec(mlsLabel)
	w.writeVec(context)

	r := hkdf.Expand(sha256.New, secret, w.bytes())
	out := make([]byte, length)
	_, err := io.ReadFull(r, out)
	return out, err
}
