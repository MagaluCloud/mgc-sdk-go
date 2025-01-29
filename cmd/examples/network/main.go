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
	waitTimeout   = 5 * time.Minute
	retryInterval = 5 * time.Second
)

func main() {
	ExampleListVPCs()
	ExamplePublicIPs()
	ExampleSecurityGroups()
	ExampleSecurityGroupRules()
	ExampleSubnetPools()
	ExampleSubnets()
	ExamplePorts()
	id := ExampleCreateVPC()
	ExampleManageVPC(id)
	ExampleManageSubnets(id)
	ExampleManagePorts(id)
	ExampleDeleteVPC(id)
}

func ExampleListVPCs() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)

	// List VPCs with pagination and expansion
	vpcs, err := networkClient.VPCs().List(context.Background(), network.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{network.SubnetsExpand, network.SecurityGroupsExpand},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Print VPC details
	for _, vpc := range vpcs {
		fmt.Printf("VPC: %s (ID: %s)\n", vpc.Name, vpc.ID)
		fmt.Printf("  Status: %s\n", vpc.Status)
		fmt.Printf("  Router ID: %s\n", vpc.RouterID)
		fmt.Printf("  External Network: %s\n", vpc.ExternalNetwork)
		fmt.Printf("  Created At: %s\n", vpc.CreatedAt)
		fmt.Printf("  Subnets: %v\n", vpc.Subnets)
		fmt.Printf("  Security Groups: %v\n", vpc.SecurityGroups)
	}
}

func ExampleCreateVPC() string {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)

	// Create a new VPC
	createReq := network.CreateVPCRequest{
		Name:        "my-test-vpc",
		Description: "Test VPC created via SDK",
	}

	id, err := networkClient.VPCs().Create(context.Background(), createReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created VPC with ID: %s\n", id)
	return id
}

func ExampleManageVPC(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// Get VPC details
	vpc, err := networkClient.VPCs().Get(ctx, id, []string{network.SubnetsExpand})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("VPC: %s (ID: %s)\n", vpc.Name, vpc.ID)

	// Rename the VPC
	if err := networkClient.VPCs().Rename(ctx, vpc.ID, "new-vpc-name"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("VPC renamed successfully")
}

func ExampleManageSubnets(vpcID string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// List subnets
	subnets, err := networkClient.VPCs().ListSubnets(ctx, vpcID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Subnets in VPC %s:\n", vpcID)
	for _, subnet := range subnets {
		fmt.Printf("  Subnet: %s (CIDR: %s)\n", subnet.ID, subnet.CIDRBlock)
	}

	// Create a new subnet
	createSubnetReq := network.SubnetCreateRequest{
		Name:        "my-subnet",
		CIDRBlock:   "192.168.1.0/24",
		Description: "Test subnet created via SDK",
	}

	subnetID, err := networkClient.VPCs().CreateSubnet(ctx, vpcID, createSubnetReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created subnet with ID: %s\n", subnetID)
}

func ExampleManagePorts(vpcID string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// List ports
	ports, err := networkClient.VPCs().ListPorts(ctx, vpcID, true, network.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ports in VPC %s:\n", vpcID)
	if portList, ok := ports.([]network.PortResponse); ok {
		for _, port := range portList {
			fmt.Printf("  Port: %s\n", port.ID)
		}
	}

	// Create a new port
	createPortReq := network.PortCreateRequest{
		Name:   "my-port",
		HasPIP: true,
		HasSG:  true,
	}

	portID, err := networkClient.VPCs().CreatePort(ctx, vpcID, createPortReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created port with ID: %s\n", portID)
}

func ExampleDeleteVPC(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)

	if err := networkClient.VPCs().Delete(context.Background(), id); err != nil {
		log.Fatal(err)
	}

	fmt.Println("VPC deleted successfully")
}

func ExampleSubnets() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// Example subnet ID - replace with an actual subnet ID
	subnetID := "030d0e77-e9f9-4af8-99db-067eba6826c0"

	// Get subnet details
	subnet, err := networkClient.Subnets().Get(ctx, subnetID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Subnet Details:\n")
	fmt.Printf("  ID: %s\n", subnet.ID)
	fmt.Printf("  Name: %s\n", subnet.Name)
	fmt.Printf("  CIDR Block: %s\n", subnet.CIDRBlock)
	fmt.Printf("  Gateway IP: %s\n", subnet.GatewayIP)
	fmt.Printf("  IP Version: %s\n", subnet.IPVersion)
	fmt.Printf("  Zone: %s\n", subnet.Zone)
	fmt.Printf("  DNS Nameservers: %v\n", subnet.DNSNameservers)

	// Update subnet DNS nameservers
	updateReq := network.SubnetPatchRequest{
		DNSNameservers: []string{"8.8.8.8", "8.8.4.4"},
	}

	updatedSubnet, err := networkClient.Subnets().Update(ctx, subnetID, updateReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated subnet ID: %s\n", updatedSubnet.ID)

	// Delete subnet
	if err := networkClient.Subnets().Delete(ctx, subnetID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Subnet deleted successfully")
}

func ExampleSubnetPools() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// List subnet pools
	pools, err := networkClient.SubnetPools().List(ctx, network.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available Subnet Pools:")
	for _, pool := range pools {
		fmt.Printf("  ID: %s\n", pool.ID)
		fmt.Printf("  Name: %s\n", pool.Name)
		fmt.Printf("  CIDR: %s\n", pool.CIDR)
		fmt.Printf("  Is Default: %v\n", pool.IsDefault)
		fmt.Printf("  Description: %s\n\n", pool.Description)
	}

	// Create a new subnet pool
	createReq := network.CreateSubnetPoolRequest{
		Name:        "my-subnet-pool",
		Description: "Test subnet pool created via SDK",
		CIDR:        "192.168.0.0/16",
	}

	poolID, err := networkClient.SubnetPools().Create(ctx, createReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created subnet pool with ID: %s\n", poolID)

	// Get subnet pool details
	pool, err := networkClient.SubnetPools().Get(ctx, poolID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nSubnet Pool Details:\n")
	fmt.Printf("  ID: %s\n", pool.ID)
	fmt.Printf("  Name: %s\n", pool.Name)
	fmt.Printf("  CIDR: %s\n", pool.CIDR)
	fmt.Printf("  IP Version: %d\n", pool.IPVersion)
	fmt.Printf("  Created At: %s\n", pool.CreatedAt)

	// Book a CIDR from the pool
	bookReq := network.BookCIDRRequest{
		Mask: helpers.IntPtr(24), // Request a /24 subnet
	}

	bookedCIDR, err := networkClient.SubnetPools().BookCIDR(ctx, poolID, bookReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nBooked CIDR: %s\n", bookedCIDR.CIDR)

	// Unbook the CIDR
	unbookReq := network.UnbookCIDRRequest{
		CIDR: bookedCIDR.CIDR,
	}

	if err := networkClient.SubnetPools().UnbookCIDR(ctx, poolID, unbookReq); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Unbooked CIDR: %s\n", bookedCIDR.CIDR)

	// Delete the subnet pool
	if err := networkClient.SubnetPools().Delete(ctx, poolID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Subnet pool deleted successfully")
}

func ExampleSecurityGroups() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// List security groups
	securityGroups, err := networkClient.SecurityGroups().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available Security Groups:")
	for _, sg := range securityGroups {
		fmt.Printf("  ID: %s\n", sg.ID)
		fmt.Printf("  Name: %s\n", sg.Name)
		fmt.Printf("  Description: %s\n", sg.Description)
		fmt.Printf("  VPC ID: %s\n", sg.VPCID)
		fmt.Printf("  Status: %s\n", sg.Status)
		fmt.Printf("  Is Default: %v\n", sg.IsDefault)
		fmt.Printf("  Created At: %s\n\n", sg.CreatedAt)
	}

	// Create a new security group
	createReq := network.SecurityGroupCreateRequest{
		Name:        "my-security-group",
		Description: "Test security group created via SDK",
	}

	sgID, err := networkClient.SecurityGroups().Create(ctx, createReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created security group with ID: %s\n", sgID)

	// Get security group details
	sg, err := networkClient.SecurityGroups().Get(ctx, sgID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nSecurity Group Details:\n")
	fmt.Printf("  ID: %s\n", sg.ID)
	fmt.Printf("  Name: %s\n", sg.Name)
	fmt.Printf("  Description: %s\n", sg.Description)
	fmt.Printf("  VPC ID: %s\n", sg.VPCID)
	fmt.Printf("  Status: %s\n", sg.Status)
	fmt.Printf("  Rules Count: %d\n", len(sg.Rules))

	// Print security rules if any exist
	if len(sg.Rules) > 0 {
		fmt.Println("\nSecurity Rules:")
		for _, rule := range sg.Rules {
			fmt.Printf("  Direction: %s\n", rule.Direction)
			fmt.Printf("  Protocol: %s\n", rule.Protocol)
			if rule.PortRangeMin != nil {
				fmt.Printf("  Port Range Min: %d\n", *rule.PortRangeMin)
			}
			if rule.PortRangeMax != nil {
				fmt.Printf("  Port Range Max: %d\n", *rule.PortRangeMax)
			}
			fmt.Printf("  Remote IP Prefix: %s\n", rule.RemoteIPPrefix)
			fmt.Printf("  Description: %s\n\n", rule.Description)
		}
	}

	// Delete the security group
	if err := networkClient.SecurityGroups().Delete(ctx, sgID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Security group deleted successfully")
}

func ExampleSecurityGroupRules() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// First, create a security group for the rules
	sgCreateReq := network.SecurityGroupCreateRequest{
		Name:        "test-sg-with-rules",
		Description: "Security group for testing rules",
	}

	sgID, err := networkClient.SecurityGroups().Create(ctx, sgCreateReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created security group with ID: %s\n", sgID)

	// Create an inbound rule allowing SSH access
	sshRule := network.RuleCreateRequest{
		Direction:      "ingress",
		PortRangeMin:   helpers.IntPtr(22),
		PortRangeMax:   helpers.IntPtr(22),
		Protocol:       "tcp",
		RemoteIPPrefix: "0.0.0.0/0",
		EtherType:      "IPv4",
		Description:    "Allow SSH access",
	}

	sshRuleID, err := networkClient.Rules().Create(ctx, sgID, sshRule)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created SSH rule with ID: %s\n", sshRuleID)

	// Create an inbound rule allowing HTTPS access
	httpsRule := network.RuleCreateRequest{
		Direction:      "ingress",
		PortRangeMin:   helpers.IntPtr(443),
		PortRangeMax:   helpers.IntPtr(443),
		Protocol:       "tcp",
		RemoteIPPrefix: "0.0.0.0/0",
		EtherType:      "IPv4",
		Description:    "Allow HTTPS access",
	}

	httpsRuleID, err := networkClient.Rules().Create(ctx, sgID, httpsRule)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created HTTPS rule with ID: %s\n", httpsRuleID)

	// List all rules in the security group
	rules, err := networkClient.Rules().List(ctx, sgID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRules in security group %s:\n", sgID)
	for _, rule := range rules {
		fmt.Printf("Rule ID: %s\n", rule.ID)
		fmt.Printf("  Direction: %s\n", rule.Direction)
		fmt.Printf("  Protocol: %s\n", rule.Protocol)
		if rule.PortRangeMin != nil {
			fmt.Printf("  Port Range: %d-%d\n", *rule.PortRangeMin, *rule.PortRangeMax)
		}
		fmt.Printf("  Remote IP Prefix: %s\n", rule.RemoteIPPrefix)
		fmt.Printf("  Description: %s\n", rule.Description)
		fmt.Printf("  Status: %s\n\n", rule.Status)
	}

	// Get details of a specific rule
	ruleDetails, err := networkClient.Rules().Get(ctx, sshRuleID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SSH Rule Details:\n")
	fmt.Printf("  ID: %s\n", ruleDetails.ID)
	fmt.Printf("  Direction: %s\n", ruleDetails.Direction)
	fmt.Printf("  Protocol: %s\n", ruleDetails.Protocol)
	fmt.Printf("  Port: %d\n", *ruleDetails.PortRangeMin)
	fmt.Printf("  Remote IP Prefix: %s\n", ruleDetails.RemoteIPPrefix)

	// Clean up - delete rules and security group
	fmt.Println("\nCleaning up resources...")

	// Delete rules
	if err := networkClient.Rules().Delete(ctx, sshRuleID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted SSH rule: %s\n", sshRuleID)

	if err := networkClient.Rules().Delete(ctx, httpsRuleID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted HTTPS rule: %s\n", httpsRuleID)

	// Delete security group
	if err := networkClient.SecurityGroups().Delete(ctx, sgID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted security group: %s\n", sgID)
}

func ExamplePublicIPs() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// List all public IPs
	publicIPs, err := networkClient.PublicIPs().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available Public IPs:")
	for _, pip := range publicIPs {
		fmt.Printf("  ID: %s\n", pip.ID)
		fmt.Printf("  Public IP: %s\n", pip.PublicIP)
		fmt.Printf("  VPC ID: %s\n", pip.VPCID)
		fmt.Printf("  Port ID: %s\n", pip.PortID)
		fmt.Printf("  Status: %s\n", pip.Status)
		fmt.Printf("  Created At: %s\n\n", pip.CreatedAt)
	}

	// Create a VPC and port for testing public IP operations
	vpcReq := network.CreateVPCRequest{
		Name:        "test-vpc-for-pip",
		Description: "VPC for testing public IP operations",
	}

	vpcID, err := networkClient.VPCs().Create(ctx, vpcReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created VPC with ID: %s\n", vpcID)

	// Create a port in the VPC
	portReq := network.PortCreateRequest{
		Name:   "test-port-for-pip",
		HasPIP: true,
		HasSG:  false,
	}

	portID, err := networkClient.VPCs().CreatePort(ctx, vpcID, portReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created port with ID: %s\n", portID)

	// Create a public IP
	pipReq := network.PublicIPCreateRequest{
		Description: "Test public IP created via SDK",
	}

	pipID, err := networkClient.VPCs().CreatePublicIP(ctx, vpcID, pipReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created public IP with ID: %s\n", pipID)

	// Get public IP details
	pip, err := networkClient.PublicIPs().Get(ctx, pipID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nPublic IP Details:\n")
	fmt.Printf("  ID: %s\n", pip.ID)
	fmt.Printf("  Public IP: %s\n", pip.PublicIP)
	fmt.Printf("  Description: %s\n", pip.Description)
	fmt.Printf("  Status: %s\n", pip.Status)

	// Attach public IP to port
	if err := networkClient.PublicIPs().AttachToPort(ctx, pipID, portID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Attached public IP %s to port %s\n", pipID, portID)

	// Get updated public IP details
	pip, err = networkClient.PublicIPs().Get(ctx, pipID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Public IP %s is now attached to port %s\n", pip.PublicIP, pip.PortID)

	// Detach public IP from port
	if err := networkClient.PublicIPs().DetachFromPort(ctx, pipID, portID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Detached public IP %s from port %s\n", pipID, portID)

	// Clean up resources
	fmt.Println("\nCleaning up resources...")

	// Delete public IP
	if err := networkClient.PublicIPs().Delete(ctx, pipID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted public IP: %s\n", pipID)

	// Delete VPC (this will also delete the port)
	if err := networkClient.VPCs().Delete(ctx, vpcID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted VPC: %s\n", vpcID)
}

func ExamplePorts() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	networkClient := network.New(c)
	ctx := context.Background()

	// List all ports
	ports, err := networkClient.Ports().List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available Ports:")
	for _, port := range ports {
		fmt.Printf("  ID: %s\n", port.ID)
		fmt.Printf("  Name: %s\n", port.Name)
		fmt.Printf("  VPC ID: %s\n", port.VPCID)
		fmt.Printf("  Security Groups: %v\n", port.SecurityGroups)
		fmt.Printf("  Created At: %s\n", port.CreatedAt)

		// Print IP addresses
		if len(port.IPAddress) > 0 {
			fmt.Println("  IP Addresses:")
			for _, ip := range port.IPAddress {
				fmt.Printf("    %s (Subnet: %s)\n", ip.IPAddress, ip.SubnetID)
			}
		}

		// Print public IPs if any
		if len(port.PublicIP) > 0 {
			fmt.Println("  Public IPs:")
			for _, pip := range port.PublicIP {
				fmt.Printf("    %s (ID: %s)\n", pip.PublicIP, pip.PublicIPID)
			}
		}
		fmt.Println()
	}

	// Create a security group for testing
	sgCreateReq := network.SecurityGroupCreateRequest{
		Name:        "test-sg-for-port",
		Description: "Security group for testing port operations",
	}

	sgID, err := networkClient.SecurityGroups().Create(ctx, sgCreateReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created security group with ID: %s\n", sgID)

	// Get port details by ID (using first port from list if available)
	if len(ports) > 0 {
		portDetails, err := networkClient.Ports().Get(ctx, ports[0].ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nPort Details:\n")
		fmt.Printf("  ID: %s\n", portDetails.ID)
		fmt.Printf("  Name: %s\n", portDetails.Name)
		fmt.Printf("  Description: %s\n", portDetails.Description)
		fmt.Printf("  VPC ID: %s\n", portDetails.VPCID)
		fmt.Printf("  Security Groups: %v\n", portDetails.SecurityGroups)

		// Attach security group to port
		if err := networkClient.Ports().AttachSecurityGroup(ctx, portDetails.ID, sgID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Attached security group %s to port %s\n", sgID, portDetails.ID)

		// Get updated port details
		portDetails, err = networkClient.Ports().Get(ctx, portDetails.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updated security groups: %v\n", portDetails.SecurityGroups)

		// Detach security group from port
		if err := networkClient.Ports().DetachSecurityGroup(ctx, portDetails.ID, sgID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Detached security group %s from port %s\n", sgID, portDetails.ID)
	}

	// Clean up resources
	fmt.Println("\nCleaning up resources...")

	// Delete security group
	if err := networkClient.SecurityGroups().Delete(ctx, sgID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted security group: %s\n", sgID)
}
