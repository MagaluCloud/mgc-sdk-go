# Lbaas

Example usage of the `lbaas` module.

Exemplos de operações com Load Balancer existente

**File:** `cmd/examples/lbaas/main.go`

```go
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

func main() {
	ExampleListLoadBalancers()

	vpcID := "b44e3fe0-b609-4906-81a3-256e0cf68cfb"
	lbID := ExampleCreateLoadBalancer(vpcID)
	time.Sleep(10 * time.Second)

	// Exemplos de operações com Load Balancer existente
	// Substitua por um ID real para testar
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

func waitCreateLoadBalancer(lbID string, status lbaas.LoadBalancerStatus) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Obter detalhes do Load Balancer
	getReq := lbaas.GetNetworkLoadBalancerRequest{
		LoadBalancerID: lbID,
	}

	for {
		lb, err := lbaasClient.NetworkLoadBalancers().Get(ctx, getReq)
		if err != nil {
			log.Fatal("Erro ao obter Load Balancer:", err)
		}
		if lb.Status != string(status) {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func ExampleCreateLoadBalancer(vpcID string) string {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Configurar health check
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

	// Configurar backend com targets raw (IPs)
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

	// Configurar listener
	listener := lbaas.NetworkListenerRequest{
		Name:        "example-listener",
		Description: helpers.StrPtr("Listener HTTP para exemplo"),
		BackendName: "example-backend",
		Protocol:    lbaas.ListenerProtocolTCP,
		Port:        80,
	}

	// Configurar ACL (opcional)
	acl := lbaas.NetworkAclRequest{
		Name:           helpers.StrPtr("allow-all"),
		Ethertype:      lbaas.AclEtherTypeIPv4,
		Protocol:       lbaas.AclProtocolTCP,
		RemoteIPPrefix: "0.0.0.0/0",
		Action:         lbaas.AclActionTypeAllow,
	}

	// Criar Load Balancer
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
		// PublicIPID:   helpers.StrPtr("your-public-ip-id"), // Opcional
		PanicThreshold: helpers.IntPtr(50),
	}

	id, err := lbaasClient.NetworkLoadBalancers().Create(ctx, createReq)
	if err != nil {
		log.Fatal("Erro ao criar Load Balancer:", err)
	}

	fmt.Printf("Load Balancer criado com sucesso! ID: %s\n", id)
	return id
}

func ExampleListLoadBalancers() {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Listar Load Balancers com paginação
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

func ExampleGetLoadBalancer(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Obter detalhes do Load Balancer
	getReq := lbaas.GetNetworkLoadBalancerRequest{
		LoadBalancerID: lbID,
	}

	lb, err := lbaasClient.NetworkLoadBalancers().Get(ctx, getReq)
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

	// Exibir IPs públicos
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

	// Exibir listeners
	if len(lb.Listeners) > 0 {
		fmt.Println("  Listeners:")
		for _, listener := range lb.Listeners {
			fmt.Printf("    - ID: %s, Nome: %s, Protocolo: %s, Porta: %d\n",
				listener.ID, listener.Name, listener.Protocol, listener.Port)
		}
	}

	// Exibir backends
	if len(lb.Backends) > 0 {
		fmt.Println("  Backends:")
		for _, backend := range lb.Backends {
			fmt.Printf("    - ID: %s, Nome: %s, Algoritmo: %s, Tipo de Targets: %s\n",
				backend.ID, backend.Name, backend.BalanceAlgorithm, backend.TargetsType)
		}
	}
}

func ExampleUpdateLoadBalancer(lbID string) {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Atualizar Load Balancer
	updateReq := lbaas.UpdateNetworkLoadBalancerRequest{
		LoadBalancerID: lbID,
		Name:           helpers.StrPtr("load-balancer-atualizado"),
		Description:    helpers.StrPtr("Descrição atualizada via SDK"),
		PanicThreshold: helpers.IntPtr(75),
	}

	err := lbaasClient.NetworkLoadBalancers().Update(ctx, updateReq)
	if err != nil {
		log.Fatal("Erro ao atualizar Load Balancer:", err)
	}

	fmt.Printf("Load Balancer %s atualizado com sucesso!\n", lbID)
}

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

func ExampleDeleteLoadBalancer(lbID string) error {
	// Criar um novo cliente
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	lbaasClient := lbaas.New(c)
	ctx := context.Background()

	// Deletar Load Balancer
	deleteReq := lbaas.DeleteNetworkLoadBalancerRequest{
		LoadBalancerID: lbID,
		DeletePublicIP: helpers.BoolPtr(false), // Manter o IP público
	}

	err := lbaasClient.NetworkLoadBalancers().Delete(ctx, deleteReq)
	if err != nil {
		return err
	}

	fmt.Printf("Load Balancer %s deletado com sucesso!\n", lbID)
	return nil
}
```
