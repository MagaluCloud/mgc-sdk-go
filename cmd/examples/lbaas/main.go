package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/lbaas"
)

// Configuration constants - modify these for your environment
const (
	// Replace with your actual VPC ID
	ExampleVPCID = "9dd2d30e-565d-42ce-a0a3-f2de1c473fed"
	// Replace with your actual subnet pool ID (optional)
	ExampleSubnetPoolID = "subnet-pool-12345678-1234-1234-1234-123456789012"
	// Replace with your actual public IP ID (optional for external LBs)
	ExamplePublicIPID = "public-ip-12345678-1234-1234-1234-123456789012"
)

// Global variables to store created resources for cleanup
var (
	createdBackends []string
	createdCerts    []string
	createdHCs      []string
	createdACLs     []string
)

func main() {
	fmt.Println("=== Magalu Cloud Load Balancer as a Service (LBaaS) SDK Examples ===")
	fmt.Println()

	// Initialize the SDK client
	client, err := initializeClient()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Create context for all operations
	ctx := context.Background()

	// Run examples in order
	fmt.Println("1. Creating a comprehensive Load Balancer...")
	lbID := runCreateLoadBalancerExample(ctx, client)
	if lbID == "" {
		log.Fatal("Failed to create load balancer, stopping examples")
	}

	fmt.Println("\n2. Listing all Load Balancers...")
	runListLoadBalancersExample(ctx, client)

	fmt.Println("\n3. Getting Load Balancer details...")
	runGetLoadBalancerExample(ctx, client, lbID)

	fmt.Println("\n4. Updating Load Balancer...")
	runUpdateLoadBalancerExample(ctx, client, lbID)

	fmt.Println("\n5. Managing Backends...")
	runManageBackendsExample(ctx, client, lbID)

	fmt.Println("\n6. Managing Listeners...")
	runManageListenersExample(ctx, client, lbID)

	fmt.Println("\n7. Managing Health Checks...")
	runManageHealthChecksExample(ctx, client, lbID)

	fmt.Println("\n8. Managing TLS Certificates...")
	runManageCertificatesExample(ctx, client, lbID)

	fmt.Println("\n9. Managing Access Control Lists (ACLs)...")
	runManageACLsExample(ctx, client, lbID)

	fmt.Println("\n10. Cleaning up resources...")
	runCleanupExample(ctx, client, lbID)

	fmt.Println("\n=== All examples completed successfully! ===")
}

// initializeClient sets up the MGC SDK client with proper configuration
func initializeClient() (*client.CoreClient, error) {
	fmt.Println("Initializing Magalu Cloud SDK client...")

	// Check for required environment variables
	apiKey := os.Getenv("MGC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MGC_API_KEY environment variable is required")
	}

	// Create the core client with configuration
	coreClient := client.NewMgcClient(
		client.WithAPIKey(apiKey),
		client.WithBaseURL(client.BrNe1),
	)

	return coreClient, nil
}

// waitForLoadBalancerStatus waits for a load balancer to reach a specific status
func waitForLoadBalancerStatus(ctx context.Context, client *client.CoreClient, lbID string, targetStatus lbaas.LoadBalancerStatus, timeout time.Duration) error {
	fmt.Printf("Waiting for load balancer %s to reach status: %s\n", lbID, targetStatus)

	lbService := lbaas.New(client).NetworkLoadBalancers()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		lb, err := lbService.Get(ctx, lbID)
		if err != nil {
			return fmt.Errorf("error checking load balancer status: %w", err)
		}

		currentStatus := lbaas.LoadBalancerStatus(lb.Status)
		fmt.Printf("Current status: %s\n", currentStatus)

		if currentStatus == targetStatus {
			fmt.Printf("✓ Load balancer reached target status: %s\n", targetStatus)
			return nil
		}

		if currentStatus == lbaas.LoadBalancerStatusFailed {
			return fmt.Errorf("load balancer creation failed")
		}

		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("timeout waiting for load balancer to reach status %s", targetStatus)
}

// runCreateLoadBalancerExample demonstrates creating a comprehensive load balancer
func runCreateLoadBalancerExample(ctx context.Context, client *client.CoreClient) string {
	fmt.Println("Creating a load balancer with backends, listeners, health checks, and certificates...")

	lbService := lbaas.New(client).NetworkLoadBalancers()

	// Generate a self-signed certificate for HTTPS listener
	cert, key, err := generateSelfSignedCertificate()
	if err != nil {
		log.Printf("Failed to generate certificate: %v", err)
		return ""
	}

	// Prepare the load balancer configuration
	createRequest := lbaas.CreateNetworkLoadBalancerRequest{
		Name:        "example-web-lb",
		Description: stringPtr("Example web application load balancer with HTTPS support"),
		Visibility:  lbaas.LoadBalancerVisibilityExternal,
		VPCID:       ExampleVPCID,

		// Optional: specify subnet pool and public IP
		// SubnetPoolID: &ExampleSubnetPoolID,
		// PublicIPID:   &ExamplePublicIPID,

		// Define backends (server pools)
		Backends: []lbaas.CreateNetworkBackendRequest{
			{
				Name:                                "web-servers",
				Description:                         stringPtr("Backend pool for web servers"),
				BalanceAlgorithm:                    lbaas.BackendBalanceAlgorithmRoundRobin,
				TargetsType:                         lbaas.BackendTypeInstance,
				PanicThreshold:                      intPtr(50), // Panic when 50% of targets are unhealthy
				CloseConnectionsOnHostHealthFailure: boolPtr(true),
				// Health check will be linked later
			},
			{
				Name:                                "api-servers",
				Description:                         stringPtr("Backend pool for API servers"),
				BalanceAlgorithm:                    lbaas.BackendBalanceAlgorithmRoundRobin,
				TargetsType:                         lbaas.BackendTypeInstance,
				PanicThreshold:                      intPtr(30),
				CloseConnectionsOnHostHealthFailure: boolPtr(false),
			},
		},

		// Define health checks
		HealthChecks: []lbaas.CreateNetworkHealthCheckRequest{
			{
				Name:                    "web-health-check",
				Description:             stringPtr("HTTP health check for web servers"),
				Protocol:                lbaas.HealthCheckProtocolHTTP,
				Port:                    80,
				Path:                    stringPtr("/health"),
				HealthyStatusCode:       intPtr(200),
				IntervalSeconds:         intPtr(30),
				TimeoutSeconds:          intPtr(5),
				InitialDelaySeconds:     intPtr(10),
				HealthyThresholdCount:   intPtr(3),
				UnhealthyThresholdCount: intPtr(3),
			},
			{
				Name:                    "api-health-check",
				Description:             stringPtr("HTTP health check for API servers"),
				Protocol:                lbaas.HealthCheckProtocolHTTP,
				Port:                    8080,
				Path:                    stringPtr("/api/health"),
				HealthyStatusCode:       intPtr(200),
				IntervalSeconds:         intPtr(15),
				TimeoutSeconds:          intPtr(3),
				InitialDelaySeconds:     intPtr(5),
				HealthyThresholdCount:   intPtr(2),
				UnhealthyThresholdCount: intPtr(2),
			},
		},

		// Define TLS certificates for HTTPS
		TLSCertificates: []lbaas.CreateNetworkCertificateRequest{
			{
				Name:        "web-ssl-cert",
				Description: stringPtr("SSL certificate for web application"),
				Certificate: cert,
				PrivateKey:  key,
			},
		},

		// Define listeners (entry points)
		Listeners: []lbaas.NetworkListenerRequest{
			{
				BackendName: "web-servers", // Links to the backend defined above
				Name:        "http-listener",
				Description: stringPtr("HTTP listener for web traffic"),
				Protocol:    lbaas.ListenerProtocolTCP,
				Port:        80,
			},
			{
				BackendName:        "web-servers",
				Name:               "https-listener",
				Description:        stringPtr("HTTPS listener for secure web traffic"),
				Protocol:           lbaas.ListenerProtocolTLS,
				Port:               443,
				TLSCertificateName: stringPtr("web-ssl-cert"), // Links to certificate by name
			},
			{
				BackendName: "api-servers",
				Name:        "api-listener",
				Description: stringPtr("API listener"),
				Protocol:    lbaas.ListenerProtocolTCP,
				Port:        8080,
			},
		},

		// Define Access Control Lists for security
		ACLs: []lbaas.CreateNetworkACLRequest{
			{
				Name:           stringPtr("allow-web-traffic"),
				Ethertype:      lbaas.AclEtherTypeIPv4,
				Protocol:       lbaas.AclProtocolTCP,
				Action:         lbaas.AclActionTypeAllow,
				RemoteIPPrefix: "0.0.0.0/0", // Allow from anywhere (adjust as needed)
			},
			{
				Name:           stringPtr("allow-api-internal"),
				Ethertype:      lbaas.AclEtherTypeIPv4,
				Protocol:       lbaas.AclProtocolTCP,
				Action:         lbaas.AclActionTypeAllow,
				RemoteIPPrefix: "10.0.0.0/8", // Allow only from internal networks
			},
		},
	}

	// Create the load balancer
	fmt.Println("Sending create request...")
	lbID, err := lbService.Create(ctx, createRequest)
	if err != nil {
		log.Printf("Failed to create load balancer: %v", err)
		return ""
	}

	fmt.Printf("✓ Load balancer created with ID: %s\n", lbID)

	// Wait for the load balancer to become ready
	if err := waitForLoadBalancerStatus(ctx, client, lbID, lbaas.LoadBalancerStatusRunning, 10*time.Minute); err != nil {
		log.Printf("Load balancer creation failed or timed out: %v", err)
		return ""
	}

	fmt.Printf("✓ Load balancer is now running and ready to serve traffic!\n")
	return lbID
}

// runListLoadBalancersExample demonstrates listing load balancers with pagination
func runListLoadBalancersExample(ctx context.Context, client *client.CoreClient) {
	fmt.Println("Listing load balancers with pagination...")

	lbService := lbaas.New(client).NetworkLoadBalancers()

	// Example 1: List with pagination options
	listOptions := lbaas.ListNetworkLoadBalancerRequest{
		Limit:  intPtr(10),                   // Get up to 10 results per page
		Offset: intPtr(0),                    // Start from the beginning
		Sort:   stringPtr("created_at:desc"), // Sort by creation date, newest first
	}

	paginatedResp, err := lbService.List(ctx, listOptions)
	if err != nil {
		log.Printf("Failed to list load balancers: %v", err)
		return
	}

	fmt.Printf("Paginated results - Found %d load balancer(s) on this page (Total: %d):\n",
		len(paginatedResp.Results), paginatedResp.Meta.Page.Total)
	for i, lb := range paginatedResp.Results {
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, lb.Name, lb.ID)
		fmt.Printf("     Status: %s, Visibility: %s\n", lb.Status, lb.Visibility)
		fmt.Printf("     Created: %s\n", lb.CreatedAt)
		fmt.Printf("     Listeners: %d, Backends: %d, Health Checks: %d\n",
			len(lb.Listeners), len(lb.Backends), len(lb.HealthChecks))
		if lb.PublicIP != nil {
			fmt.Printf("     Public IP: %s\n", *lb.PublicIP.IPAddress)
		}
		fmt.Println()
	}

	// Example 2: List all load balancers across all pages
	fmt.Println("\nFetching ALL load balancers (across all pages)...")
	allLoadBalancers, err := lbService.ListAll(ctx)
	if err != nil {
		log.Printf("Failed to list all load balancers: %v", err)
		return
	}

	fmt.Printf("✓ Retrieved %d total load balancer(s) across all pages\n", len(allLoadBalancers))
}

// runGetLoadBalancerExample demonstrates getting detailed information about a load balancer
func runGetLoadBalancerExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Getting detailed information for load balancer: %s\n", lbID)

	lbService := lbaas.New(client).NetworkLoadBalancers()

	lb, err := lbService.Get(ctx, lbID)
	if err != nil {
		log.Printf("Failed to get load balancer: %v", err)
		return
	}

	// Display comprehensive information
	fmt.Printf("=== Load Balancer Details ===\n")
	fmt.Printf("ID: %s\n", lb.ID)
	fmt.Printf("Name: %s\n", lb.Name)
	if lb.Description != nil {
		fmt.Printf("Description: %s\n", *lb.Description)
	}
	fmt.Printf("Type: %s\n", lb.Type)
	fmt.Printf("Visibility: %s\n", lb.Visibility)
	fmt.Printf("Status: %s\n", lb.Status)
	fmt.Printf("VPC ID: %s\n", lb.VPCID)
	if lb.PublicIP != nil {
		fmt.Printf("     Public IP: %s\n", *lb.PublicIP.IPAddress)
	}
	fmt.Printf("Created: %s\n", lb.CreatedAt)
	fmt.Printf("Updated: %s\n", lb.UpdatedAt)

	// Display listeners
	fmt.Printf("\nListeners (%d):\n", len(lb.Listeners))
	for _, listener := range lb.Listeners {
		fmt.Printf("  - %s (Port: %d, Protocol: %s)\n", listener.Name, listener.Port, listener.Protocol)
		if listener.TLSCertificateID != nil {
			fmt.Printf("    TLS Certificate: %s\n", *listener.TLSCertificateID)
		}
	}

	// Display backends
	fmt.Printf("\nBackends (%d):\n", len(lb.Backends))
	for _, backend := range lb.Backends {
		fmt.Printf("  - %s (Algorithm: %s, Targets: %d)\n",
			backend.Name, backend.BalanceAlgorithm, len(backend.Targets))
		if backend.HealthCheckID != nil {
			fmt.Printf("    Health Check: %s\n", *backend.HealthCheckID)
		}
	}

	// Display health checks
	fmt.Printf("\nHealth Checks (%d):\n", len(lb.HealthChecks))
	for _, hc := range lb.HealthChecks {
		fmt.Printf("  - %s (Protocol: %s, Port: %d)\n", hc.Name, hc.Protocol, hc.Port)
		if hc.Path != nil {
			fmt.Printf("    Path: %s\n", *hc.Path)
		}
	}

	// Display certificates
	fmt.Printf("\nTLS Certificates (%d):\n", len(lb.TLSCertificates))
	for _, cert := range lb.TLSCertificates {
		fmt.Printf("  - %s\n", cert.Name)
		if cert.ExpirationDate != nil {
			fmt.Printf("    Expires: %s\n", *cert.ExpirationDate)
		}
	}

	// Display ACLs
	fmt.Printf("\nACLs (%d):\n", len(lb.ACLs))
	for _, acl := range lb.ACLs {
		name := "unnamed"
		if acl.Name != nil {
			name = *acl.Name
		}
		fmt.Printf("  - %s (Action: %s, Protocol: %s, CIDR: %s)\n",
			name, acl.Action, acl.Protocol, acl.RemoteIPPrefix)
	}

	// Display public IPs
	if lb.PublicIP != nil {
		fmt.Printf("     Public IP: %s\n", *lb.PublicIP.IPAddress)
	}
}

// runUpdateLoadBalancerExample demonstrates updating a load balancer
func runUpdateLoadBalancerExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Updating load balancer: %s\n", lbID)

	lbService := lbaas.New(client).NetworkLoadBalancers()

	// Update the load balancer name and description
	updateRequest := lbaas.UpdateNetworkLoadBalancerRequest{
		Name:        stringPtr("updated-web-lb"),
		Description: stringPtr("Updated description with enhanced features"),
	}

	updatedID, err := lbService.Update(ctx, lbID, updateRequest)
	if err != nil {
		log.Printf("Failed to update load balancer: %v", err)
		return
	}

	fmt.Printf("✓ Load balancer updated successfully (ID: %s)\n", updatedID)

	// Wait for update to complete
	if err := waitForLoadBalancerStatus(ctx, client, lbID, lbaas.LoadBalancerStatusRunning, 5*time.Minute); err != nil {
		log.Printf("Update failed or timed out: %v", err)
		return
	}

	fmt.Println("✓ Load balancer update completed")
}

// runManageBackendsExample demonstrates backend management operations
func runManageBackendsExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Managing backends for load balancer: %s\n", lbID)

	backendService := lbaas.New(client).NetworkBackends()

	// Create a new backend
	fmt.Println("Creating a new backend...")
	createBackendReq := lbaas.CreateBackendRequest{
		Name:                                "new-backend-pool",
		Description:                         stringPtr("Additional backend pool for scaling"),
		BalanceAlgorithm:                    lbaas.BackendBalanceAlgorithmRoundRobin,
		TargetsType:                         lbaas.BackendTypeInstance,
		PanicThreshold:                      floatPtr(40.0),
		CloseConnectionsOnHostHealthFailure: boolPtr(true),
	}

	backendID, err := backendService.Create(ctx, lbID, createBackendReq)
	if err != nil {
		log.Printf("Failed to create backend: %v", err)
		return
	}
	fmt.Printf("✓ Backend created with ID: %s\n", backendID)
	createdBackends = append(createdBackends, backendID)

	// List all backends (using ListAll to get all pages)
	fmt.Println("\nListing all backends...")
	backends, err := backendService.ListAll(ctx, lbID)
	if err != nil {
		log.Printf("Failed to list backends: %v", err)
		return
	}

	fmt.Printf("Found %d backend(s):\n", len(backends))
	for _, backend := range backends {
		fmt.Printf("  - %s (ID: %s, Targets: %d)\n",
			backend.Name, backend.ID, len(backend.Targets))
	}

	// Get detailed backend information
	if len(backends) > 0 {
		fmt.Printf("\nGetting details for backend: %s\n", backends[0].ID)
		backend, err := backendService.Get(ctx, lbID, backends[0].ID)
		if err != nil {
			log.Printf("Failed to get backend details: %v", err)
			return
		}

		fmt.Printf("Backend: %s\n", backend.Name)
		fmt.Printf("  Algorithm: %s\n", backend.BalanceAlgorithm)
		fmt.Printf("  Targets Type: %s\n", backend.TargetsType)
		if backend.PanicThreshold != nil {
			fmt.Printf("  Panic Threshold: %d%%\n", *backend.PanicThreshold)
		}
		fmt.Printf("  Close Connections on Health Failure: %t\n", *backend.CloseConnectionsOnHostHealthFailure)
	}

	// Update the backend
	fmt.Printf("\nUpdating backend: %s\n", backendID)
	updateBackendReq := lbaas.UpdateNetworkBackendRequest{
		PanicThreshold:                      intPtr(60),
		CloseConnectionsOnHostHealthFailure: boolPtr(false),
	}

	_, err = backendService.Update(ctx, lbID, backendID, updateBackendReq)
	if err != nil {
		log.Printf("Failed to update backend: %v", err)
		return
	}
	fmt.Println("✓ Backend updated successfully")
}

// runManageListenersExample demonstrates listener management operations
func runManageListenersExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Managing listeners for load balancer: %s\n", lbID)

	listenerService := lbaas.New(client).NetworkListeners()
	backendService := lbaas.New(client).NetworkBackends()

	// First, get available backends to link to
	backends, err := backendService.ListAll(ctx, lbID)
	if err != nil || len(backends) == 0 {
		log.Printf("No backends available for listener creation")
		return
	}

	// Create a new listener
	fmt.Println("Creating a new listener...")
	createListenerReq := lbaas.CreateNetworkListenerRequest{
		Name:        "admin-listener",
		Description: stringPtr("Administrative interface listener"),
		Protocol:    lbaas.ListenerProtocolTCP,
		Port:        9090,
	}

	listener, err := listenerService.Create(ctx, lbID, backends[0].ID, createListenerReq)
	if err != nil {
		log.Printf("Failed to create listener: %v", err)
		return
	}
	fmt.Printf("✓ Listener created with ID: %s\n", listener.ID)

	// List all listeners (using ListAll to get all pages)
	fmt.Println("\nListing all listeners...")
	listeners, err := listenerService.ListAll(ctx, lbID)
	if err != nil {
		log.Printf("Failed to list listeners: %v", err)
		return
	}

	fmt.Printf("Found %d listener(s):\n", len(listeners))
	for _, l := range listeners {
		fmt.Printf("  - %s (Port: %d, Protocol: %s, Backend: %s)\n",
			l.Name, l.Port, l.Protocol, l.BackendID)
	}

	// Get detailed listener information
	fmt.Printf("\nGetting details for listener: %s\n", listener.ID)
	detailedListener, err := listenerService.Get(ctx, lbID, listener.ID)
	if err != nil {
		log.Printf("Failed to get listener details: %v", err)
		return
	}

	fmt.Printf("Listener: %s\n", detailedListener.Name)
	fmt.Printf("  Protocol: %s\n", detailedListener.Protocol)
	fmt.Printf("  Port: %d\n", detailedListener.Port)
	fmt.Printf("  Backend ID: %s\n", detailedListener.BackendID)
	if detailedListener.TLSCertificateID != nil {
		fmt.Printf("  TLS Certificate: %s\n", *detailedListener.TLSCertificateID)
	}

	// Update the listener
	fmt.Printf("\nUpdating listener: %s\n", listener.ID)
	updateListenerReq := lbaas.UpdateNetworkListenerRequest{
		Name: stringPtr("updated-admin-listener"),
	}

	err = listenerService.Update(ctx, lbID, listener.ID, updateListenerReq)
	if err != nil {
		log.Printf("Failed to update listener: %v", err)
		return
	}
	fmt.Println("✓ Listener updated successfully")
}

// runManageHealthChecksExample demonstrates health check management
func runManageHealthChecksExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Managing health checks for load balancer: %s\n", lbID)

	hcService := lbaas.New(client).NetworkHealthChecks()

	// Create a new health check
	fmt.Println("Creating a new health check...")
	createHCReq := lbaas.CreateNetworkHealthCheckRequest{
		Name:                    "tcp-health-check",
		Description:             stringPtr("Simple TCP health check"),
		Protocol:                lbaas.HealthCheckProtocolTCP,
		Port:                    3306, // MySQL port example
		IntervalSeconds:         intPtr(20),
		TimeoutSeconds:          intPtr(5),
		InitialDelaySeconds:     intPtr(15),
		HealthyThresholdCount:   intPtr(2),
		UnhealthyThresholdCount: intPtr(3),
	}

	hc, err := hcService.Create(ctx, lbID, createHCReq)
	if err != nil {
		log.Printf("Failed to create health check: %v", err)
		return
	}
	fmt.Printf("✓ Health check created with ID: %s\n", hc.ID)
	createdHCs = append(createdHCs, hc.ID)

	// List all health checks (using ListAll to get all pages)
	fmt.Println("\nListing all health checks...")
	healthChecks, err := hcService.ListAll(ctx, lbID)
	if err != nil {
		log.Printf("Failed to list health checks: %v", err)
		return
	}

	fmt.Printf("Found %d health check(s):\n", len(healthChecks))
	for _, hc := range healthChecks {
		fmt.Printf("  - %s (Protocol: %s, Port: %d)\n", hc.Name, hc.Protocol, hc.Port)
		fmt.Printf("    Interval: %ds, Timeout: %ds\n", hc.IntervalSeconds, hc.TimeoutSeconds)
		fmt.Printf("    Healthy Threshold: %d, Unhealthy Threshold: %d\n",
			hc.HealthyThresholdCount, hc.UnhealthyThresholdCount)
	}

	// Get detailed health check information
	fmt.Printf("\nGetting details for health check: %s\n", hc.ID)
	detailedHC, err := hcService.Get(ctx, lbID, hc.ID)
	if err != nil {
		log.Printf("Failed to get health check details: %v", err)
		return
	}

	fmt.Printf("Health Check: %s\n", detailedHC.Name)
	fmt.Printf("  Protocol: %s\n", detailedHC.Protocol)
	fmt.Printf("  Port: %d\n", detailedHC.Port)
	if detailedHC.Path != nil {
		fmt.Printf("  Path: %s\n", *detailedHC.Path)
	}
	fmt.Printf("  Status Code: %d\n", detailedHC.HealthyStatusCode)

	// Update the health check
	fmt.Printf("\nUpdating health check: %s\n", hc.ID)
	updateHCReq := lbaas.UpdateNetworkHealthCheckRequest{
		Protocol:        lbaas.HealthCheckProtocolHTTP,
		Port:            8080,
		Path:            stringPtr("/status"),
		IntervalSeconds: intPtr(15),
	}

	err = hcService.Update(ctx, lbID, hc.ID, updateHCReq)
	if err != nil {
		log.Printf("Failed to update health check: %v", err)
		return
	}
	fmt.Println("✓ Health check updated successfully")
}

// runManageCertificatesExample demonstrates TLS certificate management
func runManageCertificatesExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Managing TLS certificates for load balancer: %s\n", lbID)

	certService := lbaas.New(client).NetworkCertificates()

	// Generate a new certificate
	fmt.Println("Generating a new TLS certificate...")
	cert, key, err := generateSelfSignedCertificate()
	if err != nil {
		log.Printf("Failed to generate certificate: %v", err)
		return
	}

	// Create a new certificate
	fmt.Println("Creating a new TLS certificate...")
	createCertReq := lbaas.CreateNetworkCertificateRequest{
		Name:        "api-ssl-cert",
		Description: stringPtr("SSL certificate for API endpoints"),
		Certificate: cert,
		PrivateKey:  key,
	}

	certResp, err := certService.Create(ctx, lbID, createCertReq)
	if err != nil {
		log.Printf("Failed to create certificate: %v", err)
		return
	}
	fmt.Printf("✓ Certificate created with ID: %s\n", certResp.ID)
	createdCerts = append(createdCerts, certResp.ID)

	// List all certificates (using ListAll to get all pages)
	fmt.Println("\nListing all certificates...")
	certificates, err := certService.ListAll(ctx, lbID)
	if err != nil {
		log.Printf("Failed to list certificates: %v", err)
		return
	}

	fmt.Printf("Found %d certificate(s):\n", len(certificates))
	for _, cert := range certificates {
		fmt.Printf("  - %s (ID: %s)\n", cert.Name, cert.ID)
		if cert.ExpirationDate != nil {
			fmt.Printf("    Expires: %s\n", *cert.ExpirationDate)
		}
		fmt.Printf("    Created: %s\n", cert.CreatedAt)
	}

	// Get detailed certificate information
	fmt.Printf("\nGetting details for certificate: %s\n", certResp.ID)
	detailedCert, err := certService.Get(ctx, lbID, certResp.ID)
	if err != nil {
		log.Printf("Failed to get certificate details: %v", err)
		return
	}

	fmt.Printf("Certificate: %s\n", detailedCert.Name)
	if detailedCert.Description != nil {
		fmt.Printf("  Description: %s\n", *detailedCert.Description)
	}
	if detailedCert.ExpirationDate != nil {
		fmt.Printf("  Expires: %s\n", *detailedCert.ExpirationDate)
	}

	// Update the certificate (replace with new cert/key pair)
	fmt.Println("\nGenerating and updating certificate...")
	newCert, newKey, err := generateSelfSignedCertificate()
	if err != nil {
		log.Printf("Failed to generate new certificate: %v", err)
		return
	}

	updateCertReq := lbaas.UpdateNetworkCertificateRequest{
		Certificate: newCert,
		PrivateKey:  newKey,
	}

	err = certService.Update(ctx, lbID, certResp.ID, updateCertReq)
	if err != nil {
		log.Printf("Failed to update certificate: %v", err)
		return
	}
	fmt.Println("✓ Certificate updated successfully")
}

// runManageACLsExample demonstrates Access Control List management
func runManageACLsExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Managing ACLs for load balancer: %s\n", lbID)

	aclService := lbaas.New(client).NetworkACLs()

	// Create individual ACL rules
	fmt.Println("Creating ACL rules...")

	// Allow HTTP from anywhere
	httpACLID, err := aclService.Create(ctx, lbID, lbaas.CreateNetworkACLRequest{
		Name:           stringPtr("allow-http-global"),
		Ethertype:      lbaas.AclEtherTypeIPv4,
		Protocol:       lbaas.AclProtocolTCP,
		Action:         lbaas.AclActionTypeAllow,
		RemoteIPPrefix: "0.0.0.0/0",
	})
	if err != nil {
		log.Printf("Failed to create HTTP ACL: %v", err)
		return
	}
	fmt.Printf("✓ HTTP ACL created with ID: %s\n", httpACLID)
	createdACLs = append(createdACLs, httpACLID)

	// Allow HTTPS from anywhere
	httpsACLID, err := aclService.Create(ctx, lbID, lbaas.CreateNetworkACLRequest{
		Name:           stringPtr("allow-https-global"),
		Ethertype:      lbaas.AclEtherTypeIPv4,
		Protocol:       lbaas.AclProtocolTLS,
		Action:         lbaas.AclActionTypeAllow,
		RemoteIPPrefix: "0.0.0.0/0",
	})
	if err != nil {
		log.Printf("Failed to create HTTPS ACL: %v", err)
		return
	}
	fmt.Printf("✓ HTTPS ACL created with ID: %s\n", httpsACLID)
	createdACLs = append(createdACLs, httpsACLID)

	// Deny access from a specific IP range
	denyACLID, err := aclService.Create(ctx, lbID, lbaas.CreateNetworkACLRequest{
		Name:           stringPtr("deny-suspicious-range"),
		Ethertype:      lbaas.AclEtherTypeIPv4,
		Protocol:       lbaas.AclProtocolTCP,
		Action:         lbaas.AclActionTypeDeny,
		RemoteIPPrefix: "192.168.100.0/24", // Example suspicious range
	})
	if err != nil {
		log.Printf("Failed to create deny ACL: %v", err)
		return
	}
	fmt.Printf("✓ Deny ACL created with ID: %s\n", denyACLID)
	createdACLs = append(createdACLs, denyACLID)

	// Example: Replace all ACLs with a new set
	fmt.Println("\nReplacing all ACLs with a new comprehensive set...")
	replaceACLsReq := lbaas.UpdateNetworkACLRequest{
		Acls: []lbaas.CreateNetworkACLRequest{
			{
				Name:           stringPtr("allow-web-traffic"),
				Ethertype:      lbaas.AclEtherTypeIPv4,
				Protocol:       lbaas.AclProtocolTCP,
				Action:         lbaas.AclActionTypeAllow,
				RemoteIPPrefix: "0.0.0.0/0",
			},
			{
				Name:           stringPtr("allow-office-network"),
				Ethertype:      lbaas.AclEtherTypeIPv4,
				Protocol:       lbaas.AclProtocolTCP,
				Action:         lbaas.AclActionTypeAllow,
				RemoteIPPrefix: "10.0.0.0/8",
			},
			{
				Name:           stringPtr("deny-blacklisted-range"),
				Ethertype:      lbaas.AclEtherTypeIPv4,
				Protocol:       lbaas.AclProtocolTCP,
				Action:         lbaas.AclActionTypeDeny,
				RemoteIPPrefix: "203.0.113.0/24", // RFC 5737 test range
			},
		},
	}

	err = aclService.Replace(ctx, lbID, replaceACLsReq)
	if err != nil {
		log.Printf("Failed to replace ACLs: %v", err)
		return
	}
	fmt.Println("✓ ACLs replaced successfully")

	fmt.Println("\nACL management completed. The load balancer now has updated access control rules.")
}

// runCleanupExample demonstrates proper resource cleanup
func runCleanupExample(ctx context.Context, client *client.CoreClient, lbID string) {
	fmt.Printf("Cleaning up resources for load balancer: %s\n", lbID)

	// Note: In a real application, you might want to clean up individual resources
	// before deleting the load balancer, but since we created them as part of the LB,
	// they'll be cleaned up automatically when we delete the LB.

	lbService := lbaas.New(client).NetworkLoadBalancers()

	// Delete the load balancer (this will cascade delete most sub-resources)
	fmt.Println("Deleting load balancer...")
	deleteRequest := lbaas.DeleteNetworkLoadBalancerRequest{
		DeletePublicIP: boolPtr(true), // Also delete the associated public IP
	}

	err := lbService.Delete(ctx, lbID, deleteRequest)
	if err != nil {
		log.Printf("Failed to delete load balancer: %v", err)
		return
	}

	fmt.Println("✓ Load balancer deletion initiated")

	// Wait for deletion to complete
	fmt.Println("Waiting for deletion to complete...")
	deadline := time.Now().Add(10 * time.Minute)
	for time.Now().Before(deadline) {
		_, err := lbService.Get(ctx, lbID)
		if err != nil {
			// If we get an error (likely 404), the LB is deleted
			if strings.Contains(err.Error(), "404") {
				fmt.Println("✓ Load balancer successfully deleted")
				return
			}
		}
		time.Sleep(10 * time.Second)
	}

	fmt.Println("⚠ Deletion timeout reached - load balancer may still be deleting")
}

// Helper function to generate a self-signed certificate for testing
func generateSelfSignedCertificate() (string, string, error) {
	// Generate a private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Example Org"},
			Country:       []string{"BR"},
			Province:      []string{"SP"},
			Locality:      []string{"São Paulo"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: nil,
		DNSNames:    []string{"example.com", "*.example.com"},
	}

	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key to PEM
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})

	// Encode to base64 as required by the API
	certBase64 := base64.StdEncoding.EncodeToString(certPEM)
	keyBase64 := base64.StdEncoding.EncodeToString(privateKeyPEM)

	return certBase64, keyBase64, nil
}

// Helper functions for pointer conversions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}
