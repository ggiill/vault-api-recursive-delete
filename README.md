# vault-api-recursive-delete
Recursively delete paths on Vault's KV engine.

### CLI
```shell
$ cd cmd/vault-api-recursive-delete

# Create some secrets under the "test/" path
$ vault kv put secret/t hey=yo
$ vault kv put secret/t/a hey=yo
$ vault kv put secret/t/a/b hey=yo
$ vault kv put secret/t/a/b/c hey=yo

# Check the CLI args
$ go run main.go --help
  -VAULT_ADDR string
        Set the VAULT_ADDR
  -VAULT_CACERT string
        Set the VAULT_CACERT
  -VAULT_TOKEN string
        Set the VAULT_TOKEN
  -delete-metadata
        Delete metadata as well
  -path string
        Path to recursively delete

# Recursively delete secrets and path metadata
$ go run main.go --path t --delete-metadata
These paths will be deleted:
t/a/b/c
t/a/b
t/a
t
Do you want to proceed? (Y/N)
Y
deleted: t/a/b/c
deleted: t/a/b
deleted: t/a
deleted: t
```
