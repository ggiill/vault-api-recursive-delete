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
	"sort"
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

// ListResults gets keys for a path as an array of strings.
func (v *VaultClient) ListResults(dataPath string) ([]string, error) {
	res, err := v.Request("list", dataPath, nil)
	if err != nil {
		return nil, err
	}
	j := new(map[string]interface{})
	err = json.Unmarshal(res, j)
	if err != nil {
		return nil, err
	}
	data, ok := (*j)["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("no ['data']")
	}
	keysTmp, ok := data["keys"].([]interface{})
	if !ok {
		return nil, errors.New("no ['keys']")
	}
	keys := make([]string, len(keysTmp))
	for i, k := range keysTmp {
		key, ok := k.(string)
		if ok {
			keys[i] = key
		} else {
			return nil, errors.New("issue converting 'keys' interface to string")
		}
	}
	return keys, nil
}

// GetResults gets the values for a key AT a specific path.
func (v *VaultClient) GetResults(dataPath string) ([]string, error) {
	res, err := v.Request("get", dataPath, nil)
	if err != nil {
		return nil, err
	}
	j := new(map[string]interface{})
	err = json.Unmarshal(res, j)
	if err != nil {
		return nil, err
	}
	data, ok := (*j)["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("no ['data']")
	}
	_, ok = data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("no ['data']['data']")
	}
	return []string{dataPath}, nil
}

// GetPaths gets paths to delete.
func (v *VaultClient) GetPaths(dataPath string) ([]string, error) {
	var subErr error
	paths := make(map[string]bool) // This is a set, since "a/b/" and "a/b" can both exist as "a/b" paths
	var p func(dataPath string) error
	p = func(dataPath string) error {
		keys, err := v.ListResults(dataPath)
		if err != nil {
			keys, subErr = v.GetResults(dataPath)
			if subErr != nil {
				return fmt.Errorf("%s | %s", err, subErr)
			}
			if keys == nil {
				return err
			}
		}
		for _, key := range keys {
			newPath := path.Join(dataPath, key)
			paths[newPath] = true
			if strings.HasSuffix(key, "/") {
				err := p(newPath)
				if err != nil {
					return err
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
	sort.Strings(res)
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
