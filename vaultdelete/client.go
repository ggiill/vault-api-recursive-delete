package vaultdelete

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// VaultClient is a vault API client.
type VaultClient struct {
	client    *http.Client
	version   string
	address   string
	token     string
	certPaths []string
}

// Request makes an http request.
func (v *VaultClient) Request(endpoint string, dataPath string, body io.Reader) ([]byte, error) {
	info, ok := versionPaths[v.version][endpoint]
	if !ok {
		return nil, fmt.Errorf("endpoint '%s' not mapped for version '%s'", endpoint, v.version)
	}
	reqURL, err := url.Parse(v.address)
	if err != nil {
		return nil, fmt.Errorf("invalid Vault address URL: '%s'", v.address)
	}
	reqURL.Path = path.Join(reqURL.Path, info.path, dataPath)
	reqPath := reqURL.String()
	req, err := http.NewRequest(info.method, reqPath, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", v.token)
	res, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

// GetPaths gets paths to delete.
func (v *VaultClient) GetPaths(dataPath string) ([]string, error) {
	paths := make(map[string]bool) // This is a set, since "a/b/" and "a/b" can both exist as "a/b" paths
	var p func(dataPath string) error
	p = func(dataPath string) error {
		res, err := v.Request("list", dataPath, nil)
		if err != nil {
			return err
		}
		j := new(map[string]interface{})
		err = json.Unmarshal(res, j)
		if err != nil {
			return err
		}
		data, ok := (*j)["data"].(map[string]interface{})
		if !ok {
			return errors.New("no ['data']")
		}
		keysTmp, ok := data["keys"].([]interface{})
		if !ok {
			return errors.New("no ['keys']")
		}
		for _, k := range keysTmp {
			if key, ok := k.(string); ok {
				newPath := path.Join(dataPath, key)
				paths[newPath] = true
				if strings.HasSuffix(key, "/") {
					err := p(newPath)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
	err := p(dataPath)
	if err != nil {
		return nil, err
	}
	res := make([]string, len(paths))
	pos := 0
	for k := range paths {
		res[pos] = k
		pos++
	}
	return res, nil
}

// NewVaultClient creates a new VaultClient.
func NewVaultClient(version, address, token string, certPaths []string) (*VaultClient, error) {
	caCertPool := x509.NewCertPool()
	for _, v := range certPaths {
		caCert, err := ioutil.ReadFile(v)
		if err != nil {
			return nil, fmt.Errorf("could not load cert from %s", v)
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	client := &VaultClient{
		client:  httpClient,
		version: version,
		address: address,
		token:   token,
	}
	return client, nil
}
