package mls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/cloudflare/circl/hpke"
	"golang.org/x/crypto/hkdf"
)

type tlsWriter struct {
	buf []byte
}

func (w *tlsWriter) writeUint8(v uint8) {
	w.buf = append(w.buf, v)
}

func (w *tlsWriter) writeUint16(v uint16) {
	w.buf = append(w.buf, byte(v>>8), byte(v))
}

func (w *tlsWriter) writeUint32(v uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	w.buf = append(w.buf, b...)
}

func (w *tlsWriter) writeUint64(v uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	w.buf = append(w.buf, b...)
}

func (w *tlsWriter) writeVec(data []byte) {
	w.writeVarint(uint64(len(data)))
	w.buf = append(w.buf, data...)
}

func (w *tlsWriter) writeVarint(v uint64) {
	switch {
	case v <= 63:
		w.buf = append(w.buf, byte(v))
	case v <= 16383:
		w.buf = append(w.buf, byte(0x40|v>>8), byte(v))
	case v <= 1073741823:
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(v)|0x80000000)
		w.buf = append(w.buf, b...)
	default:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, v|0xC000000000000000)
		w.buf = append(w.buf, b...)
	}
}

func (w *tlsWriter) writeRaw(data []byte) {
	w.buf = append(w.buf, data...)
}

func (w *tlsWriter) bytes() []byte {
	return w.buf
}

type tlsReader struct {
	data []byte
	pos  int
	err  error
}

func (r *tlsReader) remaining() int {
	return len(r.data) - r.pos
}

func (r *tlsReader) readUint8() uint8 {
	if r.err != nil || r.pos+1 > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read uint8 at pos %d", r.pos)
		return 0
	}
	v := r.data[r.pos]
	r.pos++
	return v
}

func (r *tlsReader) readUint16() uint16 {
	if r.err != nil || r.pos+2 > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read uint16 at pos %d", r.pos)
		return 0
	}
	v := binary.BigEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return v
}

func (r *tlsReader) readUint32() uint32 {
	if r.err != nil || r.pos+4 > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read uint32 at pos %d", r.pos)
		return 0
	}
	v := binary.BigEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return v
}

func (r *tlsReader) readUint64() uint64 {
	if r.err != nil || r.pos+8 > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read uint64 at pos %d", r.pos)
		return 0
	}
	v := binary.BigEndian.Uint64(r.data[r.pos:])
	r.pos += 8
	return v
}

func (r *tlsReader) readVec() []byte {
	if r.err != nil {
		return nil
	}
	length := r.readVarint()
	if r.err != nil {
		return nil
	}
	n := int(length)
	if r.pos+n > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read vec len=%d at pos %d", n, r.pos)
		return nil
	}
	out := make([]byte, n)
	copy(out, r.data[r.pos:r.pos+n])
	r.pos += n
	return out
}

func (r *tlsReader) readVarint() uint64 {
	if r.err != nil || r.pos >= len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read varint at pos %d", r.pos)
		return 0
	}
	first := r.data[r.pos]
	kind := first >> 6
	switch kind {
	case 0:
		r.pos++
		return uint64(first & 0x3F)
	case 1:
		if r.pos+2 > len(r.data) {
			r.err = fmt.Errorf("tlsReader: short read varint(2) at pos %d", r.pos)
			return 0
		}
		v := uint64(first&0x3F)<<8 | uint64(r.data[r.pos+1])
		r.pos += 2
		return v
	case 2:
		if r.pos+4 > len(r.data) {
			r.err = fmt.Errorf("tlsReader: short read varint(4) at pos %d", r.pos)
			return 0
		}
		v := uint64(first&0x3F)<<24 | uint64(r.data[r.pos+1])<<16 | uint64(r.data[r.pos+2])<<8 | uint64(r.data[r.pos+3])
		r.pos += 4
		return v
	default:
		if r.pos+8 > len(r.data) {
			r.err = fmt.Errorf("tlsReader: short read varint(8) at pos %d", r.pos)
			return 0
		}
		v := uint64(first&0x3F)<<56 | uint64(r.data[r.pos+1])<<48 | uint64(r.data[r.pos+2])<<40 | uint64(r.data[r.pos+3])<<32 |
			uint64(r.data[r.pos+4])<<24 | uint64(r.data[r.pos+5])<<16 | uint64(r.data[r.pos+6])<<8 | uint64(r.data[r.pos+7])
		r.pos += 8
		return v
	}
}

func (r *tlsReader) readRaw(n int) []byte {
	if r.err != nil || r.pos+n > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short read raw(%d) at pos %d", n, r.pos)
		return nil
	}
	out := make([]byte, n)
	copy(out, r.data[r.pos:r.pos+n])
	r.pos += n
	return out
}

func (r *tlsReader) readRemaining() []byte {
	if r.err != nil {
		return nil
	}
	out := make([]byte, len(r.data)-r.pos)
	copy(out, r.data[r.pos:])
	r.pos = len(r.data)
	return out
}

func (r *tlsReader) skip(n int) {
	if r.err != nil || r.pos+n > len(r.data) {
		r.err = fmt.Errorf("tlsReader: short skip(%d) at pos %d", n, r.pos)
		return
	}
	r.pos += n
}

type hpkeKeyPair struct {
	pub  []byte
	priv []byte
}

func generateHPKEKeyPair() (*hpkeKeyPair, error) {
	kemID := hpke.KEM_P256_HKDF_SHA256
	scheme := kemID.Scheme()
	pub, priv, err := scheme.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generating HPKE key pair: %w", err)
	}
	pubBytes, err := pub.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshaling HPKE public key: %w", err)
	}
	privBytes, err := priv.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshaling HPKE private key: %w", err)
	}
	return &hpkeKeyPair{pub: pubBytes, priv: privBytes}, nil
}

func hpkeDecrypt(privKeyBytes, kemOutput, info, aad, ciphertext []byte) ([]byte, error) {
	kemID := hpke.KEM_P256_HKDF_SHA256
	scheme := kemID.Scheme()

	priv, err := scheme.UnmarshalBinaryPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling HPKE private key: %w", err)
	}

	suite := hpke.NewSuite(hpke.KEM_P256_HKDF_SHA256, hpke.KDF_HKDF_SHA256, hpke.AEAD_AES128GCM)
	opener, err := suite.NewReceiver(priv, info)
	if err != nil {
		return nil, fmt.Errorf("creating HPKE receiver: %w", err)
	}
	ctx, err := opener.Setup(kemOutput)
	if err != nil {
		return nil, fmt.Errorf("HPKE setup: %w", err)
	}
	plaintext, err := ctx.Open(ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("HPKE open: %w", err)
	}
	return plaintext, nil
}

type signatureKeyPair struct {
	pub  []byte
	priv *ecdsa.PrivateKey
}

func generateSignatureKeyPair() (*signatureKeyPair, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key: %w", err)
	}
	pub := elliptic.Marshal(elliptic.P256(), priv.PublicKey.X, priv.PublicKey.Y)
	return &signatureKeyPair{pub: pub, priv: priv}, nil
}

func signWithLabel(key *signatureKeyPair, label string, content []byte) ([]byte, error) {
	mlsLabel := []byte("MLS 1.0 " + label)

	w := &tlsWriter{}
	w.writeVec(mlsLabel)
	w.writeVec(content)

	hash := sha256.Sum256(w.bytes())
	r, s, err := ecdsa.Sign(rand.Reader, key.priv, hash[:])
	if err != nil {
		return nil, err
	}
	return marshalECDSASignature(r, s), nil
}

func marshalECDSASignature(r, s *big.Int) []byte {
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	if len(rBytes) > 0 && rBytes[0]&0x80 != 0 {
		rBytes = append([]byte{0}, rBytes...)
	}
	if len(sBytes) > 0 && sBytes[0]&0x80 != 0 {
		sBytes = append([]byte{0}, sBytes...)
	}
	totalLen := 2 + len(rBytes) + 2 + len(sBytes)
	sig := make([]byte, 0, 2+totalLen)
	sig = append(sig, 0x30, byte(totalLen))
	sig = append(sig, 0x02, byte(len(rBytes)))
	sig = append(sig, rBytes...)
	sig = append(sig, 0x02, byte(len(sBytes)))
	sig = append(sig, sBytes...)
	return sig
}

func deriveSecret(secret []byte, label string) ([]byte, error) {
	return ExpandWithLabel(secret, label, []byte{}, 32)
}

func buildLeafNodeContent(encryptionKey, signatureKey, identity []byte) []byte {
	w := &tlsWriter{}

	w.writeVec(encryptionKey)
	w.writeVec(signatureKey)

	w.writeUint16(uint16(credentialTypeBasic))
	w.writeVec(identity)

	versionsData := make([]byte, 2)
	binary.BigEndian.PutUint16(versionsData, version10)
	w.writeVec(versionsData)

	csData := make([]byte, 2)
	binary.BigEndian.PutUint16(csData, cipherSuiteID)
	w.writeVec(csData)

	w.writeVec(nil)
	w.writeVec(nil)

	credData := make([]byte, 2)
	binary.BigEndian.PutUint16(credData, uint16(credentialTypeBasic))
	w.writeVec(credData)

	w.writeUint8(leafNodeSourceKeyPackage)

	w.writeUint64(0)
	w.writeUint64(^uint64(0))

	w.writeVec(nil)

	return w.bytes()
}

func buildKeyPackageContent(initKey []byte, leafNode []byte) []byte {
	w := &tlsWriter{}
	w.writeUint16(version10)
	w.writeUint16(cipherSuiteID)
	w.writeVec(initKey)
	w.writeRaw(leafNode)
	w.writeVec(nil)
	return w.bytes()
}

func keyPackageRef(serializedKeyPackage []byte) []byte {
	return refHash("MLS 1.0 KeyPackage Reference", serializedKeyPackage)
}

func refHash(label string, value []byte) []byte {
	w := &tlsWriter{}
	w.writeVec([]byte(label))
	w.writeVec(value)
	h := sha256.Sum256(w.bytes())
	return h[:]
}

func hkdfExtract(salt, ikm []byte) []byte {
	if salt == nil {
		salt = make([]byte, 32)
	}
	return hkdf.Extract(sha256.New, ikm, salt)
}

func parseGroupInfoForContext(data []byte) (groupContext []byte, epoch uint64, groupID []byte, err error) {
	r := &tlsReader{data: data}

	ctxStart := r.pos

	_ = r.readUint16()
	_ = r.readUint16()
	groupID = r.readVec()
	epoch = r.readUint64()
	_ = r.readVec()
	_ = r.readVec()
	_ = r.readVec()

	if r.err != nil {
		return nil, 0, nil, fmt.Errorf("parsing GroupContext: %w", r.err)
	}

	ctxEnd := r.pos
	groupContext = make([]byte, ctxEnd-ctxStart)
	copy(groupContext, data[ctxStart:ctxEnd])

	return groupContext, epoch, groupID, nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
