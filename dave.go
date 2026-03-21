package discordgo

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo/mls"
)

var opusSilencePacket = [3]byte{0xF8, 0xFF, 0xFE}

type daveReceiver struct {
	userID            string
	baseSecret        []byte
	currentGeneration uint32
	key               []byte
	aesBlock          cipher.Block
	frameCipher       cipher.AEAD
}

type DAVESession struct {
	mu                  sync.Mutex
	protocolVersion     int
	epoch               uint64
	pendingTransitionID uint16
	pendingVersion      int

	exporterSecret    []byte
	senderKey         []byte
	senderNonce       uint32
	frameCipher       cipher.AEAD
	userID            string
	active            bool
	ratchetBaseSecret []byte
	currentGeneration uint32
	hasPendingKey     bool

	ssrcToUserID map[uint32]string
	receivers    map[uint32]*daveReceiver

	kpBundle *mls.KeyPackageBundle
}

func NewDAVESession(userID string) *DAVESession {
	return &DAVESession{
		userID: userID,
	}
}

func (d *DAVESession) GenerateKeyPackage() ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.generateKeyPackageLocked()
}

func (d *DAVESession) ResetForReWelcome() ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.exporterSecret = nil
	d.hasPendingKey = false

	return d.generateKeyPackageLocked()
}

func (d *DAVESession) generateKeyPackageLocked() ([]byte, error) {
	userIDNum, err := strconv.ParseUint(d.userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing user ID for credential: %w", err)
	}
	identity := make([]byte, 8)
	binary.BigEndian.PutUint64(identity, userIDNum)

	bundle, err := mls.GenerateKeyPackage(identity)
	if err != nil {
		return nil, fmt.Errorf("generating key package: %w", err)
	}
	d.kpBundle = bundle
	return bundle.Serialized, nil
}

func (d *DAVESession) HandleExternalSenderPackage(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return nil
}

func (d *DAVESession) HandleWelcome(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.kpBundle == nil {
		return fmt.Errorf("no key package generated")
	}

	result, err := mls.ProcessWelcome(data, d.kpBundle)
	if err != nil {
		return fmt.Errorf("processing welcome: %w", err)
	}

	d.exporterSecret = result.ExporterSecret
	d.epoch = result.Epoch
	d.hasPendingKey = true
	d.receivers = nil
	return nil
}

func (d *DAVESession) HandleCommit(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return nil
}

func (d *DAVESession) HandlePrepareTransition(transitionID uint16, protocolVersion int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pendingTransitionID = transitionID
	d.pendingVersion = protocolVersion
}

func (d *DAVESession) HandleExecuteTransition(transitionID uint16) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if transitionID != d.pendingTransitionID {
		if d.senderKey != nil {
			d.active = true
		}
		return nil
	}

	if d.pendingVersion > 0 {
		derivedNewKey := false
		if d.hasPendingKey && d.exporterSecret != nil {
			if err := d.deriveSenderKeyLocked(); err != nil {
				return err
			}
			d.hasPendingKey = false
			derivedNewKey = true
		}
		if d.senderKey == nil {
			return nil
		}

		if !derivedNewKey && !d.hasPendingKey {
			d.active = false
			d.senderKey = nil
			d.frameCipher = nil
			d.ratchetBaseSecret = nil
			d.currentGeneration = 0
			return nil
		}

		d.active = true
	} else {
		d.active = false
		d.senderKey = nil
		d.frameCipher = nil
		d.hasPendingKey = false
	}
	return nil
}

func (d *DAVESession) HandlePrepareEpoch(epoch uint64, protocolVersion int) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.epoch = epoch
	d.active = false
	d.senderKey = nil
	d.frameCipher = nil
	d.exporterSecret = nil
	d.receivers = nil

	return d.generateKeyPackageLocked()
}

func (d *DAVESession) DeriveSenderKey() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.deriveSenderKeyLocked()
}

func (d *DAVESession) deriveSenderKeyLocked() error {
	if d.exporterSecret == nil {
		return fmt.Errorf("no exporter secret")
	}

	userIDNum, err := strconv.ParseUint(d.userID, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing user ID: %w", err)
	}
	context := make([]byte, 8)
	binary.LittleEndian.PutUint64(context, userIDNum)

	baseSecret, err := mls.Export(d.exporterSecret, daveExportLabel, context, daveKeySize)
	if err != nil {
		return fmt.Errorf("exporting base secret: %w", err)
	}

	d.ratchetBaseSecret = baseSecret
	d.currentGeneration = 0
	d.senderNonce = 0

	key, err := hashRatchetGetKey(baseSecret, 0)
	if err != nil {
		return fmt.Errorf("deriving ratchet key: %w", err)
	}
	d.senderKey = key

	frameCipher, err := newDAVECipher(key)
	if err != nil {
		return fmt.Errorf("creating frame cipher: %w", err)
	}
	d.frameCipher = frameCipher
	return nil
}

func (d *DAVESession) EncryptFrame(opusData []byte) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.frameCipher == nil {
		return nil, fmt.Errorf("no frame cipher")
	}

	d.senderNonce++

	generation := d.senderNonce >> 24
	if generation != d.currentGeneration {
		d.currentGeneration = generation
		key, err := hashRatchetGetKey(d.ratchetBaseSecret, generation)
		if err != nil {
			return nil, fmt.Errorf("ratcheting key for generation %d: %w", generation, err)
		}
		d.senderKey = key
		frameCipher, err := newDAVECipher(key)
		if err != nil {
			return nil, fmt.Errorf("creating cipher for generation %d: %w", generation, err)
		}
		d.frameCipher = frameCipher
	}

	encrypted := encryptSecureFrame(d.frameCipher, d.senderNonce, opusData)
	return encrypted, nil
}

func (d *DAVESession) SetSSRC(ssrc uint32, userID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ssrcToUserID == nil {
		d.ssrcToUserID = make(map[uint32]string)
	}
	d.ssrcToUserID[ssrc] = userID
}

func (d *DAVESession) DecryptFrame(ssrc uint32, data []byte) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(data) == 3 && data[0] == opusSilencePacket[0] && data[1] == opusSilencePacket[1] && data[2] == opusSilencePacket[2] {
		return data, nil
	}

	ciphertext, truncatedTag, nonce, err := parseSecureFrame(data)
	if err == errNotDAVEFrame {
		return data, nil
	}
	if err != nil {
		return nil, err
	}

	recv := d.receivers[ssrc]
	if recv == nil {
		userID, ok := d.ssrcToUserID[ssrc]
		if !ok {
			return nil, fmt.Errorf("unknown SSRC %d", ssrc)
		}
		recv, err = d.createReceiverLocked(ssrc, userID)
		if err != nil {
			return nil, err
		}
	}

	generation := nonce >> 24
	if generation != recv.currentGeneration {
		key, err := hashRatchetGetKey(recv.baseSecret, generation)
		if err != nil {
			return nil, fmt.Errorf("ratcheting receiver key for generation %d: %w", generation, err)
		}
		recv.key = key
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		recv.aesBlock = block
		fc, err := newDAVECipher(key)
		if err != nil {
			return nil, err
		}
		recv.frameCipher = fc
		recv.currentGeneration = generation
	}

	return decryptSecureFrame(recv.aesBlock, recv.frameCipher, nonce, ciphertext, truncatedTag)
}

func (d *DAVESession) createReceiverLocked(ssrc uint32, userID string) (*daveReceiver, error) {
	if d.exporterSecret == nil {
		return nil, fmt.Errorf("no exporter secret")
	}

	userIDNum, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing user ID: %w", err)
	}
	context := make([]byte, 8)
	binary.LittleEndian.PutUint64(context, userIDNum)

	baseSecret, err := mls.Export(d.exporterSecret, daveExportLabel, context, daveKeySize)
	if err != nil {
		return nil, fmt.Errorf("exporting receiver base secret: %w", err)
	}

	key, err := hashRatchetGetKey(baseSecret, 0)
	if err != nil {
		return nil, fmt.Errorf("deriving receiver ratchet key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	fc, err := newDAVECipher(key)
	if err != nil {
		return nil, err
	}

	recv := &daveReceiver{
		userID:      userID,
		baseSecret:  baseSecret,
		key:         key,
		aesBlock:    block,
		frameCipher: fc,
	}

	if d.receivers == nil {
		d.receivers = make(map[uint32]*daveReceiver)
	}
	d.receivers[ssrc] = recv
	return recv, nil
}

func (d *DAVESession) clearReceiversLocked() {
	d.receivers = nil
}

func (d *DAVESession) CanEncrypt() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.frameCipher != nil
}

func (d *DAVESession) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.exporterSecret = nil
	d.senderKey = nil
	d.senderNonce = 0
	d.frameCipher = nil
	d.active = false
	d.kpBundle = nil
	d.pendingTransitionID = 0
	d.pendingVersion = 0
	d.ratchetBaseSecret = nil
	d.currentGeneration = 0
	d.hasPendingKey = false
	d.ssrcToUserID = nil
	d.clearReceiversLocked()
}
