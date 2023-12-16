package encrypt

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "right-key"
	encryptString, err := Encrypt(key, "teste")
	if err != nil {
		t.Fatal(err)
	}

	decryptValue, err := Decrypt(key, encryptString)
	if err != nil {
		t.Fatal(err)
	}

	if decryptValue != "teste" {
		t.Fatalf("Expected %s, got %s", "teste", decryptValue)
	}

	wrongValue, err := Decrypt("wrong-key", encryptString)
	if err != nil {
		t.Fatal(err)
	}

	if wrongValue == "teste" {
		t.Fatal("Expected wrong value, got teste")
	}
}
