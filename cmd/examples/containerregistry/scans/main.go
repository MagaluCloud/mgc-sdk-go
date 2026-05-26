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

	repositoryID := os.Getenv("MGC_REPOSITORY_ID")
	if repositoryID == "" {
		log.Fatal("MGC_REPOSITORY_ID environment variable is not set")
	}

	digestOrTag := os.Getenv("MGC_DIGEST_OR_TAG")
	if digestOrTag == "" {
		log.Fatal("MGC_DIGEST_OR_TAG environment variable is not set")
	}

	c := client.NewMgcClient(
		client.WithAPIKey(apiToken),
		client.WithBaseURL(client.BrSe1),
		client.WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
	)

	crClient := containerregistry.New(c)

	imageRef := containerregistry.ScheduleScanInput{
		RegistryID:   registryID,
		RepositoryID: repositoryID,
		DigestOrTag:  digestOrTag,
	}

	scheduled := scheduleScan(crClient, imageRef)
	listScansByImage(crClient, imageRef)
	listAllScansByImage(crClient, imageRef)
	getScan(crClient, scheduled.ID)
	listVulnerabilities(crClient, scheduled.ID)
	listAllVulnerabilities(crClient, scheduled.ID)
	stopScan(crClient, scheduled.ID)
}

func scheduleScan(crClient *containerregistry.ContainerRegistryClient, in containerregistry.ScheduleScanInput) *containerregistry.ImageScanScheduleResponse {
	scan, err := crClient.Scans().Schedule(context.Background(), in)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to schedule scan")
	}

	fmt.Printf("Scan scheduled:\n")
	fmt.Printf("ID: %s\n", scan.ID)
	fmt.Printf("Digest: %s\n", scan.Digest)
	fmt.Printf("Registry ID: %s\n", scan.RegistryID)
	fmt.Printf("Repository ID: %s\n", scan.RepositoryID)
	if scan.InputTag != nil {
		fmt.Printf("Input tag: %s\n", *scan.InputTag)
	}
	fmt.Printf("Status: %s\n", scan.Status)
	fmt.Printf("Created at: %s\n", scan.CreatedAt)
	fmt.Printf("Updated at: %s\n", scan.UpdatedAt)
	fmt.Println("-------------------------------------------")

	return scan
}

func listScansByImage(crClient *containerregistry.ContainerRegistryClient, in containerregistry.ScheduleScanInput) {
	listInput := containerregistry.ListByImageInput(in)
	resp, err := crClient.Scans().ListByImage(context.Background(), listInput, containerregistry.ListImageScansOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		ImageScansFilterOptions: containerregistry.ImageScansFilterOptions{
			Sort:   helpers.StrPtr("created_at:desc"),
			Status: helpers.StrPtr("completed"),
		},
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to list scans by image")
	}

	fmt.Printf("Scans (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, scan := range resp.Results {
		printScan(&scan)
	}
}

func listAllScansByImage(crClient *containerregistry.ContainerRegistryClient, in containerregistry.ScheduleScanInput) {
	listInput := containerregistry.ListByImageInput(in)
	scans, err := crClient.Scans().ListByImageAll(context.Background(), listInput, containerregistry.ImageScansFilterOptions{
		Sort: helpers.StrPtr("created_at:desc"),
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to list all scans by image")
	}

	fmt.Printf("All scans:\n\n")
	for _, scan := range scans {
		printScan(&scan)
	}
}

func getScan(crClient *containerregistry.ContainerRegistryClient, scanID string) {
	scan, err := crClient.Scans().Get(context.Background(), scanID)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to get scan")
	}

	fmt.Printf("Scan details:\n")
	printScan(scan)

	for _, child := range scan.ChildScans {
		fmt.Printf("Child scan:\n")
		fmt.Printf("  ID: %s\n", child.ID)
		fmt.Printf("  Digest: %s\n", child.Digest)
		fmt.Printf("  Status: %s\n", child.Status)
		if child.OS != nil {
			fmt.Printf("  OS: %s\n", *child.OS)
		}
		if child.Architecture != nil {
			fmt.Printf("  Architecture: %s\n", *child.Architecture)
		}
		printSeverity(child.SeveritySummary)
		fmt.Println("  -------------------------------------------")
	}
}

func listVulnerabilities(crClient *containerregistry.ContainerRegistryClient, scanID string) {
	resp, err := crClient.Scans().ListVulnerabilities(context.Background(), scanID, containerregistry.ListVulnerabilitiesOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		VulnerabilitiesFilterOptions: containerregistry.VulnerabilitiesFilterOptions{
			Sort:     helpers.StrPtr("severity:desc"),
			Severity: []string{"high", "critical"},
			Fixable:  helpers.BoolPtr(true),
		},
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to list vulnerabilities")
	}

	fmt.Printf("Vulnerabilities (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, vuln := range resp.Results {
		printVulnerability(&vuln)
	}
}

func listAllVulnerabilities(crClient *containerregistry.ContainerRegistryClient, scanID string) {
	vulns, err := crClient.Scans().ListVulnerabilitiesAll(context.Background(), scanID, containerregistry.VulnerabilitiesFilterOptions{
		Sort:     helpers.StrPtr("severity:desc"),
		Severity: []string{"critical"},
	})
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to list all vulnerabilities")
	}

	fmt.Printf("All vulnerabilities:\n\n")
	for _, vuln := range vulns {
		printVulnerability(&vuln)
	}
}

func stopScan(crClient *containerregistry.ContainerRegistryClient, scanID string) {
	err := crClient.Scans().Stop(context.Background(), scanID)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
		log.Fatal("Failed to stop scan")
	}

	fmt.Printf("Scan stopped: %s\n", scanID)
}

func printScan(scan *containerregistry.ImageScanResponse) {
	fmt.Printf("ID: %s\n", scan.ID)
	fmt.Printf("Digest: %s\n", scan.Digest)
	if scan.InputTag != nil {
		fmt.Printf("Input tag: %s\n", *scan.InputTag)
	}
	fmt.Printf("Status: %s\n", scan.Status)
	if scan.OS != nil {
		fmt.Printf("OS: %s\n", *scan.OS)
	}
	if scan.Architecture != nil {
		fmt.Printf("Architecture: %s\n", *scan.Architecture)
	}
	printSeverity(scan.SeveritySummary)
	fmt.Printf("Created at: %s\n", scan.CreatedAt)
	fmt.Printf("Updated at: %s\n", scan.UpdatedAt)
	if scan.FinishedAt != nil {
		fmt.Printf("Finished at: %s\n", *scan.FinishedAt)
	}
	fmt.Println("-------------------------------------------")
}

func printSeverity(s *containerregistry.SeverityResponse) {
	if s == nil {
		return
	}
	fmt.Printf("Severity summary: total=%d, critical=%d, high=%d, medium=%d, low=%d, fixable=%d\n",
		s.Total, s.Critical, s.High, s.Medium, s.Low, s.Fixable)
}

func printVulnerability(v *containerregistry.VulnerabilityResponse) {
	fmt.Printf("CVE: %s\n", v.CveID)
	fmt.Printf("Severity: %s\n", v.Severity)
	fmt.Printf("Package: %s (current: %s)\n", v.PackageName, v.CurrentVersion)
	if v.FixedVersion != nil {
		fmt.Printf("Fixed version: %s\n", *v.FixedVersion)
	}
	fmt.Printf("Fixable: %t\n", v.Fixable)
	fmt.Printf("Allowlisted: %t\n", v.IsAllowlisted)
	if v.TargetPath != "" {
		fmt.Printf("Target path: %s\n", v.TargetPath)
	}
	if v.Cvss.Preferred != nil && v.Cvss.Preferred.Score != nil {
		fmt.Printf("CVSS preferred score: %.1f\n", *v.Cvss.Preferred.Score)
	}
	fmt.Println("-------------------------------------------")
}
