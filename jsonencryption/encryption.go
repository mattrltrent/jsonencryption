package jsonencryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
)

func NewSecret(maskSecret string) *Secret {
	return &Secret{
		val: []byte(maskSecret),
	}
}

type Secret struct {
	val []byte
}

func Hash(input uint) string {
	hash := sha256.Sum256([]byte(fmt.Sprint(input)))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (e Secret) Mask(id uint) (string, error) {
	if len(e.val) != 16 {
		return "", ErrInvalidKey
	}
	block, err := aes.NewCipher(e.val)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(fmt.Sprintf("%d", id)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	ctr := cipher.NewCTR(block, iv)
	plaintext := []byte(fmt.Sprintf("%d", id))
	ctr.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e Secret) Unmask(ciphertext string) (uint, error) {
	if len(e.val) != 16 {
		return 0, ErrInvalidKey
	}
	block, err := aes.NewCipher(e.val)
	if err != nil {
		return 0, err
	}

	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return 0, err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return 0, ErrInvalidKey
	}

	iv := decodedCiphertext[:aes.BlockSize]
	if len(decodedCiphertext) <= aes.BlockSize {
		return 0, ErrInvalidKey
	}

	ctr := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(decodedCiphertext)-aes.BlockSize)
	ctr.XORKeyStream(plaintext, decodedCiphertext[aes.BlockSize:])

	decryptedID, err := strconv.Atoi(string(plaintext))
	if err != nil {
		return 0, err
	}
	return uint(decryptedID), nil
}
