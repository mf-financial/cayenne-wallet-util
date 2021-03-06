package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

// Crypt is for cipher config data
type Crypt struct {
	cipher cipher.Block
	iv     []byte
}

// Creates a new encryption/decryption object
// with a given key of a given size
// (16, 24 or 32 for AES-128, AES-192 and AES-256 respectively,
// as per http://golang.org/pkg/crypto/aes/#NewCipher)
//
// The key will be padded to the given size if needed.
// An IV is created as a series of NULL bytes of necessary length
// when there is no iv string passed as 3rd value to function.

// NewCryptWithParam is to create crypt instance
// key size should be 16,24,32
// iv size should be 16
func NewCrypt(key, iv string) (Crypter, error) {
	if len(iv) != aes.BlockSize {
		return nil, errors.Errorf("iv size should be %d", aes.BlockSize)
	}

	padded := make([]byte, len(key))
	copy(padded, []byte(key))

	bIv := []byte(iv)
	block, err := aes.NewCipher(padded)
	if err != nil {
		return nil, err
	}

	cryptInfo := Crypt{block, bIv}

	return &cryptInfo, nil
}

func (c *Crypt) padSlice(src []byte) []byte {
	// src must be a multiple of block size
	mult := int((len(src) / aes.BlockSize) + 1)
	leng := aes.BlockSize * mult

	srcPadded := make([]byte, leng)
	copy(srcPadded, src)
	return srcPadded
}

// Encrypt is encrypt a slice of bytes, producing a new, freshly allocated slice
// Source will be padded with null bytes if necessary
func (c *Crypt) Encrypt(src []byte) []byte {
	if len(src)%aes.BlockSize != 0 {
		src = c.padSlice(src)
	}
	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(c.cipher, c.iv).CryptBlocks(dst, src)
	return dst
}

// EncryptBase64 is encrypt and encode by base64 string
func (c *Crypt) EncryptBase64(plainText string) string {
	encryptedBytes := c.Encrypt([]byte(plainText))
	base64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	return base64
}

// Decrypt is to decrypt a slice of bytes, producing a new, freshly allocated slice
// Source will be padded with null bytes if necessary
func (c *Crypt) Decrypt(src []byte) []byte {
	if len(src)%aes.BlockSize != 0 {
		src = c.padSlice(src)
	}
	dst := make([]byte, len(src))
	cipher.NewCBCDecrypter(c.cipher, c.iv).CryptBlocks(dst, src)
	trimmed := bytes.Trim(dst, "\x00")
	return trimmed
}

// DecryptBase64 is to decrypt decoded Base64 string
func (c *Crypt) DecryptBase64(base64String string) (string, error) {
	unbase64, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	decryptedBytes := c.Decrypt(unbase64)
	return string(decryptedBytes[:]), nil
}
