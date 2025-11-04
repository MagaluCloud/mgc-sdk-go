package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/sshkeys"
)

func main() {
	ExampleListSSHKeys()
	id := ExampleCreateSSHKey()
	ExampleGetSSHKey(id)
	ExampleDeleteSSHKey(id)
}

func ExampleListSSHKeys() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	sshClient := sshkeys.New(c)

	keys, err := sshClient.Keys().List(context.Background(), sshkeys.ListOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d SSH keys:\n", len(keys))
	for _, key := range keys {
		fmt.Printf("SSH Key: %s (ID: %s)\n", key.Name, key.ID)
		fmt.Printf("  Type: %s\n", key.KeyType)
	}
}

func ExampleCreateSSHKey() string {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	sshClient := sshkeys.New(c)

	key, err := sshClient.Keys().Create(context.Background(), sshkeys.CreateSSHKeyRequest{
		Name: "example-key",
		Key:  "ssh-rsa AAAA... example@localhost",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created SSH key: %s (ID: %s)\n", key.Name, key.ID)
	return key.ID
}

func ExampleGetSSHKey(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	sshClient := sshkeys.New(c)

	key, err := sshClient.Keys().Get(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SSH Key Details:\n")
	fmt.Printf("  ID: %s\n", key.ID)
	fmt.Printf("  Name: %s\n", key.Name)
	fmt.Printf("  Type: %s\n", key.KeyType)
}

func ExampleDeleteSSHKey(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	sshClient := sshkeys.New(c)

	key, err := sshClient.Keys().Delete(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully deleted SSH key: %s (ID: %s)\n", key.Name, key.ID)
}
