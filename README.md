# vault-api-recursive-delete
Recursively delete paths on Vault's KV engine.

### CLI
```shell
$ cd cmd/vault-api-recursive-delete

# Create some secrets under the "test/" path
$ vault kv put secret/test hey=yo
$ vault kv put secret/test/a hey=yo
$ vault kv put secret/test/a/b hey=yo
$ vault kv put secret/test/a/b/c hey=yo

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
$ go run main.go --path test --delete-metadata
These paths will be deleted:
test/a
test/a/b
test/a/b/c
Do you want to proceed? (Y/N)
Y
deleted: test/a
deleted: test/a/b
deleted: test/a/b/c
```
