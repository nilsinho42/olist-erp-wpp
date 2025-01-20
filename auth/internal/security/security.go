package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Security package encrypts and decrypts the token

// EncryptAES encrypts the token using a env key (hex-encoded), returns the token (string) AES encrypted
func EncryptAES(plaintext []byte) (string, error) {
	// encrypt> 1) get the key for encryption
	// "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	byte_securitykey := os.Getenv("ENCRYPTION_KEY")

	// encrypt> 2) decode the key from hex-encoded
	key, err := hex.DecodeString(byte_securitykey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize()) // nonce is an arrery of bytes of length gcm.NonceSize()
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil) // when creating the ciphertext, Seal appends the nonce to begin of it

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(ciphertext []byte) (string, error) {
	// decrypt> 1) get the key for encryption
	// "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	byte_securitykey := os.Getenv("ENCRYPTION_KEY")
	key, err := hex.DecodeString(byte_securitykey)
	if err != nil {
		return "", err
	}

	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	noncesize := gcm.NonceSize() // not the nonce itself, but the size of the nonce!
	if len(ciphertext) < noncesize {
		return "", fmt.Errorf("invalid input")
	}
	nonce, cipherbytes := ciphertext[:noncesize], ciphertext[noncesize:]

	plaintext, err := gcm.Open(nil, nonce, cipherbytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
