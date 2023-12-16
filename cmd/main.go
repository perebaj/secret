package main

import (
	secert "github.com/perebaj/secret"
)

func main() {
	v := secert.NewVault("my-fake-key")
	err := v.Set("demo", "this is some demo")
	if err != nil {
		panic(err)
	}
	plain, err := v.Get("demo")
	if err != nil {
		panic(err)
	}
	println(plain)
}
