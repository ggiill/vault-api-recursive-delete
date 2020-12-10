package vaultdelete

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	log.Print("**SETUP**")
	u := uuid.New()
	uu := u.String()
	log.Print("Docker container name:", uu)
	vaultToken := Setup(uu)
	err := os.Setenv("VAULT_TOKEN_TESTING", vaultToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("vaultToken:", vaultToken)
	log.Print("**END SETUP**")
	exitVal := m.Run()
	log.Print("**TEARDOWN**")
	TearDown(uu)
	log.Print("**END TEARDOWN**")
	os.Exit(exitVal)
}

func Test(t *testing.T) {
	client, err := NewVaultClient("v2", "http://0.0.0.0:8200", os.Getenv("VAULT_TOKEN_TESTING"), nil)
	if err != nil {
		t.Error(err)
	}
	data := map[string]map[string]string{
		"data": {
			"hey": "yo",
			"sup": "hi",
		},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	pathsToCreate := []string{
		"a",
		"a/b",
		"a/b/c",
		"a/b/f",
		"a/b/c/d",
		"a/b/c/e",
	}
	for _, ptc := range pathsToCreate {
		resp, err := client.Request("post", ptc, bytes.NewReader(dataBytes))
		if err != nil {
			t.Error(err)
		}
		t.Logf("created secret @ %v: %v | %v", ptc, data, string(resp))
	}
	paths, err := client.GetPaths("a")
	if err != nil {
		t.Error(err)
	}
	deleteTypes := map[string]string{
		"delete":         "deleted:",
		"deleteMetadata": "metadata deleted:",
	}
	for method, message := range deleteTypes {
		for _, path := range paths {
			_, err := client.Request(method, path, nil)
			t.Log(message, path)
			if err != nil {
				t.Error(err)
			}
		}
		pr, err := client.GetPaths("a")
		if err != nil {
			if fmt.Sprint(err) != "no ['data']" {
				t.Error(err)
			}
		}
		t.Log("Paths remaining:", pr)
	}
}

func Setup(u string) string {
	commandText := fmt.Sprintf("docker run --cap-add=IPC_LOCK -p 8200:8200 --name %v -l %v vault server -dev", u, u)
	cmd := exec.Command("bash", "-c", commandText)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer cmdReader.Close()
	scanner := bufio.NewScanner(cmdReader)
	scanner.Split(bufio.ScanRunes)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	output := ""
	re, err := regexp.Compile("Root Token: ([A-Za-z0-9.]+)\n")
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		output += scanner.Text()
		rei := re.FindAllStringSubmatch(output, -1)
		if len(rei) > 0 {
			match := rei[0][1]
			cmd.Process.Release()
			return match
		}
	}
	log.Fatal("did not find vault token")
	return ""
}

func TearDown(containerName string) {
	commandText := fmt.Sprintf("docker kill %v && docker container prune -f --filter label=%v", containerName, containerName)
	cmd := exec.Command("bash", "-c", commandText)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Docker killed & pruned container:", string(out))
}
