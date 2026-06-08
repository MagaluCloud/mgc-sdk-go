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
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}

	registryID := os.Getenv("MGC_REGISTRY_ID")
	if registryID == "" {
		log.Fatal("MGC_REGISTRY_ID environment variable is not set")
	}

	userID := os.Getenv("MGC_USER_ID")
	if userID == "" {
		log.Fatal("MGC_USER_ID environment variable is not set")
	}

	c := client.NewMgcClient(
		client.WithAPIKey(apiToken),
		client.WithBaseURL(client.BrSe1),
		client.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	crClient := containerregistry.New(c)

	member := addMember(crClient, registryID, userID)
	listMembers(crClient, registryID)
	listAllMembers(crClient, registryID)
	getMember(crClient, registryID, member.ID)
	updateMember(crClient, registryID, member.ID)
	deleteMember(crClient, registryID, member.ID)
}

func addMember(crClient *containerregistry.ContainerRegistryClient, registryID, userID string) *containerregistry.MemberResponse {
	member, err := crClient.Members().Add(context.Background(), registryID, containerregistry.MemberRequest{
		UserID: userID,
		Role:   helpers.StrPtr("guest"),
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to add member")
	}

	fmt.Printf("Member added:\n")
	printMember(member)

	return member
}

func listMembers(crClient *containerregistry.ContainerRegistryClient, registryID string) {
	resp, err := crClient.Members().List(context.Background(), registryID, containerregistry.MemberListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		MemberFilterOptions: containerregistry.MemberFilterOptions{
			Sort: helpers.StrPtr("created_at:desc"),
		},
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to list members")
	}

	fmt.Printf("Members (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, member := range resp.Results {
		printMember(&member)
	}
}

func listAllMembers(crClient *containerregistry.ContainerRegistryClient, registryID string) {
	resp, err := crClient.Members().ListAll(context.Background(), registryID, containerregistry.MemberFilterOptions{
		Sort: helpers.StrPtr("created_at:desc"),
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to list all members")
	}

	fmt.Printf("All members:\n\n")
	for _, member := range resp {
		printMember(&member)
	}
}

func getMember(crClient *containerregistry.ContainerRegistryClient, registryID, memberID string) {
	member, err := crClient.Members().Get(context.Background(), registryID, memberID)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to get member")
	}

	fmt.Printf("Member:\n")
	printMember(member)
}

func updateMember(crClient *containerregistry.ContainerRegistryClient, registryID, memberID string) {
	member, err := crClient.Members().Update(context.Background(), registryID, memberID, containerregistry.MemberUpdateRequest{
		Role: "developer",
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to update member")
	}

	fmt.Printf("Member updated:\n")
	printMember(member)
}

func deleteMember(crClient *containerregistry.ContainerRegistryClient, registryID, memberID string) {
	err := crClient.Members().Delete(context.Background(), registryID, memberID)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to delete member")
	}

	fmt.Printf("Member deleted: %s\n", memberID)
}

func printMember(member *containerregistry.MemberResponse) {
	fmt.Printf("ID: %s\n", member.ID)
	fmt.Printf("Registry ID: %s\n", member.RegistryID)
	fmt.Printf("User ID: %s\n", member.UserID)
	fmt.Printf("Role: %s\n", member.Role)
	fmt.Printf("Created at: %s\n", member.CreatedAt)
	fmt.Printf("Updated at: %s\n", member.UpdatedAt)
	fmt.Println("-------------------------------------------")
}
