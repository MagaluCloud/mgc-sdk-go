package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/containerregistry"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken), client.WithBaseURL(client.BrSe1), client.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))))

	getCredentials(c)
	resetPassword(c)
	listRegistries(c)
	id := createRegistry(c)
	getRegistry(c, id)

	/*
		// Unique test - we need push an image to the registry to test the repository list
		preId := "a6f44d78-30b4-44ea-bce6-4f11b9adac2e"
		name := listRepositories(c, preId)
		getRepository(c, preId, name)
		getImages(c, preId, name)

		deleteRepository(c, id, name)
	*/
	deleteRegistry(c, id)

	// Create Container Registry client
	containerRegistryClient := containerregistry.New(c)

	id = createProxyCache(containerRegistryClient)
	listProxyCaches(containerRegistryClient)
	listAllProxyCaches(containerRegistryClient)
	updateProxyCache(containerRegistryClient, id)
	getProxyCache(containerRegistryClient, id)
	deleteProxyCache(containerRegistryClient, id)
}

func getCredentials(c *client.CoreClient) {
	containerRegistryClient := containerregistry.New(c)
	credentials, err := containerRegistryClient.Credentials().Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Credentials: ")
	fmt.Println("User: ", credentials.Username, "Password: ", "[REDACTED]", "Email: ", credentials.Email)
}

func resetPassword(c *client.CoreClient) {
	containerRegistryClient := containerregistry.New(c)
	credentials, err := containerRegistryClient.Credentials().ResetPassword(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Password reset successfully, new credentials:")
	fmt.Println("User: ", credentials.Username, "Password: ", "[REDACTED]", "Email: ", credentials.Email)
}

func listRegistries(c *client.CoreClient) {
	containerRegistryClient := containerregistry.New(c)
	registries, err := containerRegistryClient.Registries().List(context.Background(), containerregistry.RegistryListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, registry := range registries.Results {
		fmt.Println("Registry: ", registry.Name, "Storage: ", registry.Storage, "CreatedAt: ", registry.CreatedAt, "UpdatedAt: ", registry.UpdatedAt)
	}
}

func createRegistry(c *client.CoreClient) string {
	containerRegistryClient := containerregistry.New(c)
	registry, err := containerRegistryClient.Registries().Create(context.Background(), &containerregistry.RegistryRequest{
		Name: "test-registry" + strconv.Itoa(int(time.Now().Unix())),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registry created: ", registry.Name)
	return registry.ID

}

func getRegistry(c *client.CoreClient, id string) {
	containerRegistryClient := containerregistry.New(c)
	registry, err := containerRegistryClient.Registries().Get(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registry: ", registry.Name, "Storage: ", registry.Storage, "CreatedAt: ", registry.CreatedAt, "UpdatedAt: ", registry.UpdatedAt)
}

func deleteRegistry(c *client.CoreClient, id string) {
	containerRegistryClient := containerregistry.New(c)
	err := containerRegistryClient.Registries().Delete(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registry deleted: ", id)
}

func listRepositories(c *client.CoreClient, id string) string {
	containerRegistryClient := containerregistry.New(c)
	repositories, err := containerRegistryClient.Repositories().List(context.Background(), id, containerregistry.RepositoryListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, repository := range repositories.Results {
		fmt.Println("Repository: ", repository.Name, "ImageCount: ", repository.ImageCount, "CreatedAt: ", repository.CreatedAt, "UpdatedAt: ", repository.UpdatedAt)
	}
	if len(repositories.Results) == 0 {
		return ""
	}
	return repositories.Results[0].Name
}

func getRepository(c *client.CoreClient, id string, name string) {
	containerRegistryClient := containerregistry.New(c)
	repository, err := containerRegistryClient.Repositories().Get(context.Background(), id, name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Repository: ", repository.Name, "ImageCount: ", repository.ImageCount, "CreatedAt: ", repository.CreatedAt, "UpdatedAt: ", repository.UpdatedAt)
}

func deleteRepository(c *client.CoreClient, id string, name string) {
	containerRegistryClient := containerregistry.New(c)
	err := containerRegistryClient.Repositories().Delete(context.Background(), id, name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Repository deleted: ", name)
}

func getImages(c *client.CoreClient, id string, name string) {
	containerRegistryClient := containerregistry.New(c)
	images, err := containerRegistryClient.Images().List(context.Background(), id, name, containerregistry.ImageListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Images: ", images)
	for _, image := range images.Results {
		fmt.Println("Image: ", image.Digest, "SizeBytes: ", image.SizeBytes, "PushedAt: ", image.PushedAt, "PulledAt: ", image.PulledAt, "ManifestMediaType: ", image.ManifestMediaType, "MediaType: ", image.MediaType, "Tags: ", image.Tags, "TagsDetails: ", image.TagsDetails, "ExtraAttr: ", image.ExtraAttr)
	}
}

func listProxyCaches(crClient *containerregistry.ContainerRegistryClient) {
	// List proxy caches with pagination
	resp, err := crClient.ProxyCaches().List(context.Background(), containerregistry.ProxyCacheListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Sort:   helpers.StrPtr("name:desc"),
	})

	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}

		log.Fatal("Failed to list proxy-caches")
	}

	// Print pagination metadata
	fmt.Printf("Proxy-caches (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	fmt.Println()

	// Print proxy-caches details
	for _, proxy := range resp.Results {

		fmt.Printf("ID: %s\n", proxy.ID)
		fmt.Printf("Name: %s\n", proxy.Name)
		fmt.Printf("Provider: %s\n", proxy.Provider)
		fmt.Printf("URL: %s\n", proxy.URL)
		fmt.Printf("Created at: %s\n", proxy.CreatedAt)
		fmt.Printf("Updated at: %s\n", proxy.UpdatedAt)

		fmt.Printf("-------------------------------------------\n")
	}

}

func listAllProxyCaches(crClient *containerregistry.ContainerRegistryClient) {
	// List proxy caches with pagination
	resp, err := crClient.ProxyCaches().ListAll(context.Background(), containerregistry.ProxyCacheListAllOptions{
		Sort: helpers.StrPtr("name:desc"),
	})

	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}

		log.Fatal("Failed to list all proxy-caches")
	}

	fmt.Printf("All proxy-caches:\n\n")

	// Print proxy-caches details
	for _, proxy := range resp {
		fmt.Printf("ID: %s\n", proxy.ID)
		fmt.Printf("Name: %s\n", proxy.Name)
		fmt.Printf("Provider: %s\n", proxy.Provider)
		fmt.Printf("URL: %s\n", proxy.URL)
		fmt.Printf("Created at: %s\n", proxy.CreatedAt)
		fmt.Printf("Updated at: %s\n", proxy.UpdatedAt)

		fmt.Printf("-------------------------------------------\n")
	}
}

func createProxyCache(crClient *containerregistry.ContainerRegistryClient) string {
	resp, err := crClient.ProxyCaches().Create(context.Background(), containerregistry.CreateProxyCacheRequest{
		Name:        "test-proxy-cache",
		Provider:    "docker-hub",
		URL:         "https://hub.docker.com/repositories",
		Description: helpers.StrPtr("Proxy-cache description"),
	})

	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}

		log.Fatal("Failed to create a proxy-cache")
	}

	fmt.Printf("Created proxy-cache with ID: %s\n", resp.ID)

	return resp.ID
}

func deleteProxyCache(crClient *containerregistry.ContainerRegistryClient, id string) {
	err := crClient.ProxyCaches().Delete(context.Background(), id)

	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}

		log.Fatal("Failed to delete proxy-cache")
	}

	fmt.Printf("Deleted proxy-cache with ID: %s\n", id)
}

func updateProxyCache(crClient *containerregistry.ContainerRegistryClient, id string) {
	resp, err := crClient.ProxyCaches().Update(context.Background(), id, containerregistry.UpdateProxyCacheRequest{
		Name: helpers.StrPtr("new-test-proxy-cache"),
	})

	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}

		log.Fatal("Failed to update proxy-cache")
	}

	fmt.Printf("Updated proxy-cache:\n")

	fmt.Printf("ID: %s\n", resp.ID)
	fmt.Printf("Name: %s\n", resp.Name)
	fmt.Printf("Provider: %s\n", resp.Provider)
	fmt.Printf("URL: %s\n", resp.URL)
	fmt.Printf("Created at: %s\n", resp.CreatedAt)
	fmt.Printf("Updated at: %s\n", resp.UpdatedAt)
}

func getProxyCache(crClient *containerregistry.ContainerRegistryClient, id string) {
	resp, err := crClient.ProxyCaches().Get(context.Background(), id)

	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}

		log.Fatal("Failed to get a proxy-cache")
	}

	fmt.Printf("ID: %s\n", resp.ID)
	fmt.Printf("Name: %s\n", resp.Name)
	fmt.Printf("Description: %s\n", resp.Description)
	fmt.Printf("Provider: %s\n", resp.Provider)
	fmt.Printf("URL: %s\n", resp.URL)
	fmt.Printf("Created at: %s\n", resp.CreatedAt)
	fmt.Printf("Updated at: %s\n", resp.UpdatedAt)
}
