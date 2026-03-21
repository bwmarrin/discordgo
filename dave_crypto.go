package discordgo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"fmt"

	"github.com/bwmarrin/discordgo/mls"
)

var errNotDAVEFrame = fmt.Errorf("not a DAVE frame")

const (
	daveTagSize              = 8
	daveKeySize              = 16
	daveExportLabel          = "Discord Secure Frames v0"
	minSupplementalBytesSize = daveTagSize + 1 + 1 + 2 // tag + nonce(min 1) + sizeB + magic = 12
)

func encryptSecureFrame(frameCipher cipher.AEAD, nonce uint32, opusData []byte) []byte {
	fullNonce := buildNonce(nonce)
	sealed := frameCipher.Seal(nil, fullNonce, opusData, nil)

	ciphertext := sealed[:len(opusData)]
	fullTag := sealed[len(opusData):]
	truncatedTag := fullTag[:daveTagSize]

	nonceBytes := encodeULEB128(nonce)

	supplementalSize := byte(daveTagSize + len(nonceBytes) + 1 + 2)

	result := make([]byte, 0, len(ciphertext)+daveTagSize+len(nonceBytes)+3)
	result = append(result, ciphertext...)
	result = append(result, truncatedTag...)
	result = append(result, nonceBytes...)
	result = append(result, supplementalSize)
	result = append(result, 0xFA, 0xFA)
	return result
}

func buildNonce(counter uint32) []byte {
	nonce := make([]byte, 12)
	binary.LittleEndian.PutUint32(nonce[8:], counter)
	return nonce
}

func encodeULEB128(value uint32) []byte {
	if value == 0 {
		return []byte{0}
	}
	var result []byte
	for value > 0 {
		b := byte(value & 0x7F)
		value >>= 7
		if value > 0 {
			b |= 0x80
		}
		result = append(result, b)
	}
	return result
}

func decodeULEB128(data []byte) (uint32, int) {
	var result uint32
	var shift uint
	for i, b := range data {
		result |= uint32(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, i + 1
		}
		shift += 7
	}
	return result, len(data)
}

func parseSecureFrame(data []byte) (ciphertext, truncatedTag []byte, nonce uint32, err error) {
	if len(data) < 2+1+1+daveTagSize {
		err = fmt.Errorf("secure frame too short: %d bytes", len(data))
		return
	}

	if data[len(data)-1] != 0xFA || data[len(data)-2] != 0xFA {
		err = errNotDAVEFrame
		return
	}

	supplementalSize := int(data[len(data)-3])
	supplementalStart := len(data) - supplementalSize

	if supplementalStart < 0 || supplementalSize < minSupplementalBytesSize {
		err = fmt.Errorf("invalid supplemental size %d for frame of %d bytes", supplementalSize, len(data))
		return
	}

	ciphertext = data[:supplementalStart]

	nonceBytes := data[supplementalStart+daveTagSize : len(data)-3]
	nonce, _ = decodeULEB128(nonceBytes)

	truncatedTag = data[supplementalStart : supplementalStart+daveTagSize]

	return
}

func decryptSecureFrame(aesBlock cipher.Block, frameCipher cipher.AEAD, nonce uint32, ciphertext, truncatedTag []byte) ([]byte, error) {
	gcmNonce := buildNonce(nonce)

	ctrIV := make([]byte, aes.BlockSize)
	copy(ctrIV, gcmNonce)
	binary.BigEndian.PutUint32(ctrIV[12:], 2)

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCTR(aesBlock, ctrIV).XORKeyStream(plaintext, ciphertext)

	sealed := frameCipher.Seal(nil, gcmNonce, plaintext, nil)
	fullTag := sealed[len(plaintext):]

	if subtle.ConstantTimeCompare(fullTag[:daveTagSize], truncatedTag) != 1 {
		return nil, fmt.Errorf("DAVE tag verification failed")
	}

	return plaintext, nil
}

func newDAVECipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func hashRatchetGetKey(baseSecret []byte, generation uint32) ([]byte, error) {
	secret := baseSecret
	for i := uint32(0); i < generation; i++ {
		genCtx := make([]byte, 4)
		binary.BigEndian.PutUint32(genCtx, i)
		next, err := mls.ExpandWithLabel(secret, "secret", genCtx, 32)
		if err != nil {
			return nil, err
		}
		secret = next
	}
	genCtx := make([]byte, 4)
	binary.BigEndian.PutUint32(genCtx, generation)
	return mls.ExpandWithLabel(secret, "key", genCtx, 16)
}
