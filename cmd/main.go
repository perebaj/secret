package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/perebaj/secret"
)

/*
Usage
secret set -encodingKey=mykey -path=./.vault -key=foo -value=bar
secret get -encodingKey=mykey -path=./.vault -key=foo

*/

func main() {
	help := flag.Bool("h", false, "show help")
	flag.Parse()

	home, _ := os.UserHomeDir()
	path := home + "/.vault"

	setCmd := flag.NewFlagSet("set", flag.ExitOnError)

	pathSet := setCmd.String("p", path, "the path to the vault file")
	keySet := setCmd.String("k", "", "the key to get from the vault")
	encodingKeySet := setCmd.String("e", "", "the encoding key")
	valueSet := setCmd.String("v", "", "the value to store in the vault")

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	pathGet := getCmd.String("p", path, "the path to the vault file")
	keyGet := getCmd.String("k", "", "the key to get from the vault")
	encodingKeyGet := getCmd.String("e", "", "the encoding key")

	if len(os.Args) < 2 || *help {
		fmt.Println("expected 'set' or 'get' subcommands")
		setCmd.Usage()
		getCmd.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "set":
		setCmd.Parse(os.Args[2:])

		if *keySet == "" || *encodingKeySet == "" || *valueSet == "" {
			fmt.Println("Wrong usage")
			setCmd.Usage()
			os.Exit(1)
		}

		v, err := secret.NewFileVault(*encodingKeySet, *pathSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = v.Set(*keySet, *valueSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "get":
		getCmd.Parse(os.Args[2:])
		if *keyGet == "" || *encodingKeyGet == "" {
			fmt.Println("Wrong usage")
			getCmd.Usage()
			os.Exit(1)
		}

		v, err := secret.NewFileVault(*encodingKeyGet, *pathGet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		value, err := v.Get(*keyGet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s=%s\n", *keyGet, value)
	}
}
