package vaultdelete

import (
	"errors"
	"fmt"
)

// RunConfig contains the configuration for the run.
type RunConfig struct {
	Version, Address, Token     string
	CertPaths                   []string
	Path                        string
	Interactive, DeleteMetadata bool
}

// Run runs the app.
func Run(r RunConfig) error {
	client, err := NewVaultClient(r.Version, r.Address, r.Token, r.CertPaths)
	if err != nil {
		return err
	}
	paths, err := client.GetPaths(r.Path)
	if err != nil {
		if fmt.Sprint(err) != "no ['data']" {
			return err
		}
	}
	if len(paths) == 0 {
		return errors.New("no paths found")
	}
	if r.Interactive {
		fmt.Println("These paths will be deleted:")
		for _, path := range paths {
			fmt.Println(path)
		}
		fmt.Println("Do you want to proceed? (Y/N)")
		var proceed string
		fmt.Scanln(&proceed)
		if proceed != "Y" {
			return nil
		}
	}
	for _, path := range paths {
		msg := "deleted:"
		_, err := client.Request("delete", path, nil)
		if err != nil {
			return err
		}
		if r.DeleteMetadata {
			_, err := client.Request("deleteMetadata", path, nil)
			if err != nil {
				fmt.Println(msg, path)
				return err
			}
			msg = "deleted (w/ metadata):"
		}
		fmt.Println(msg, path)
	}
	return nil
}
