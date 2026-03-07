package discordgo

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

const (
	daveTagSize     = 8
	daveKeySize     = 16
	daveExportLabel = "Discord Secure Frames v0"
)

func encryptSecureFrame(frameCipher cipher.AEAD, nonce uint32, opusData []byte) []byte {
	fullNonce := buildNonce(nonce)
	sealed := frameCipher.Seal(nil, fullNonce, opusData, nil)

	ciphertext := sealed[:len(opusData)]
	fullTag := sealed[len(opusData):]
	truncatedTag := fullTag[:daveTagSize]

	nonceBytes := encodeULEB128(nonce)

	supplementalSize := byte(daveTagSize + len(nonceBytes) + 0 + 1 + 2)

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
		next, err := mlsExpandWithLabel(secret, "secret", genCtx, 32)
		if err != nil {
			return nil, err
		}
		secret = next
	}
	genCtx := make([]byte, 4)
	binary.BigEndian.PutUint32(genCtx, generation)
	return mlsExpandWithLabel(secret, "key", genCtx, 16)
}
