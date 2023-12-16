package secert

import (
	"fmt"

	"github.com/perebaj/secret/encrypt"
)

type Vault struct {
	encodingKey string
	keyValues   map[string]string
}

func NewVault(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
	}
}

func (v Vault) Get(key string) (string, error) {
	hex, ok := v.keyValues[key]
	if !ok {
		return "", fmt.Errorf("no value for key: %s", key)
	}
	ret, err := encrypt.Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", fmt.Errorf("error decrypting value: %s", err.Error())
	}

	return ret, nil
}

func (v Vault) Set(key, value string) error {
	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return fmt.Errorf("error encrypting value: %s", err.Error())
	}
	v.keyValues[key] = encryptedValue
	return nil
}
