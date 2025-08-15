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
	"net"
	"os"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/lbaas"
)

// main is an orchestrator showcasing the recommended execution order of the examples.
// It demonstrates: listing LBs, creating a new LB, waiting for status changes,
// updating the LB, managing sub-resources, and finally deleting the LB.
//
// Tip:
// - Replace vpcID with a valid VPC ID in your account.
// - You can comment out sections you don't need.
// - For long-running CI samples, consider skipping creation/deletion and only run read-only samples.
func main() {
	ExampleListLoadBalancers()

	vpcID := "b44e3fe0-b609-4906-81a3-256e0cf68cfb"
	lbID := ExampleCreateLoadBalancer(vpcID)
	time.Sleep(10 * time.Second)

	// Examples of operations with an existing Load Balancer
	// Replace with a real ID to test
	// lbID := "" // comment and uncomment to run the examples

	if lbID != "" {
		waitCreateLoadBalancer(lbID, lbaas.LoadBalancerStatusCreating)
		ExampleGetLoadBalancer(lbID)

		ExampleUpdateLoadBalancer(lbID)
		waitCreateLoadBalancer(lbID, lbaas.LoadBalancerStatusUpdating)

		ExampleManageBackends(lbID)
		ExampleManageListeners(lbID)
		ExampleManageHealthChecks(lbID)
		ExampleManageCertificates(lbID)
		ExampleManageACLs(lbID)
		ExampleGetLoadBalancer(lbID)

		err := deleteWithRetry(lbID)
		if err != nil {
			log.Fatal("Erro ao deletar Load Balancer:", err)
		}
	}
}

// deleteWithRetry attempts to delete a load balancer and retries if the platform responds
// with a 409 Conflict (usually indicating a concurrent operation in progress).
//
// Recommended usage:
// - Use this helper when tearing down test infrastructure that may still be reconciling.
// - In production, consider exponential backoff and a max retry budget.
func deleteWithRetry(lbID string) error {
	err := ExampleDeleteLoadBalancer(lbID)
	if err != nil {
		if httpErr, ok := err.(*client.HTTPError); ok {
			if httpErr.StatusCode == 409 {
				time.Sleep(10 * time.Second)
				return deleteWithRetry(lbID)
			}
		}
		return err
	}

	return nil
}

// waitCreateLoadBalancer polls a load balancer until it is no longer in the provided status.
// Common use cases include waiting for "creating" to finish or "updating" to complete.
//
// Parameters:
// - lbID: the load balancer ID returned by Create
// - status: the transitional status to wait out (e.g., lbaas.LoadBalancerStatusCreating)
//
// Caution:
//   - This is a simple polling loop for demonstration. In production, set timeouts, backoff,
//     and maximum attempts to avoid infinite waits.
func waitCreateLoadBalancer(lbID string, status lbaas.LoadBalancerStatus) {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Get Load Balancer details
	for {
		lb, err := lbaasClient.NetworkLoadBalancers().Get(ctx, lbID)
		if err != nil {
			log.Fatal("Erro ao obter Load Balancer:", err)
		}
		if lb.Status != string(status) {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

// ExampleCreateLoadBalancer shows how to create a network load balancer with:
// - a health check
// - a backend (with raw IP targets)
// - a listener
// - an optional ACL
//
// Inputs:
// - vpcID: the VPC where the load balancer will be provisioned
//
// Returns:
// - The created Load Balancer ID on success
//
// After creation:
// - The platform continues provisioning asynchronously. Use waitCreateLoadBalancer to poll.
// - You can optionally assign a Public IP at creation or later.
func ExampleCreateLoadBalancer(vpcID string) string {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Configure health check
	healthCheck := lbaas.NetworkHealthCheckRequest{
		Name:                    "example-health-check",
		Description:             helpers.StrPtr("Health check para exemplo"),
		Protocol:                lbaas.HealthCheckProtocolHTTP,
		Path:                    helpers.StrPtr("/health"),
		Port:                    80,
		HealthyStatusCode:       helpers.IntPtr(200),
		IntervalSeconds:         helpers.IntPtr(30),
		TimeoutSeconds:          helpers.IntPtr(5),
		InitialDelaySeconds:     helpers.IntPtr(10),
		HealthyThresholdCount:   helpers.IntPtr(3),
		UnhealthyThresholdCount: helpers.IntPtr(3),
	}

	// Configure backend with raw targets (IPs)
	backend := lbaas.NetworkBackendRequest{
		Name:             "example-backend",
		Description:      helpers.StrPtr("Backend para exemplo"),
		BalanceAlgorithm: lbaas.BackendBalanceAlgorithmRoundRobin,
		TargetsType:      lbaas.BackendTypeRaw,
		HealthCheckName:  helpers.StrPtr("example-health-check"),
		Targets: &lbaas.TargetsRawOrInstancesRequest{
			TargetsRaw: []lbaas.NetworkBackendRawTargetRequest{
				{
					IPAddress: "192.168.1.10",
					Port:      80,
				},
				{
					IPAddress: "192.168.1.11",
					Port:      80,
				},
			},
		},
	}

	// Configure listener
	listener := lbaas.NetworkListenerRequest{
		Name:        "example-listener",
		Description: helpers.StrPtr("Listener HTTP para exemplo"),
		BackendName: "example-backend",
		Protocol:    lbaas.ListenerProtocolTCP,
		Port:        80,
	}

	// Configure ACL (optional)
	acl := lbaas.NetworkAclRequest{
		Name:           helpers.StrPtr("allow-all"),
		Ethertype:      lbaas.AclEtherTypeIPv4,
		Protocol:       lbaas.AclProtocolTCP,
		RemoteIPPrefix: "0.0.0.0/0",
		Action:         lbaas.AclActionTypeAllow,
	}

	// Create Load Balancer
	createReq := lbaas.CreateNetworkLoadBalancerRequest{
		Name:         "example-load-balancer-" + time.Now().Format("20060102-150405"),
		Description:  helpers.StrPtr("Load Balancer de exemplo criado via SDK"),
		Type:         helpers.StrPtr("proxy"),
		Visibility:   lbaas.LoadBalancerVisibilityExternal,
		Listeners:    []lbaas.NetworkListenerRequest{listener},
		Backends:     []lbaas.NetworkBackendRequest{backend},
		HealthChecks: []lbaas.NetworkHealthCheckRequest{healthCheck},
		ACLs:         []lbaas.NetworkAclRequest{acl},
		VPCID:        vpcID,
		// PublicIPID:   helpers.StrPtr("your-public-ip-id"), // Optional
	}

	id, err := lbaasClient.NetworkLoadBalancers().Create(ctx, createReq)
	if err != nil {
		log.Fatal("Erro ao criar Load Balancer:", err)
	}

	fmt.Printf("Load Balancer criado com sucesso! ID: %s\n", id)
	return id
}

// ExampleListLoadBalancers demonstrates listing load balancers with common pagination options.
// It prints summarized information for each result.
//
// Notes:
// - For large projects, always use pagination to control API and memory impact.
// - Sorting is supported on common fields like created_at.
func ExampleListLoadBalancers() {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// List Load Balancers with pagination
	listReq := lbaas.ListNetworkLoadBalancerRequest{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Sort:   helpers.StrPtr("created_at:desc"),
	}

	loadBalancers, err := lbaasClient.NetworkLoadBalancers().List(ctx, listReq)
	if err != nil {
		log.Fatal("Erro ao listar Load Balancers:", err)
	}

	fmt.Printf("Encontrados %d Load Balancers:\n", len(loadBalancers))
	for _, lb := range loadBalancers {
		fmt.Printf("- ID: %s\n", lb.ID)
		fmt.Printf("  Nome: %s\n", lb.Name)
		fmt.Printf("  Tipo: %s\n", lb.Type)
		fmt.Printf("  Visibilidade: %s\n", lb.Visibility)
		fmt.Printf("  Status: %s\n", lb.Status)
		fmt.Printf("  VPC ID: %s\n", lb.VPCID)
		if lb.IPAddress != nil {
			fmt.Printf("  IP Address: %s\n", *lb.IPAddress)
		}
		fmt.Printf("  Criado em: %s\n", lb.CreatedAt)
		fmt.Printf("  Listeners: %d\n", len(lb.Listeners))
		fmt.Printf("  Backends: %d\n", len(lb.Backends))
		fmt.Printf("  Health Checks: %d\n", len(lb.HealthChecks))
		fmt.Printf("  Certificados TLS: %d\n", len(lb.TLSCertificates))
		fmt.Printf("  ACLs: %d\n", len(lb.ACLs))
		fmt.Println("  ---")
	}
}

// ExampleGetLoadBalancer fetches a single load balancer by ID and prints key details,
// including public IPs, listeners, and backends.
//
// Use cases:
// - Inspect current state during workflows
// - Gather identifiers for sub-resource operations
func ExampleGetLoadBalancer(lbID string) {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Get Load Balancer details
	lb, err := lbaasClient.NetworkLoadBalancers().Get(ctx, lbID)
	if err != nil {
		log.Fatal("Erro ao obter Load Balancer:", err)
	}

	fmt.Printf("Detalhes do Load Balancer %s:\n", lb.ID)
	fmt.Printf("  Nome: %s\n", lb.Name)
	if lb.Description != nil {
		fmt.Printf("  Descrição: %s\n", *lb.Description)
	}
	fmt.Printf("  Tipo: %s\n", lb.Type)
	fmt.Printf("  Visibilidade: %s\n", lb.Visibility)
	fmt.Printf("  Status: %s\n", lb.Status)
	fmt.Printf("  VPC ID: %s\n", lb.VPCID)
	if lb.SubnetPoolID != nil {
		fmt.Printf("  Subnet Pool ID: %s\n", *lb.SubnetPoolID)
	}
	if lb.IPAddress != nil {
		fmt.Printf("  IP Address: %s\n", *lb.IPAddress)
	}
	fmt.Printf("  Criado em: %s\n", lb.CreatedAt)
	fmt.Printf("  Atualizado em: %s\n", lb.UpdatedAt)

	// Show public IPs
	if len(lb.PublicIPs) > 0 {
		fmt.Println("  IPs Públicos:")
		for _, ip := range lb.PublicIPs {
			fmt.Printf("    - ID: %s, IP: %s, External ID: %s\n",
				ip.ID,
				func() string {
					if ip.IPAddress != nil {
						return *ip.IPAddress
					} else {
						return "N/A"
					}
				}(),
				ip.ExternalID)
		}
	}

	// Show listeners
	if len(lb.Listeners) > 0 {
		fmt.Println("  Listeners:")
		for _, listener := range lb.Listeners {
			fmt.Printf("    - ID: %s, Nome: %s, Protocolo: %s, Porta: %d\n",
				listener.ID, listener.Name, listener.Protocol, listener.Port)
		}
	}

	// Show backends
	if len(lb.Backends) > 0 {
		fmt.Println("  Backends:")
		for _, backend := range lb.Backends {
			fmt.Printf("    - ID: %s, Nome: %s, Algoritmo: %s, Tipo de Targets: %s\n",
				backend.ID, backend.Name, backend.BalanceAlgorithm, backend.TargetsType)
		}
	}
}

// ExampleUpdateLoadBalancer performs an in-place update on a load balancer,
// such as changing its name, description, or panic threshold.
//
// Lifecycle:
// - Updates are applied asynchronously by the platform. Poll until the LB leaves "updating".
// - Some fields may be immutable post-creation; consult the API docs for constraints.
func ExampleUpdateLoadBalancer(lbID string) {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Update Load Balancer
	updateReq := lbaas.UpdateNetworkLoadBalancerRequest{
		Name:        helpers.StrPtr("load-balancer-atualizado"),
		Description: helpers.StrPtr("Descrição atualizada via SDK"),
	}

	err := lbaasClient.NetworkLoadBalancers().Update(ctx, lbID, updateReq)
	if err != nil {
		log.Fatal("Erro ao atualizar Load Balancer:", err)
	}

	fmt.Printf("Load Balancer %s atualizado com sucesso!\n", lbID)
}

// ExampleManageBackends demonstrates how to:
// - List existing backends for a load balancer
// - Create a new backend with raw targets
// - Retrieve the newly created backend’s details
//
// Tips:
// - Use instances or raw targets according to your environment topology.
// - Keep health checks aligned with backend target behavior.
func ExampleManageBackends(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Listar backends existentes
	listReq := lbaas.ListNetworkBackendRequest{
		LoadBalancerID: lbID,
	}

	backends, err := lbaasClient.NetworkBackends().List(ctx, listReq)
	if err != nil {
		log.Fatal("Erro ao listar backends:", err)
	}

	fmt.Printf("Backends do Load Balancer %s:\n", lbID)
	for _, backend := range backends {
		fmt.Printf("  - ID: %s, Nome: %s, Algoritmo: %s\n",
			backend.ID, backend.Name, backend.BalanceAlgorithm)
	}

	// Criar um novo backend
	createBackendReq := lbaas.CreateNetworkBackendRequest{
		LoadBalancerID:   lbID,
		Name:             "novo-backend-" + time.Now().Format("150405"),
		Description:      helpers.StrPtr("Backend criado via exemplo"),
		BalanceAlgorithm: lbaas.BackendBalanceAlgorithmRoundRobin,
		TargetsType:      lbaas.BackendTypeRaw,
		Targets: &lbaas.TargetsRawOrInstancesRequest{
			TargetsRaw: []lbaas.NetworkBackendRawTargetRequest{
				{
					IPAddress: "192.168.1.20",
					Port:      8080,
				},
			},
		},
	}

	backendID, err := lbaasClient.NetworkBackends().Create(ctx, createBackendReq)
	if err != nil {
		log.Printf("Erro ao criar backend: %v\n", err)
	} else {
		fmt.Printf("Backend criado com sucesso! ID: %s\n", backendID)

		// Obter detalhes do backend criado
		getBackendReq := lbaas.GetNetworkBackendRequest{
			LoadBalancerID: lbID,
			BackendID:      backendID,
		}

		backend, err := lbaasClient.NetworkBackends().Get(ctx, getBackendReq)
		if err != nil {
			log.Printf("Erro ao obter backend: %v\n", err)
		} else {
			fmt.Printf("Detalhes do backend %s:\n", backend.ID)
			fmt.Printf("  Nome: %s\n", backend.Name)
			fmt.Printf("  Algoritmo: %s\n", backend.BalanceAlgorithm)
			fmt.Printf("  Tipo de Targets: %s\n", backend.TargetsType)
		}
	}
}

// ExampleManageListeners demonstrates how to list listeners associated with a load balancer.
// Use this to discover IDs and properties before updating or deleting listeners.
//
// Extending this example:
// - Add create/update/delete flows similar to backends and health checks.
// - Enforce protocol/port policies according to your security standards.
func ExampleManageListeners(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Listar listeners existentes
	listReq := lbaas.ListNetworkListenerRequest{
		LoadBalancerID: lbID,
	}

	listeners, err := lbaasClient.NetworkListeners().List(ctx, listReq)
	if err != nil {
		log.Fatal("Erro ao listar listeners:", err)
	}

	fmt.Printf("Listeners do Load Balancer %s:\n", lbID)
	for _, listener := range listeners {
		fmt.Printf("  - ID: %s, Nome: %s, Protocolo: %s, Porta: %d\n",
			listener.ID, listener.Name, listener.Protocol, listener.Port)
	}
}

// ExampleManageHealthChecks demonstrates how to:
// - List health checks on a load balancer
// - Create a new health check with common settings
//
// Guidance:
// - Tune intervals, timeouts, and thresholds to match your application’s responsiveness.
// - For HTTP checks, consider path and expected status codes.
//
// Returns:
// - Prints created health check ID for follow-up operations.
func ExampleManageHealthChecks(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Listar health checks existentes
	listReq := lbaas.ListNetworkHealthCheckRequest{
		LoadBalancerID: lbID,
	}

	healthChecks, err := lbaasClient.NetworkHealthChecks().List(ctx, listReq)
	if err != nil {
		log.Fatal("Erro ao listar health checks:", err)
	}

	fmt.Printf("Health Checks do Load Balancer %s:\n", lbID)
	for _, hc := range healthChecks {
		fmt.Printf("  - ID: %s, Nome: %s, Protocolo: %s, Porta: %d\n",
			hc.ID, hc.Name, hc.Protocol, hc.Port)
		if hc.Path != nil {
			fmt.Printf("    Path: %s\n", *hc.Path)
		}
	}

	// Criar um novo health check
	createHCReq := lbaas.CreateNetworkHealthCheckRequest{
		LoadBalancerID:          lbID,
		Name:                    "novo-health-check-" + time.Now().Format("150405"),
		Description:             helpers.StrPtr("Health check criado via exemplo"),
		Protocol:                lbaas.HealthCheckProtocolHTTP,
		Port:                    8080,
		IntervalSeconds:         helpers.IntPtr(15),
		TimeoutSeconds:          helpers.IntPtr(3),
		HealthyThresholdCount:   helpers.IntPtr(2),
		UnhealthyThresholdCount: helpers.IntPtr(2),
	}

	hcID, err := lbaasClient.NetworkHealthChecks().Create(ctx, createHCReq)
	if err != nil {
		log.Printf("Erro ao criar health check: %v\n", err)
	} else {
		fmt.Printf("Health Check criado com sucesso! ID: %s\n", hcID.ID)
	}
}

// generateSelfSignedCertificate gera um certificado autoassinado em memória para testes
//
// generateSelfSignedCertificate produces an in-memory self-signed certificate and private key.
// This is intended for demo purposes only and should NOT be used in production.
//
// Returns:
// - PEM-encoded certificate string
// - PEM-encoded private key string
// - error, if any
//
// Security note:
// - In production, obtain certificates from a trusted CA and store private keys securely.
func generateSelfSignedCertificate() (string, string, error) {
	// Gerar chave privada RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("erro ao gerar chave privada: %v", err)
	}

	// Template do certificado
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Organization"},
			Country:       []string{"BR"},
			Province:      []string{"SP"},
			Locality:      []string{"São Paulo"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // Válido por 1 ano
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost", "example.com", "*.example.com"},
	}

	// Criar o certificado
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("erro ao criar certificado: %v", err)
	}

	// Codificar certificado em PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Codificar chave privada em PEM
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("erro ao codificar chave privada: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})

	return string(certPEM), string(privateKeyPEM), nil
}

// ExampleManageCertificates demonstrates how to:
// - List existing TLS certificates associated with a load balancer
// - Generate a temporary self-signed certificate (demo only)
// - Create a certificate on the load balancer and fetch its details
//
// Best practices:
// - Use real certificates from a trusted CA in production.
// - Avoid logging or printing sensitive key material.
// - Consider certificate rotation policies and expiration monitoring.
func ExampleManageCertificates(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Listar certificados existentes
	listReq := lbaas.ListNetworkCertificateRequest{
		LoadBalancerID: lbID,
	}

	certificates, err := lbaasClient.NetworkCertificates().List(ctx, listReq)
	if err != nil {
		log.Fatal("Erro ao listar certificados:", err)
	}

	fmt.Printf("Certificados TLS do Load Balancer %s:\n", lbID)
	for _, cert := range certificates {
		fmt.Printf("  - ID: %s, Nome: %s\n", cert.ID, cert.Name)
		if cert.Description != nil {
			fmt.Printf("    Descrição: %s\n", *cert.Description)
		}
	}

	// Gerar certificado autoassinado em memória
	fmt.Println("Gerando certificado autoassinado para teste...")
	certPEM, keyPEM, err := generateSelfSignedCertificate()
	if err != nil {
		log.Printf("Erro ao gerar certificado: %v\n", err)
		return
	}

	fmt.Println("Certificado gerado com sucesso!")
	fmt.Printf("Certificado (primeiras 100 chars): %s...\n", certPEM[:100])
	fmt.Printf("Chave privada (primeiras 100 chars): %s...\n", keyPEM[:100])

	certPEM = base64.StdEncoding.EncodeToString([]byte(certPEM))
	keyPEM = base64.StdEncoding.EncodeToString([]byte(keyPEM))

	// Criar certificado no Load Balancer
	createCertReq := lbaas.CreateNetworkCertificateRequest{
		LoadBalancerID: lbID,
		Name:           "certificado-teste-" + time.Now().Format("20060102-150405"),
		Description:    helpers.StrPtr("Certificado autoassinado gerado para teste"),
		Certificate:    certPEM,
		PrivateKey:     keyPEM,
	}

	certID, err := lbaasClient.NetworkCertificates().Create(ctx, createCertReq)
	if err != nil {
		log.Printf("Erro ao criar certificado: %v\n", err)
	} else {
		fmt.Printf("Certificado criado com sucesso! ID: %s\n", certID.ID)

		// Obter detalhes do certificado criado
		getCertReq := lbaas.GetNetworkCertificateRequest{
			LoadBalancerID:   lbID,
			TLSCertificateID: certID.ID,
		}

		cert, err := lbaasClient.NetworkCertificates().Get(ctx, getCertReq)
		if err != nil {
			log.Printf("Erro ao obter certificado: %v\n", err)
		} else {
			fmt.Printf("Detalhes do certificado %s:\n", cert.ID)
			fmt.Printf("  Nome: %s\n", cert.Name)
			if cert.Description != nil {
				fmt.Printf("  Descrição: %s\n", *cert.Description)
			}
			fmt.Printf("  Criado em: %s\n", cert.CreatedAt)
			fmt.Printf("  Atualizado em: %s\n", cert.UpdatedAt)
		}
	}
}

// ExampleManageACLs demonstrates how to create an ACL entry for a load balancer.
// The sample shows an allow rule for a CIDR range and protocol.
//
// Operational tips:
// - Start with deny-all and explicitly allow only what’s required (principle of least privilege).
// - Keep network policies versioned and auditable.
func ExampleManageACLs(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// // Listar ACLs existentes
	// listReq := lbaas.ListNetworkACLRequest{
	// 	LoadBalancerID: lbID,
	// }

	// acls, err := lbaasClient.NetworkACLs().List(ctx, listReq)
	// if err != nil {
	// 	log.Fatal("Erro ao listar ACLs:", err)
	// }

	// fmt.Printf("ACLs do Load Balancer %s:\n", lbID)
	// for _, acl := range acls {
	// 	fmt.Printf("  - ID: %s", acl.ID)
	// 	if acl.Name != nil {
	// 		fmt.Printf(", Nome: %s", *acl.Name)
	// 	}
	// 	fmt.Printf(", Protocolo: %s, Ação: %s, IP: %s\n",
	// 		acl.Protocol, acl.Action, acl.RemoteIPPrefix)
	// }

	// Criar uma nova ACL
	createACLReq := lbaas.CreateNetworkACLRequest{
		LoadBalancerID: lbID,
		Name:           helpers.StrPtr("nova-acl-" + time.Now().Format("150405")),
		Ethertype:      lbaas.AclEtherTypeIPv4,
		Protocol:       lbaas.AclProtocolTCP,
		RemoteIPPrefix: "10.0.0.0/8",
		Action:         lbaas.AclActionTypeAllow,
	}

	aclID, err := lbaasClient.NetworkACLs().Create(ctx, createACLReq)
	if err != nil {
		log.Printf("Erro ao criar ACL: %v\n", err)
	} else {
		fmt.Printf("ACL criada com sucesso! ID: %s\n", aclID)
	}
}

// ExampleDeleteLoadBalancer deletes a load balancer by ID.
// Optionally, it can also delete the attached public IP (set DeletePublicIP to true if desired).
//
// Recommendations:
// - Ensure no dependent resources (e.g., listeners, backends) require cleanup policies.
// - In test environments, wrap deletion with retries to handle transient 409 conflicts.
func ExampleDeleteLoadBalancer(lbID string) error {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Delete Load Balancer
	deleteReq := lbaas.DeleteNetworkLoadBalancerRequest{
		DeletePublicIP: helpers.BoolPtr(false), // Keep the public IP
	}

	err := lbaasClient.NetworkLoadBalancers().Delete(ctx, lbID, deleteReq)
	if err != nil {
		return err
	}

	fmt.Printf("Load Balancer %s deletado com sucesso!\n", lbID)
	return nil
}
