// Package compute provides functionality to interact with the MagaluCloud compute service.
// This package allows managing virtual machine instances, images, instance types, and snapshots.
package compute

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// Constants for expanding related resources in instance responses.
const (
	// InstanceImageExpand is used to include image information in instance responses
	InstanceImageExpand = "image"
	// InstanceMachineTypeExpand is used to include machine type information in instance responses
	InstanceMachineTypeExpand = "machine-type"
	// InstanceNetworkExpand is used to include network information in instance responses
	InstanceNetworkExpand = "network"
)

// Constants for API version headers.
const (
	// VmInstanceHeaderVersionName is the header name for API version
	VmInstanceHeaderVersionName = "x-api-version"
	// VmInstanceHeaderVersion is the API version value
	VmInstanceHeaderVersion = "1.1"
)

// ListInstancesResponse represents the response from listing instances.
// This structure encapsulates the API response format for instances.
type ListInstancesResponse struct {
	// Instances contains the list of instances
	Instances []Instance `json:"instances"`
}

// InstanceTypes represents the machine type configuration of an instance.
type InstanceTypes struct {
	// ID is the unique identifier of the machine type
	ID string `json:"id"`
	// Name is the display name of the machine type
	Name *string `json:"name"`
	// Vcpus is the number of virtual CPUs
	Vcpus *int `json:"vcpus"`
	// Ram is the amount of RAM in MB
	Ram *int `json:"ram"`
	// Disk is the disk size in GB
	Disk *int `json:"disk"`
}

// VmImage represents the image configuration of an instance.
type VmImage struct {
	// ID is the unique identifier of the image
	ID string `json:"id"`
	// Name is the display name of the image
	Name *string `json:"name"`
	// Platform indicates the operating system platform
	Platform *string `json:"platform,omitempty,omitzero"`
}

// Instance represents a virtual machine instance.
// An instance is a running virtual machine with its configuration and state.
type Instance struct {
	// ID is the unique identifier of the instance
	ID string `json:"id"`
	// Name is the display name of the instance
	Name *string `json:"name,omitempty"`
	// MachineType contains the machine type configuration
	MachineType *InstanceTypes `json:"machine_type"`
	// Image contains the image configuration
	Image *VmImage `json:"image"`
	// Status indicates the current status of the instance
	Status string `json:"status"`
	// State indicates the current state of the instance
	State string `json:"state"`
	// CreatedAt is the timestamp when the instance was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the instance was last updated
	UpdatedAt *time.Time `json:"updated_at,omitempty,omitzero"`
	// SSHKeyName is the name of the SSH key associated with the instance
	SSHKeyName *string `json:"ssh_key_name,omitempty"`
	// AvailabilityZone is the availability zone where the instance is located
	AvailabilityZone *string `json:"availability_zone,omitempty,omitzero"`
	// Network contains the network configuration
	Network *Network `json:"network"`
	// UserData contains user data passed to the instance
	UserData *string `json:"user_data,omitempty"`
	// Labels contains tags associated with the instance
	Labels *[]string `json:"labels"`
	// Error contains error information if the instance is in an error state
	Error *Error `json:"error,omitempty"`
}

// Error represents an error that occurred with an instance.
type Error struct {
	// Message is the error message
	Message string `json:"message"`
	// Slug is the error identifier
	Slug string `json:"slug"`
}

// CreateRequest represents the request to create a new instance.
type CreateRequest struct {
	// AvailabilityZone specifies the availability zone for the instance (optional)
	AvailabilityZone *string `json:"availability_zone,omitempty"`
	// Image specifies the image to use for the instance
	Image IDOrName `json:"image"`
	// Labels contains tags to associate with the instance (optional)
	Labels *[]string `json:"labels,omitempty"`
	// MachineType specifies the machine type for the instance
	MachineType IDOrName `json:"machine_type"`
	// Name is the display name for the instance
	Name string `json:"name"`
	// Network specifies the network configuration (optional)
	Network *CreateParametersNetwork `json:"network,omitempty"`
	// SshKeyName specifies the SSH key to use (optional)
	SshKeyName *string `json:"ssh_key_name,omitempty"`
	// UserData contains user data to pass to the instance (optional)
	UserData *string `json:"user_data,omitempty"`
}

// CreateParametersNetwork represents network configuration for instance creation.
type CreateParametersNetwork struct {
	// AssociatePublicIp specifies whether to associate a public IP (optional)
	AssociatePublicIp *bool `json:"associate_public_ip,omitempty"`
	// Interface specifies the network interface configuration (optional)
	Interface *CreateParametersNetworkInterface `json:"interface,omitempty"`
	// Vpc specifies the VPC to use (optional)
	Vpc *IDOrName `json:"vpc,omitempty"`
}

// CreateParametersNetworkInterface represents network interface configuration.
type CreateParametersNetworkInterface struct {
	// Interface specifies the network interface (optional)
	Interface *IDOrName `json:"interface,omitempty"`
	// SecurityGroups specifies the security groups to associate (optional)
	SecurityGroups *[]CreateParametersNetworkInterfaceSecurityGroupsItem `json:"security_groups,omitempty"`
}

// CreateParametersNetworkInterfaceSecurityGroupsItem represents a security group item.
type CreateParametersNetworkInterfaceSecurityGroupsItem struct {
	// Id is the security group identifier
	Id string `json:"id"`
}

// IDOrName represents a resource that can be identified by ID or name.
type IDOrName struct {
	// ID is the unique identifier (optional)
	ID *string `json:"id,omitempty,omitzero"`
	// Name is the display name (optional)
	Name *string `json:"name,omitempty,omitzero"`
}

// UpdateNameRequest represents the request to update an instance name.
type UpdateNameRequest struct {
	// Name is the new display name
	Name string `json:"name"`
}

// RetypeRequest represents the request to change an instance's machine type.
type RetypeRequest struct {
	// MachineType specifies the new machine type
	MachineType IDOrName `json:"machine_type"`
}

// WindowsPasswordResponse represents the response from getting Windows password.
type WindowsPasswordResponse struct {
	// Instance contains the password information
	Instance WindowsPasswordInstance `json:"instance"`
}

// WindowsPasswordInstance represents Windows password information.
type WindowsPasswordInstance struct {
	// ID is the instance identifier
	ID string `json:"id"`
	// Password is the Windows administrator password
	Password string `json:"password"`
	// CreatedAt is the timestamp when the password was created
	CreatedAt time.Time `json:"created_at"`
	// User is the username for the password (optional)
	User string `json:"user,omitempty"`
}

// NICRequest represents the request to attach or detach a network interface.
type NICRequest struct {
	// Instance specifies the instance to operate on
	Instance IDOrName `json:"instance"`
	// Network specifies the network interface configuration
	Network NICRequestInterface `json:"network"`
}

// NICRequestInterface represents network interface configuration for NIC operations.
type NICRequestInterface struct {
	// Interface specifies the network interface
	Interface IDOrName `json:"interface"`
}

// IpAddressNewExpand represents IP address information for network interfaces.
type IpAddressNewExpand struct {
	// PrivateIpv4 is the private IPv4 address
	PrivateIpv4 string `json:"private_ipv4"`
	// PublicIpv6 is the public IPv6 address (optional)
	PublicIpv6 string `json:"public_ipv6,omitempty"`
}

// NetworkInterface represents a network interface attached to an instance.
type NetworkInterface struct {
	// ID is the unique identifier of the network interface
	ID string `json:"id"`
	// Name is the display name of the network interface
	Name string `json:"name"`
	// SecurityGroups contains the associated security groups
	SecurityGroups *[]string `json:"security_groups"`
	// Primary indicates whether this is the primary network interface
	Primary *bool `json:"primary"`
	// AssociatedPublicIpv4 is the associated public IPv4 address (optional)
	AssociatedPublicIpv4 *string `json:"associated_public_ipv4,omitempty"`
	// IpAddresses contains the IP address configuration
	IpAddresses IpAddressNewExpand `json:"ip_addresses"`
}

// Network represents the network configuration of an instance.
type Network struct {
	// Vpc specifies the VPC (optional)
	Vpc *IDOrName `json:"vpc,omitempty"`
	// Interfaces contains the network interfaces (optional)
	Interfaces *[]NetworkInterface `json:"interfaces,omitempty"`
}

// InitLogResponse represents the response from getting instance initialization logs.
type InitLogResponse struct {
	// Logs contains the log lines
	Logs []string `json:"logs"`
}

// InstanceService provides operations for managing virtual machine instances.
// This interface allows creating, listing, retrieving, and managing instances.
type InstanceService interface {
	// List returns a slice of instances based on the provided listing options.
	// Use ListOptions to control pagination, sorting, and expansion of related resources.
	List(ctx context.Context, opts ListOptions) ([]Instance, error)

	// Create provisions a new virtual machine instance with the specified configuration.
	// Returns the ID of the newly created instance.
	Create(ctx context.Context, req CreateRequest) (string, error)

	// Get retrieves detailed information about a specific instance.
	// The expand parameter allows fetching related resources in the same request.
	Get(ctx context.Context, id string, expand []string) (*Instance, error)

	// Delete terminates and removes an instance.
	// If deletePublicIP is true, any associated public IP will also be released.
	Delete(ctx context.Context, id string, deletePublicIP bool) error

	// Rename updates the display name of an existing instance.
	// Returns an error if the operation fails or if the ID is empty.
	Rename(ctx context.Context, id string, newName string) error

	// Retype changes the machine type (size) of an instance.
	// The instance must be in a stopped state for this operation to succeed.
	Retype(ctx context.Context, id string, req RetypeRequest) error

	// Start powers on a stopped instance.
	// Returns an error if the instance is already running or if the operation fails.
	Start(ctx context.Context, id string) error

	// Stop gracefully powers off a running instance.
	// Returns an error if the instance is already stopped or if the operation fails.
	Stop(ctx context.Context, id string) error

	// Suspend pauses the execution of an instance while maintaining its state in memory.
	// Returns an error if the instance cannot be suspended or if the operation fails.
	Suspend(ctx context.Context, id string) error

	// GetFirstWindowsPassword retrieves the initial Windows administrator password for an instance
	GetFirstWindowsPassword(ctx context.Context, id string) (*WindowsPasswordResponse, error)

	// AttachNetworkInterface connects a network interface to an instance
	AttachNetworkInterface(ctx context.Context, req NICRequest) error

	// DetachNetworkInterface removes a non-primary network interface from an instance
	DetachNetworkInterface(ctx context.Context, req NICRequest) error

	// Retrieve instance init log output
	InitLog(ctx context.Context, id string, maxLines *int) (*InitLogResponse, error)
}

// instanceService implements the InstanceService interface.
// This is an internal implementation that should not be used directly.
type instanceService struct {
	client *VirtualMachineClient
}

// ListOptions defines the parameters for filtering and pagination of instance lists.
// All fields are optional and allow controlling the listing behavior.
type ListOptions struct {
	// Limit specifies the maximum number of results to return (1-1000)
	Limit *int
	// Offset specifies the number of results to skip for pagination
	Offset *int
	// Sort defines the field and direction for result ordering (e.g., "name:asc")
	Sort *string
	// Expand lists related resources to include in the response
	Expand []string
	// Name filters listed resources based on name field
	Name *string
}

// List retrieves all instances.
// This method makes an HTTP request to get the list of instances
// and applies the filters specified in the options.
//
// Parameters:
//   - ctx: Request context
//   - opts: Options to control pagination, sorting, and expansion
//
// Returns:
//   - []Instance: List of instances
//   - error: Error if there's a failure in the request
func (s *instanceService) List(ctx context.Context, opts ListOptions) ([]Instance, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/instances", nil)
	if err != nil {
		return nil, err
	}

	// Set API version header
	req.Header.Set(VmInstanceHeaderVersionName, VmInstanceHeaderVersion)

	q := req.URL.Query()
	if opts.Limit != nil {
		q.Add("_limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Add("_offset", strconv.Itoa(*opts.Offset))
	}
	if opts.Sort != nil {
		q.Add("_sort", *opts.Sort)
	}
	if len(opts.Expand) > 0 {
		q.Add("expand", strings.Join(opts.Expand, ","))
	}
	if opts.Name != nil {
		q.Add("name", *opts.Name)
	}

	req.URL.RawQuery = q.Encode()

	var response struct {
		Instances []Instance `json:"instances"`
	}

	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}
	return resp.Instances, nil
}

// Create creates a new instance.
// This method makes an HTTP request to provision a new virtual machine instance
// and returns the ID of the created instance.
//
// Parameters:
//   - ctx: Request context
//   - createReq: Request containing instance creation parameters
//
// Returns:
//   - string: ID of the created instance
//   - error: Error if there's a failure in the request
func (s *instanceService) Create(ctx context.Context, createReq CreateRequest) (string, error) {
	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[struct{ ID string }](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/instances",
		createReq,
		nil,
	)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

// Get retrieves a specific instance.
// This method makes an HTTP request to get detailed information about an instance
// and optionally expands related resources.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to retrieve
//   - expand: List of related resources to expand in the response
//
// Returns:
//   - *Instance: The requested instance
//   - error: Error if there's a failure in the request
func (s *instanceService) Get(ctx context.Context, id string, expand []string) (*Instance, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/instances/%s", id), nil)
	if err != nil {
		return nil, err
	}

	// Set API version header
	req.Header.Set(VmInstanceHeaderVersionName, VmInstanceHeaderVersion)

	if len(expand) > 0 {
		q := req.URL.Query()
		q.Add("expand", strings.Join(expand, ","))
		req.URL.RawQuery = q.Encode()
	}

	var instance Instance
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &instance)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete removes an instance.
// This method makes an HTTP request to terminate and remove an instance.
// If deletePublicIP is true, any associated public IP will also be released.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to delete
//   - deletePublicIP: Whether to also delete associated public IP
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) Delete(ctx context.Context, id string, deletePublicIP bool) error {
	req, err := s.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/instances/%s", id), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("delete_public_ip", strconv.FormatBool(deletePublicIP))
	req.URL.RawQuery = q.Encode()

	_, err = mgc_http.Do[any](s.client.GetConfig(), ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// Rename changes the instance name.
// This method makes an HTTP request to update the display name of an existing instance.
// Returns an error if the operation fails or if the ID is empty.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to rename
//   - newName: New display name for the instance
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) Rename(ctx context.Context, id string, newName string) error {
	if id == "" {
		return &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}
	path := fmt.Sprintf("/v1/instances/%s/rename", id)
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPatch,
		path,
		UpdateNameRequest{Name: newName},
		nil,
	)
}

// Retype changes the instance machine type.
// This method makes an HTTP request to change the machine type (size) of an instance.
// The instance must be in a stopped state for this operation to succeed.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to retype
//   - retypeReq: Request containing the new machine type
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) Retype(ctx context.Context, id string, retypeReq RetypeRequest) error {
	if id == "" {
		return &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}
	path := fmt.Sprintf("/v1/instances/%s/retype", id)
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		path,
		retypeReq,
		nil,
	)
}

// Start starts the instance.
// This method makes an HTTP request to power on a stopped instance.
// Returns an error if the instance is already running or if the operation fails.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to start
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) Start(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "start")
}

// Stop stops the instance.
// This method makes an HTTP request to gracefully power off a running instance.
// Returns an error if the instance is already stopped or if the operation fails.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to stop
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) Stop(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "stop")
}

// Suspend suspends the instance.
// This method makes an HTTP request to pause the execution of an instance
// while maintaining its state in memory.
// Returns an error if the instance cannot be suspended or if the operation fails.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to suspend
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) Suspend(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "suspend")
}

// executeInstanceAction handles common instance state change operations.
// This is an internal method that should not be called directly by SDK users.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to operate on
//   - action: The action to perform (start, stop, suspend)
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) executeInstanceAction(ctx context.Context, id string, action string) error {
	if id == "" {
		return &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}
	path := fmt.Sprintf("/v1/instances/%s/%s", id, action)
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		path,
		nil,
		nil,
	)
}

// GetFirstWindowsPassword retrieves the initial Windows administrator password for an instance.
// This method makes an HTTP request to get the Windows password for instances
// that were created with Windows images.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to get the password for
//
// Returns:
//   - *WindowsPasswordResponse: The password information
//   - error: Error if there's a failure in the request
func (s *instanceService) GetFirstWindowsPassword(ctx context.Context, id string) (*WindowsPasswordResponse, error) {
	if id == "" {
		return nil, &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}
	path := fmt.Sprintf("/v1/instances/config/%s/first-windows-password", id)
	result, err := mgc_http.ExecuteSimpleRequestWithRespBody[WindowsPasswordResponse](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodGet,
		path,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AttachNetworkInterface connects a network interface to an instance.
// This method makes an HTTP request to attach a network interface to an instance.
//
// Parameters:
//   - ctx: Request context
//   - req: Request containing the instance and network interface information
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) AttachNetworkInterface(ctx context.Context, req NICRequest) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/instances/network-interface/attach",
		req,
		nil,
	)
}

// DetachNetworkInterface removes a non-primary network interface from an instance.
// This method makes an HTTP request to detach a network interface from an instance.
//
// Parameters:
//   - ctx: Request context
//   - req: Request containing the instance and network interface information
//
// Returns:
//   - error: Error if there's a failure in the request
func (s *instanceService) DetachNetworkInterface(ctx context.Context, req NICRequest) error {
	return mgc_http.ExecuteSimpleRequest(
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/instances/network-interface/detach",
		req,
		nil,
	)
}

// InitLog retrieves instance initialization log output.
// This method makes an HTTP request to get the initialization logs for an instance.
//
// Parameters:
//   - ctx: Request context
//   - id: ID of the instance to get logs for
//   - maxLines: Maximum number of log lines to return (optional)
//
// Returns:
//   - *InitLogResponse: The log output
//   - error: Error if there's a failure in the request
func (s *instanceService) InitLog(ctx context.Context, id string, maxLines *int) (*InitLogResponse, error) {
	if id == "" {
		return nil, &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/instances/%s/init-logs", id), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if maxLines != nil {
		q.Add("max-lines-count", strconv.Itoa(*maxLines))
	}
	req.URL.RawQuery = q.Encode()

	var response InitLogResponse
	resp, err := mgc_http.Do(s.client.GetConfig(), ctx, req, &response)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
