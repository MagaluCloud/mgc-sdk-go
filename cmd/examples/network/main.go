package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/network"
)

const (
	defaultZone    = "a"
	defaultTimeout = 30 * time.Second
)

func main() {
	networkClient := createNetworkClient()

	fmt.Println("\n=== VPC Examples ===")
	demoVPCOperations(networkClient)

	fmt.Println("\n=== Subnet Examples ===")
	demoSubnetOperations(networkClient)

	fmt.Println("\n=== Subnet Pool Examples ===")
	demoSubnetPoolOperations(networkClient)

	fmt.Println("\n=== Security Group Examples ===")
	demoSecurityGroupOperations(networkClient)

	fmt.Println("\n=== Security Group Rule Examples ===")
	demoSecurityGroupRuleOperations(networkClient)

	fmt.Println("\n=== Public IP Examples ===")
	demoPublicIPOperations(networkClient)

	fmt.Println("\n=== Port Examples ===")
	demoPortOperations(networkClient)

	fmt.Println("\n=== NAT Gateway Examples ===")
	demoNATGatewayOperations(networkClient)

	fmt.Println("\n=== Private IP Examples ===")
	demoPrivateIPOperations(networkClient)
}

func createNetworkClient() *network.NetworkClient {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	return network.New(c)
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

func demoVPCOperations(networkClient *network.NetworkClient) {
	listVPCs(networkClient)

	vpcID := createVPC(networkClient)
	getVPCDetails(networkClient, vpcID)
	updateVPC(networkClient, vpcID)
	deleteVPC(networkClient, vpcID)
}

func listVPCs(networkClient *network.NetworkClient) {
	ctx, cancel := getContext()
	defer cancel()

	vpcs, err := networkClient.VPCs().List(ctx)
	if err != nil {
		log.Fatalf("Failed to list VPCs: %v", err)
	}

	fmt.Printf("Found %d VPCs\n", len(vpcs))
	for _, vpc := range vpcs {
		fmt.Printf("VPC: %s (ID: %s)\n", *vpc.Name, *vpc.ID)
		fmt.Printf("  Status: %s\n", vpc.Status)
		fmt.Printf("  External Network: %s\n", *vpc.ExternalNetwork)
		fmt.Printf("  Created At: %s\n", vpc.CreatedAt)
		fmt.Printf("  Subnets: %v\n", vpc.Subnets)
		fmt.Printf("  Security Groups: %v\n\n", vpc.SecurityGroups)
	}
}

func createVPC(networkClient *network.NetworkClient) string {
	ctx, cancel := getContext()
	defer cancel()

	createReq := network.CreateVPCRequest{
		Name:        "example-vpc",
		Description: helpers.StrPtr("VPC created via SDK example"),
	}

	id, err := networkClient.VPCs().Create(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create VPC: %v", err)
	}

	fmt.Printf("Created VPC with ID: %s\n", id)
	return id
}

func getVPCDetails(networkClient *network.NetworkClient, vpcID string) {
	ctx, cancel := getContext()
	defer cancel()

	vpc, err := networkClient.VPCs().Get(ctx, vpcID)
	if err != nil {
		log.Fatalf("Failed to get VPC details: %v", err)
	}

	fmt.Printf("VPC Details for %s:\n", vpcID)
	fmt.Printf("  Name: %s\n", *vpc.Name)
	fmt.Printf("  Status: %s\n", vpc.Status)
	fmt.Printf("  External Network: %s\n", *vpc.ExternalNetwork)
	fmt.Printf("  Created At: %s\n\n", vpc.CreatedAt)
}

func updateVPC(networkClient *network.NetworkClient, vpcID string) {
	ctx, cancel := getContext()
	defer cancel()

	newName := "updated-example-vpc"
	if err := networkClient.VPCs().Rename(ctx, vpcID, newName); err != nil {
		log.Fatalf("Failed to update VPC: %v", err)
	}

	fmt.Printf("VPC %s renamed to '%s'\n", vpcID, newName)
}

func deleteVPC(networkClient *network.NetworkClient, vpcID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.VPCs().Delete(ctx, vpcID); err != nil {
		log.Fatalf("Failed to delete VPC: %v", err)
	}

	fmt.Printf("VPC %s deleted successfully\n", vpcID)
}

func demoSubnetOperations(networkClient *network.NetworkClient) {
	vpcID := createVPC(networkClient)
	defer cleanupVPC(networkClient, vpcID)

	subnetID := createSubnet(networkClient, vpcID)

	listSubnets(networkClient, vpcID)

	getSubnetDetails(networkClient, subnetID)

	updateSubnet(networkClient, subnetID)

	deleteSubnet(networkClient, subnetID)
}

func createSubnet(networkClient *network.NetworkClient, vpcID string) string {
	ctx, cancel := getContext()
	defer cancel()

	createReq := network.SubnetCreateRequest{
		Name:        "example-subnet",
		CIDRBlock:   "172.18.106.0/24",
		IPVersion:   4,
		Description: helpers.StrPtr("Subnet created via SDK example"),
	}

	options := network.SubnetCreateOptions{
		Zone: helpers.StrPtr(defaultZone),
	}

	subnetID, err := networkClient.VPCs().CreateSubnet(ctx, vpcID, createReq, options)
	if err != nil {
		log.Fatalf("Failed to create subnet: %v", err)
	}

	fmt.Printf("Created subnet %s in VPC %s\n", subnetID, vpcID)
	return subnetID
}

func listSubnets(networkClient *network.NetworkClient, vpcID string) {
	ctx, cancel := getContext()
	defer cancel()

	subnets, err := networkClient.VPCs().ListSubnets(ctx, vpcID)
	if err != nil {
		log.Fatalf("Failed to list subnets: %v", err)
	}

	fmt.Printf("Subnets in VPC %s:\n", vpcID)
	for _, subnet := range subnets {
		fmt.Printf("  Subnet: %s (CIDR: %s)\n", subnet.ID, subnet.CIDRBlock)
		fmt.Printf("    Name: %s\n", *subnet.Name)
		fmt.Printf("    Zone: %s\n\n", subnet.Zone)
	}
}

func getSubnetDetails(networkClient *network.NetworkClient, subnetID string) {
	ctx, cancel := getContext()
	defer cancel()

	subnet, err := networkClient.Subnets().Get(ctx, subnetID)
	if err != nil {
		log.Fatalf("Failed to get subnet details: %v", err)
	}

	fmt.Printf("Subnet Details for %s:\n", subnetID)
	fmt.Printf("  Name: %s\n", *subnet.Name)
	fmt.Printf("  CIDR Block: %s\n", subnet.CIDRBlock)
	fmt.Printf("  Gateway IP: %s\n", subnet.GatewayIP)
	fmt.Printf("  IP Version: %s\n", subnet.IPVersion)
	fmt.Printf("  Zone: %s\n", subnet.Zone)
	fmt.Printf("  DNS Nameservers: %v\n\n", subnet.DNSNameservers)
}

func updateSubnet(networkClient *network.NetworkClient, subnetID string) {
	ctx, cancel := getContext()
	defer cancel()

	updateReq := network.SubnetPatchRequest{
		DNSNameservers: &[]string{"8.8.8.8", "8.8.4.4"},
	}

	updatedSubnet, err := networkClient.Subnets().Update(ctx, subnetID, updateReq)
	if err != nil {
		log.Fatalf("Failed to update subnet: %v", err)
	}

	fmt.Printf("Updated subnet %s with new DNS servers\n", updatedSubnet.ID)
}

func deleteSubnet(networkClient *network.NetworkClient, subnetID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.Subnets().Delete(ctx, subnetID); err != nil {
		log.Fatalf("Failed to delete subnet: %v", err)
	}

	fmt.Printf("Subnet %s deleted successfully\n", subnetID)
}

func demoSubnetPoolOperations(networkClient *network.NetworkClient) {
	listSubnetPools(networkClient)

	poolID := createSubnetPool(networkClient)

	getSubnetPoolDetails(networkClient, poolID)

	cidrInfo := bookCIDR(networkClient, poolID)

	unbookCIDR(networkClient, poolID, cidrInfo.CIDR)

	deleteSubnetPool(networkClient, poolID)
}

func listSubnetPools(networkClient *network.NetworkClient) {
	ctx, cancel := getContext()
	defer cancel()

	pools, err := networkClient.SubnetPools().List(ctx, network.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
	})
	if err != nil {
		log.Fatalf("Failed to list subnet pools: %v", err)
	}

	fmt.Printf("Found %d subnet pools\n", len(pools))
	for _, pool := range pools {
		fmt.Printf("  ID: %s\n", pool.ID)
		fmt.Printf("  Name: %s\n", pool.Name)
		fmt.Printf("  CIDR: %s\n", *pool.CIDR)
		fmt.Printf("  Is Default: %v\n\n", pool.IsDefault)
	}
}

func createSubnetPool(networkClient *network.NetworkClient) string {
	ctx, cancel := getContext()
	defer cancel()

	createReq := network.CreateSubnetPoolRequest{
		Name:        "example-subnet-pool",
		Description: "Subnet pool created via SDK example",
		CIDR:        helpers.StrPtr("192.168.0.0/16"),
	}

	poolID, err := networkClient.SubnetPools().Create(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create subnet pool: %v", err)
	}

	fmt.Printf("Created subnet pool with ID: %s\n", poolID)
	return poolID
}

func getSubnetPoolDetails(networkClient *network.NetworkClient, poolID string) {
	ctx, cancel := getContext()
	defer cancel()

	pool, err := networkClient.SubnetPools().Get(ctx, poolID)
	if err != nil {
		log.Fatalf("Failed to get subnet pool details: %v", err)
	}

	fmt.Printf("Subnet Pool Details for %s:\n", poolID)
	fmt.Printf("  Name: %s\n", pool.Name)
	fmt.Printf("  CIDR: %s\n", *pool.CIDR)
	fmt.Printf("  IP Version: %d\n", pool.IPVersion)
	fmt.Printf("  Created At: %s\n\n", pool.CreatedAt)
}

func bookCIDR(networkClient *network.NetworkClient, poolID string) network.BookCIDRResponse {
	ctx, cancel := getContext()
	defer cancel()

	bookReq := network.BookCIDRRequest{
		Mask: helpers.IntPtr(24),
	}

	bookedCIDR, err := networkClient.SubnetPools().BookCIDR(ctx, poolID, bookReq)
	if err != nil {
		log.Fatalf("Failed to book CIDR: %v", err)
	}

	fmt.Printf("Booked CIDR: %s from pool %s\n", bookedCIDR.CIDR, poolID)
	return *bookedCIDR
}

func unbookCIDR(networkClient *network.NetworkClient, poolID, cidr string) {
	ctx, cancel := getContext()
	defer cancel()

	unbookReq := network.UnbookCIDRRequest{
		CIDR: cidr,
	}

	if err := networkClient.SubnetPools().UnbookCIDR(ctx, poolID, unbookReq); err != nil {
		log.Fatalf("Failed to unbook CIDR: %v", err)
	}

	fmt.Printf("Unbooked CIDR: %s from pool %s\n", cidr, poolID)
}

func deleteSubnetPool(networkClient *network.NetworkClient, poolID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.SubnetPools().Delete(ctx, poolID); err != nil {
		log.Fatalf("Failed to delete subnet pool: %v", err)
	}

	fmt.Printf("Subnet pool %s deleted successfully\n", poolID)
}

func demoSecurityGroupOperations(networkClient *network.NetworkClient) {
	listSecurityGroups(networkClient)

	sgID := createSecurityGroup(networkClient)

	getSecurityGroupDetails(networkClient, sgID)

	deleteSecurityGroup(networkClient, sgID)
}

func listSecurityGroups(networkClient *network.NetworkClient) {
	ctx, cancel := getContext()
	defer cancel()

	securityGroups, err := networkClient.SecurityGroups().List(ctx)
	if err != nil {
		log.Fatalf("Failed to list security groups: %v", err)
	}

	fmt.Printf("Found %d security groups\n", len(securityGroups))
	for _, sg := range securityGroups {
		fmt.Printf("  ID: %s\n", *sg.ID)
		fmt.Printf("  Name: %s\n", *sg.Name)
		fmt.Printf("  VPC ID: %s\n", *sg.VPCID)
		fmt.Printf("  Is Default: %v\n", sg.IsDefault)
		fmt.Printf("  Ports: %v\n\n", sg.Ports)
	}
}

func createSecurityGroup(networkClient *network.NetworkClient) string {
	ctx, cancel := getContext()
	defer cancel()

	createReq := network.SecurityGroupCreateRequest{
		Name:             "example-security-group",
		Description:      helpers.StrPtr("Security group created via SDK example"),
		SkipDefaultRules: helpers.BoolPtr(true),
	}

	sgID, err := networkClient.SecurityGroups().Create(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create security group: %v", err)
	}

	fmt.Printf("Created security group with ID: %s\n", sgID)
	return sgID
}

func getSecurityGroupDetails(networkClient *network.NetworkClient, sgID string) {
	ctx, cancel := getContext()
	defer cancel()

	sg, err := networkClient.SecurityGroups().Get(ctx, sgID)
	if err != nil {
		log.Fatalf("Failed to get security group details: %v", err)
	}

	fmt.Printf("Security Group Details for %s:\n", sgID)
	fmt.Printf("  Name: %s\n", *sg.Name)
	fmt.Printf("  Description: %s\n", *sg.Description)
	fmt.Printf("  VPC ID: %s\n", *sg.VPCID)
	fmt.Printf("  Status: %s\n", sg.Status)

	if sg.Rules != nil && len(*sg.Rules) > 0 {
		fmt.Println("  Rules:")
		for _, rule := range *sg.Rules {
			fmt.Printf("    Direction: %s, Protocol: %s\n", *rule.Direction, *rule.Protocol)
		}
	} else {
		fmt.Println("  No rules defined")
	}
	fmt.Println()
}

func deleteSecurityGroup(networkClient *network.NetworkClient, sgID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.SecurityGroups().Delete(ctx, sgID); err != nil {
		log.Fatalf("Failed to delete security group: %v", err)
	}

	fmt.Printf("Security group %s deleted successfully\n", sgID)
}

func demoSecurityGroupRuleOperations(networkClient *network.NetworkClient) {
	sgID := createSecurityGroup(networkClient)
	defer deleteSecurityGroup(networkClient, sgID)

	sshRuleID := createSSHSecurityRule(networkClient, sgID)
	httpsRuleID := createHTTPSSecurityRule(networkClient, sgID)

	listSecurityGroupRules(networkClient, sgID)

	getSecurityGroupRuleDetails(networkClient, sshRuleID)

	deleteSecurityGroupRule(networkClient, sshRuleID)
	deleteSecurityGroupRule(networkClient, httpsRuleID)
}

func createSSHSecurityRule(networkClient *network.NetworkClient, sgID string) string {
	ctx, cancel := getContext()
	defer cancel()

	sshRule := network.RuleCreateRequest{
		Direction:      helpers.StrPtr("ingress"),
		PortRangeMin:   helpers.IntPtr(22),
		PortRangeMax:   helpers.IntPtr(22),
		Protocol:       helpers.StrPtr("tcp"),
		RemoteIPPrefix: helpers.StrPtr("0.0.0.0/0"),
		EtherType:      "IPv4",
		Description:    helpers.StrPtr("Allow SSH access"),
	}

	ruleID, err := networkClient.Rules().Create(ctx, sgID, sshRule)
	if err != nil {
		log.Fatalf("Failed to create SSH security rule: %v", err)
	}

	fmt.Printf("Created SSH rule with ID: %s in security group %s\n", ruleID, sgID)
	return ruleID
}

func createHTTPSSecurityRule(networkClient *network.NetworkClient, sgID string) string {
	ctx, cancel := getContext()
	defer cancel()

	httpsRule := network.RuleCreateRequest{
		Direction:      helpers.StrPtr("ingress"),
		PortRangeMin:   helpers.IntPtr(443),
		PortRangeMax:   helpers.IntPtr(443),
		Protocol:       helpers.StrPtr("tcp"),
		RemoteIPPrefix: helpers.StrPtr("0.0.0.0/0"),
		EtherType:      "IPv4",
		Description:    helpers.StrPtr("Allow HTTPS access"),
	}

	ruleID, err := networkClient.Rules().Create(ctx, sgID, httpsRule)
	if err != nil {
		log.Fatalf("Failed to create HTTPS security rule: %v", err)
	}

	fmt.Printf("Created HTTPS rule with ID: %s in security group %s\n", ruleID, sgID)
	return ruleID
}

func listSecurityGroupRules(networkClient *network.NetworkClient, sgID string) {
	ctx, cancel := getContext()
	defer cancel()

	rules, err := networkClient.Rules().List(ctx, sgID)
	if err != nil {
		log.Fatalf("Failed to list security group rules: %v", err)
	}

	fmt.Printf("Rules in security group %s:\n", sgID)
	for _, rule := range rules {
		fmt.Printf("  Rule ID: %s\n", *rule.ID)
		fmt.Printf("    Direction: %s\n", *rule.Direction)
		fmt.Printf("    Protocol: %s\n", *rule.Protocol)
		if rule.PortRangeMin != nil {
			fmt.Printf("    Port Range: %d-%d\n", *rule.PortRangeMin, *rule.PortRangeMax)
		}
		fmt.Printf("    Remote IP Prefix: %s\n\n", *rule.RemoteIPPrefix)
	}
}

func getSecurityGroupRuleDetails(networkClient *network.NetworkClient, ruleID string) {
	ctx, cancel := getContext()
	defer cancel()

	rule, err := networkClient.Rules().Get(ctx, ruleID)
	if err != nil {
		log.Fatalf("Failed to get security group rule details: %v", err)
	}

	fmt.Printf("Security Rule Details for %s:\n", ruleID)
	fmt.Printf("  Direction: %s\n", *rule.Direction)
	fmt.Printf("  Protocol: %s\n", *rule.Protocol)
	if rule.PortRangeMin != nil {
		fmt.Printf("  Port Range: %d-%d\n", *rule.PortRangeMin, *rule.PortRangeMax)
	}
	fmt.Printf("  Remote IP Prefix: %s\n", *rule.RemoteIPPrefix)
	fmt.Printf("  Description: %s\n\n", *rule.Description)
}

func deleteSecurityGroupRule(networkClient *network.NetworkClient, ruleID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.Rules().Delete(ctx, ruleID); err != nil {
		log.Fatalf("Failed to delete security group rule: %v", err)
	}

	fmt.Printf("Security rule %s deleted successfully\n", ruleID)
}

func demoPublicIPOperations(networkClient *network.NetworkClient) {
	vpcID := createVPC(networkClient)
	defer cleanupVPC(networkClient, vpcID)

	portID := createPort(networkClient, vpcID, "test-port-for-pip", true, nil, nil)

	listPublicIPs(networkClient)

	pipID := createPublicIP(networkClient, vpcID)

	getPublicIPDetails(networkClient, pipID)

	attachPublicIPToPort(networkClient, pipID, portID)

	detachPublicIPFromPort(networkClient, pipID, portID)

	deletePublicIP(networkClient, pipID)
}

func listPublicIPs(networkClient *network.NetworkClient) {
	ctx, cancel := getContext()
	defer cancel()

	publicIPs, err := networkClient.PublicIPs().List(ctx)
	if err != nil {
		log.Fatalf("Failed to list public IPs: %v", err)
	}

	fmt.Printf("Found %d public IPs\n", len(publicIPs))
	for i, pip := range publicIPs {
		if i >= 5 {
			fmt.Printf("And %d more...\n", len(publicIPs)-5)
			break
		}
		fmt.Printf("  ID: %s\n", *pip.ID)
		fmt.Printf("  Public IP: %s\n", *pip.PublicIP)
		fmt.Printf("  VPC ID: %s\n\n", *pip.VPCID)
	}
}

func createPublicIP(networkClient *network.NetworkClient, vpcID string) string {
	ctx, cancel := getContext()
	defer cancel()

	pipReq := network.PublicIPCreateRequest{
		Description: helpers.StrPtr("Public IP created via SDK example"),
	}

	pipID, err := networkClient.VPCs().CreatePublicIP(ctx, vpcID, pipReq)
	if err != nil {
		log.Fatalf("Failed to create public IP: %v", err)
	}

	fmt.Printf("Created public IP with ID: %s in VPC %s\n", pipID, vpcID)
	return pipID
}

func getPublicIPDetails(networkClient *network.NetworkClient, pipID string) {
	ctx, cancel := getContext()
	defer cancel()

	pip, err := networkClient.PublicIPs().Get(ctx, pipID)
	if err != nil {
		log.Fatalf("Failed to get public IP details: %v", err)
	}

	fmt.Printf("Public IP Details for %s:\n", pipID)
	fmt.Printf("  Public IP: %s\n", *pip.PublicIP)
	fmt.Printf("  Description: %s\n", *pip.Description)
	fmt.Printf("  Status: %s\n", *pip.Status)
	fmt.Printf("  VPC ID: %s\n", *pip.VPCID)
	if pip.PortID != nil {
		fmt.Printf("  Port ID: %s\n", *pip.PortID)
	} else {
		fmt.Println("  Not attached to any port")
	}
	fmt.Println()
}

func attachPublicIPToPort(networkClient *network.NetworkClient, pipID, portID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.PublicIPs().AttachToPort(ctx, pipID, portID); err != nil {
		log.Fatalf("Failed to attach public IP to port: %v", err)
	}

	fmt.Printf("Attached public IP %s to port %s\n", pipID, portID)
}

func detachPublicIPFromPort(networkClient *network.NetworkClient, pipID, portID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.PublicIPs().DetachFromPort(ctx, pipID, portID); err != nil {
		log.Fatalf("Failed to detach public IP from port: %v", err)
	}

	fmt.Printf("Detached public IP %s from port %s\n", pipID, portID)
}

func deletePublicIP(networkClient *network.NetworkClient, pipID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.PublicIPs().Delete(ctx, pipID); err != nil {
		log.Fatalf("Failed to delete public IP: %v", err)
	}

	fmt.Printf("Public IP %s deleted successfully\n", pipID)
}

func demoPortOperations(networkClient *network.NetworkClient) {
	vpcID := createVPC(networkClient)
	defer cleanupVPC(networkClient, vpcID)

	sgID := createSecurityGroup(networkClient)
	defer deleteSecurityGroup(networkClient, sgID)

	portID := createPort(networkClient, vpcID, "example-port", false, nil, nil)

	getPortDetails(networkClient, portID)

	attachSecurityGroupToPort(networkClient, portID, sgID)
	detachSecurityGroupFromPort(networkClient, portID, sgID)

	listPorts(networkClient, vpcID)

	deletePort(networkClient, portID)

	updatePort(networkClient, portID)
}

func createPort(networkClient *network.NetworkClient, vpcID, name string, hasPIP bool, ipAddress *string, subnets *[]string) string {
	ctx, cancel := getContext()
	defer cancel()

	portReq := network.PortCreateRequest{
		Name:   name,
		HasPIP: helpers.BoolPtr(hasPIP),
		HasSG:  helpers.BoolPtr(true),
	}

	if ipAddress != nil {
		portReq.IPAddress = ipAddress
	}

	if subnets != nil {
		portReq.Subnets = subnets
	}

	options := network.PortCreateOptions{
		Zone: helpers.StrPtr(defaultZone),
	}

	portID, err := networkClient.VPCs().CreatePort(ctx, vpcID, portReq, options)
	if err != nil {
		log.Fatalf("Failed to create port: %v", err)
	}

	fmt.Printf("Created port with ID: %s in VPC %s\n", portID, vpcID)
	return portID
}

func getPortDetails(networkClient *network.NetworkClient, portID string) {
	ctx, cancel := getContext()
	defer cancel()

	port, err := networkClient.Ports().Get(ctx, portID)
	if err != nil {
		log.Fatalf("Failed to get port details: %v", err)
	}

	fmt.Printf("Port Details for %s:\n", portID)
	fmt.Printf("  Name: %s\n", *port.Name)
	fmt.Printf("  VPC ID: %s\n", *port.VPCID)
	fmt.Printf("  Security Groups: %v\n", *port.SecurityGroups)
	fmt.Printf("  IP Spoofing Guard: %t\n", *port.IPSpoofingGuard)
	if port.IPAddress != nil {
		fmt.Println("  IP Addresses:")
		for _, ip := range *port.IPAddress {
			fmt.Printf("    %s (Subnet: %s)\n", ip.IPAddress, ip.SubnetID)
		}
	}
	if port.PublicIP != nil {
		fmt.Println("  Public IPs:")
		for _, pip := range *port.PublicIP {
			fmt.Printf("    %s (ID: %s)\n", *pip.PublicIP, *pip.PublicIPID)
		}
	}
	fmt.Println()
}

func attachSecurityGroupToPort(networkClient *network.NetworkClient, portID, sgID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.Ports().AttachSecurityGroup(ctx, portID, sgID); err != nil {
		log.Fatalf("Failed to attach security group to port: %v", err)
	}

	fmt.Printf("Attached security group %s to port %s\n", sgID, portID)
}

func detachSecurityGroupFromPort(networkClient *network.NetworkClient, portID, sgID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.Ports().DetachSecurityGroup(ctx, portID, sgID); err != nil {
		log.Fatalf("Failed to detach security group from port: %v", err)
	}

	fmt.Printf("Detached security group %s from port %s\n", sgID, portID)
}

func listPorts(networkClient *network.NetworkClient, vpcID string) {
	ctx, cancel := getContext()
	defer cancel()

	ports, err := networkClient.VPCs().ListPorts(ctx, vpcID, true, network.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list ports: %v", err)
	}

	fmt.Printf("Ports in VPC %s:\n", vpcID)
	if ports.Ports != nil {
		for _, port := range *ports.Ports {
			fmt.Printf("  Port: %s\n", *port.ID)
		}
	}
}

func deletePort(networkClient *network.NetworkClient, portID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.Ports().Delete(ctx, portID); err != nil {
		log.Fatalf("Failed to delete port: %v", err)
	}

	fmt.Printf("Port %s deleted successfully\n", portID)
}

func cleanupVPC(networkClient *network.NetworkClient, vpcID string) {

	for {
		vpc, err := networkClient.VPCs().Get(context.Background(), vpcID)
		if err != nil {
			log.Fatalf("Failed to get VPC details: %v", err)
		}
		fmt.Printf("VPC status: %s\n", vpc.Status)
		if vpc.Status == "created" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	deleteVPC(networkClient, vpcID)
}

func demoNATGatewayOperations(networkClient *network.NetworkClient) {
	vpcID := createVPC(networkClient)
	defer cleanupVPC(networkClient, vpcID)

	natGatewayID := createNATGateway(networkClient, vpcID)

	getNATGatewayDetails(networkClient, natGatewayID)
	listNATGateways(networkClient, vpcID)
	deleteNATGateway(networkClient, natGatewayID)
}

func createNATGateway(networkClient *network.NetworkClient, vpcID string) string {
	ctx, cancel := getContext()
	defer cancel()

	natGatewayID, err := networkClient.NatGateways().Create(ctx, network.CreateNatGatewayRequest{
		Name:        "example-nat-gateway",
		Description: helpers.StrPtr("NAT gateway created via SDK example"),
		Zone:        defaultZone,
		VPCID:       vpcID,
	})
	if err != nil {
		log.Fatalf("Failed to create NAT gateway: %v", err)
	}

	return natGatewayID
}

func listNATGateways(networkClient *network.NetworkClient, vpcID string) {
	ctx, cancel := getContext()
	defer cancel()

	natGateways, err := networkClient.NatGateways().List(ctx, vpcID, network.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list NAT gateways: %v", err)
	}

	fmt.Printf("Found %d NAT gateways\n", len(natGateways))
	for _, natGateway := range natGateways {
		fmt.Printf("  NAT Gateway: %s\n", *natGateway.ID)
		fmt.Printf("    Name: %s\n", *natGateway.Name)
		fmt.Printf("    Zone: %s\n", *natGateway.Zone)
		fmt.Printf("    VPC ID: %s\n", *natGateway.VPCID)
		fmt.Printf("    Description: %s\n", *natGateway.Description)
	}
}

func getNATGatewayDetails(networkClient *network.NetworkClient, natGatewayID string) {
	ctx, cancel := getContext()
	defer cancel()

	natGateway, err := networkClient.NatGateways().Get(ctx, natGatewayID)
	if err != nil {
		log.Fatalf("Failed to get NAT gateway details: %v", err)
	}

	fmt.Printf("NAT Gateway Details for %s:\n", natGatewayID)
	fmt.Printf("  Name: %s\n", *natGateway.Name)

}

func deleteNATGateway(networkClient *network.NetworkClient, natGatewayID string) {
	ctx, cancel := getContext()
	defer cancel()

	if err := networkClient.NatGateways().Delete(ctx, natGatewayID); err != nil {
		log.Fatalf("Failed to delete NAT gateway: %v", err)
	}

	fmt.Printf("NAT Gateway %s deleted successfully\n", natGatewayID)
}

func updatePort(networkClient *network.NetworkClient, portID string) {
	ctx, cancel := getContext()
	defer cancel()

	portUpdateRequest := &network.PortUpdateRequest{
		IPSpoofingGuard: helpers.BoolPtr(false),
	}

	if err := networkClient.Ports().Update(ctx, portID, *portUpdateRequest); err != nil {
		log.Fatalf("Failed to update port: %v", err)
	}

	fmt.Printf("Port %s updated successfully\n", portID)
}

func demoPrivateIPOperations(networkClient *network.NetworkClient) {
	vpcID := createVPC(networkClient)
	defer cleanupVPC(networkClient, vpcID)

	subnetID := createSubnet(networkClient, vpcID)

	subnets := []string{subnetID}

	ipAddress := "172.18.106.10"

	createPort(networkClient, vpcID, "port-example", false, &ipAddress, &subnets)
}
