// aes
package controllers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	//"encoding/hex"
	"errors"
	"io"
)

func Encrypt(plaintext []byte, key string) ([]byte, error) {
	if len(plaintext)%aes.BlockSize != 0 {
		b := make([]byte, aes.BlockSize-len(plaintext)%aes.BlockSize+len(plaintext)) // padding
		copy(b, plaintext)
		plaintext = b

		//return nil, errors.New("plaintext is not a multiple of the block size")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(cblock, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func Decrypt(ciphertext []byte, key string) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	cblock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(cblock, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// trim suffix padding byte 0
	if index := bytes.IndexByte(ciphertext, 0); index > 0 {
		ciphertext = ciphertext[:index]
	}

	return ciphertext, nil
}
