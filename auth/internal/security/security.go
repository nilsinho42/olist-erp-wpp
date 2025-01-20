package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Security package encrypts and decrypts the token

func EncryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize()) // nonce is an arrery of bytes of length gcm.NonceSize()
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil) // when creating the ciphertext, Seal appends the nonce to begin of it

	return ciphertext, nil
}

func DecryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	noncesize := gcm.NonceSize() // not the nonce itself, but the size of the nonce!
	if len(ciphertext) < noncesize {
		return nil, fmt.Errorf("invalid input")
	}
	nonce, cipherbytes := ciphertext[:noncesize], ciphertext[noncesize:]

	plaintext, err := gcm.Open(nil, nonce, cipherbytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
