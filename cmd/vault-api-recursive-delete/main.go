package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gturetsky/vault-api-recursive-delete/vaultdelete"
)

func main() {
	vaultAddress := flag.String("VAULT_ADDR", os.Getenv("VAULT_ADDR"), "Set the VAULT_ADDR")
	vaultToken := flag.String("VAULT_TOKEN", os.Getenv("VAULT_TOKEN"), "Set the VAULT_TOKEN")
	vaultCert := flag.String("VAULT_CACERT", os.Getenv("VAULT_CACERT"), "Set the VAULT_CACERT")
	path := flag.String("path", "", "Path to recursively delete")
	deleteMetadata := flag.Bool("delete-metadata", false, "Delete metadata as well")
	flag.Parse()

	client, err := vaultdelete.NewVaultClient("v2", *vaultAddress, *vaultToken, []string{*vaultCert})
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	paths, err := client.GetPaths(*path)
	if err != nil {
		if fmt.Sprint(err) != "no ['data']" {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}

	if len(paths) == 0 {
		fmt.Println("No paths found")
		return
	}

	fmt.Println("These paths will be deleted:")

	for _, path := range paths {
		fmt.Println(path)
	}

	fmt.Println("Do you want to proceed? (Y/N)")

	var proceed string
	fmt.Scanln(&proceed)
	if proceed != "Y" {
		return
	}

	for _, path := range paths {
		_, err := client.Request("delete", path, nil)
		fmt.Println("deleted:", path)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if *deleteMetadata {
			_, err := client.Request("deleteMetadata", path, nil)
			if err != nil {
				fmt.Println("error:", err)
				os.Exit(1)
			}
		}
	}
}
