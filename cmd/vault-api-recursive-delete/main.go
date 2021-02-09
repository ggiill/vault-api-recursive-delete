package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ggiill/vault-api-recursive-delete/vaultdelete"
)

func main() {
	vaultAddress := flag.String("VAULT_ADDR", "", "Set the VAULT_ADDR")
	vaultToken := flag.String("VAULT_TOKEN", "", "Set the VAULT_TOKEN")
	vaultCert := flag.String("VAULT_CACERT", "", "Set the VAULT_CACERT")
	path := flag.String("path", "", "Path to recursively delete")
	deleteMetadata := flag.Bool("delete-metadata", false, "Delete metadata as well")
	flag.Parse()

	if *vaultAddress == "" {
		*vaultAddress = os.Getenv("VAULT_ADDR")
	}
	if *vaultToken == "" {
		*vaultToken = os.Getenv("VAULT_TOKEN")
	}
	if *vaultCert == "" {
		*vaultCert = os.Getenv("VAULT_CACERT")
	}

	r := vaultdelete.RunConfig{
		Version:        "v2",
		Address:        *vaultAddress,
		Token:          *vaultToken,
		CertPaths:      []string{*vaultCert},
		Path:           *path,
		Interactive:    true,
		DeleteMetadata: *deleteMetadata,
	}

	err := vaultdelete.Run(r)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
