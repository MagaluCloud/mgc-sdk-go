package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/containerregistry"
)

func main() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}

	c := client.NewMgcClient(
		client.WithAPIKey(apiToken),
		client.WithBaseURL(client.BrSe1),
		client.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	crClient := containerregistry.New(c)

	user := createUser(crClient)
	getUser(crClient)
	deleteUser(crClient, user.ID)
}

func createUser(crClient *containerregistry.ContainerRegistryClient) *containerregistry.UserResponse {
	user, err := crClient.Users().Create(context.Background())
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to create user")
	}

	fmt.Printf("User created:\n")
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Created at: %s\n", user.CreatedAt)
	fmt.Println("-------------------------------------------")

	return user
}

func getUser(crClient *containerregistry.ContainerRegistryClient) {
	user, err := crClient.Users().Get(context.Background())
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to get user")
	}

	fmt.Printf("Current user:\n")
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Created at: %s\n", user.CreatedAt)
	fmt.Println("-------------------------------------------")
}

func deleteUser(crClient *containerregistry.ContainerRegistryClient, userID string) {
	err := crClient.Users().Delete(context.Background(), userID)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to delete user")
	}

	fmt.Printf("User deleted: %s\n", userID)
}
