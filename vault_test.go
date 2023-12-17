package secret

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestVault_Load(t *testing.T) {
	f, err := os.CreateTemp("", "vault")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(`{"key1":"value1","key2":"value2"}`))
	if err != nil {
		t.Fatal(err)
	}

	want := make(map[string]string)
	want["key1"] = "value1"
	want["key2"] = "value2"

	got, err := NewFileVault("my-fake-key", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	eq := reflect.DeepEqual(want, got.keyValues)
	if !eq {
		t.Errorf("got %v, want %v", got.keyValues, want)
	}
}

func TestVault_Load_Empty(t *testing.T) {
	f, err := os.CreateTemp("", "vault")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	want := make(map[string]string)

	got, err := NewFileVault("my-fake-key", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	eq := reflect.DeepEqual(want, got.keyValues)
	if !eq {
		t.Errorf("got %v, want %v", got.keyValues, want)
	}
}

func TestVault_Write(t *testing.T) {
	f, err := os.CreateTemp("", "vault")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	want, err := NewFileVault("my-fake-key", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	want.keyValues["key1"] = "value1"
	want.keyValues["key2"] = "value2"

	err = want.write(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	byteValue, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	got := make(map[string]string)
	err = json.Unmarshal(byteValue, &got)
	if err != nil {
		t.Fatal(err)
	}

	eq := reflect.DeepEqual(want.keyValues, got)
	if !eq {
		t.Errorf("got %v, want %v", got, want.keyValues)
	}
}

func TestVault_SetGet(t *testing.T) {
	f, err := os.CreateTemp("", "vault")
	if err != nil {
		t.Fatal(err)
	}
	v, err := NewFileVault("my-fake-key", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	want := make(map[string]string)

	eq := reflect.DeepEqual(want, v.keyValues)
	if !eq {
		t.Errorf("got %v, want %v", v.keyValues, want)
	}

	err = v.Set("key1", "value1")
	if err != nil {
		t.Fatal(err)
	}

	s, err := v.Get("key1")
	if err != nil {
		t.Fatal(err)
	}

	if s != "value1" {
		t.Errorf("got %s, want %s", s, "value1")
	}
}
