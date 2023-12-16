package secret

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/perebaj/secret/encrypt"
)

type Vault struct {
	// encodingKey is used to encrypt and decrypt values
	encodingKey string
	// keyValues is a map that holds the encrypted values
	keyValues map[string]string
	// path is the location of the file that holds the vault data
	path string
}

func NewVault(encodingKey, path string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
		path:        path,
	}
}

func NewFileVault(encodingKey, path string) (*Vault, error) {
	v := NewVault(encodingKey, path)
	err := v.load(path)
	return v, err
}

func (v *Vault) Get(key string) (string, error) {
	err := v.load(v.path)
	if err != nil {
		return "", fmt.Errorf("error loading vault: %s", err.Error())
	}

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

	err = v.load(v.path)
	if err != nil {
		return fmt.Errorf("error loading vault: %s", err.Error())
	}
	v.keyValues[key] = encryptedValue
	return v.write(v.path)
}

func (v *Vault) write(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err.Error())
	}
	defer f.Close()

	byteValue, err := json.Marshal(v.keyValues)
	if err != nil {
		return fmt.Errorf("error encoding vault: %s", err.Error())
	}
	_, err = f.Write(byteValue)
	if err != nil {
		return fmt.Errorf("error writing vault: %s", err.Error())
	}
	return nil
}

// load reads a file and decodes the JSON-encoded key-value pairs into the
// Vault's keyValues map. If the file does not exist, the keyValues map is
// initialized as an empty map.
func (v *Vault) load(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err.Error())
	}
	defer f.Close()

	byteValue, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err.Error())
	}

	if len(byteValue) == 0 {
		// empty file
		v.keyValues = make(map[string]string)
		return nil
	}

	err = json.Unmarshal(byteValue, &v.keyValues)
	if err != nil {
		return fmt.Errorf("error decoding vault: %s", err.Error())
	}
	return nil
}
