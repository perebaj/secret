package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Encrypt will take in a key and plaintext and return a hex representation of the encrypted value.
func Encrypt(key, plaintext string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	initializationVector := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, initializationVector); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, initializationVector)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt will take in a key and cipherHex (hex representation of the ciphertext) and decrypt it
func Decrypt(key, cipherHex string) (string, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// initializationVector = initialization vector
	initializationVector := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, initializationVector)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	io.WriteString(hasher, key)
	return aes.NewCipher(hasher.Sum(nil))
}
