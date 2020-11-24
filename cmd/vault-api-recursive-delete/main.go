package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gturetsky/vault-api-recursive-delete/vaultdelete"
)

func main() {
	vaultAddress := flag.String("VAULT_ADDR", os.Getenv("VAULT_ADDR"), "Set the VAULT_ADDR")
	vaultToken := flag.String("VAULT_TOKEN", os.Getenv("VAULT_TOKEN"), "Set the VAULT_TOKEN")
	vaultCert := flag.String("VAULT_CACERT", os.Getenv("VAULT_CACERT"), "Set the VAULT_CACERT")
	path := flag.String("path", "", "Path to recursively delete")
	flag.Parse()

	client, err := vaultdelete.NewVaultClient("v2", *vaultAddress, *vaultToken, []string{*vaultCert})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("These paths will be deleted:")

	err = client.RecursiveDelete(*path, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Do you want to proceed? (Y/N)")

	var proceed string
	fmt.Scanln(&proceed)
	if proceed != "Y" {
		return
	}

	err = client.RecursiveDelete(*path, true)
	if err != nil {
		log.Fatal(err)
	}
}
