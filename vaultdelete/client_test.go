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
	fmt.Println("**SETUP**")
	u := uuid.New()
	uu := u.String()
	fmt.Println("Docker container name:", uu)
	vaultToken := Setup(uu)
	err := os.Setenv("VAULT_TOKEN_TESTING", vaultToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("vaultToken:", vaultToken)
	fmt.Println("**END SETUP**")
	exitVal := m.Run()
	fmt.Println("**TEARDOWN**")
	TearDown(uu)
	fmt.Println("**END TEARDOWN**")
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
	resp, err := client.Request("post", "a/b/c/d", bytes.NewReader(dataBytes))
	fmt.Println("created secret:", string(resp))
	err = client.RecursiveDelete("a", true)
	if err != nil {
		t.Error(err)
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
	fmt.Println("Docker killed & pruned container:", string(out))

}
