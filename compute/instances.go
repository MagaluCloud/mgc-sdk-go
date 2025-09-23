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
	InstanceImageExpand       = "image"
	InstanceMachineTypeExpand = "machine-type"
	InstanceNetworkExpand     = "network"
)

// Constants for API version headers.
const (
	VmInstanceHeaderVersionName = "x-api-version"
	VmInstanceHeaderVersion     = "1.1"
)

// ListInstancesResponse represents the response from listing instances.
type ListInstancesResponse struct {
	Instances []Instance `json:"instances"`
}

// InstanceTypes represents the machine type configuration of an instance.
type InstanceTypes struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Vcpus *int    `json:"vcpus"`
	Ram   *int    `json:"ram"`
	Disk  *int    `json:"disk"`
}

// VmImage represents the image configuration of an instance.
type VmImage struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Platform *string `json:"platform,omitempty,omitzero"`
}

// Instance represents a virtual machine instance.
type Instance struct {
	ID               string         `json:"id"`
	Name             *string        `json:"name,omitempty"`
	MachineType      *InstanceTypes `json:"machine_type"`
	Image            *VmImage       `json:"image"`
	Status           string         `json:"status"`
	State            string         `json:"state"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at,omitempty,omitzero"`
	SSHKeyName       *string        `json:"ssh_key_name,omitempty"`
	AvailabilityZone *string        `json:"availability_zone,omitempty,omitzero"`
	Network          *Network       `json:"network"`
	UserData         *string        `json:"user_data,omitempty"`
	Labels           *[]string      `json:"labels"`
	Error            *Error         `json:"error,omitempty"`
}

// Error represents an error that occurred with an instance.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// CreateRequest represents the request to create a new instance.
type CreateRequest struct {
	AvailabilityZone *string                  `json:"availability_zone,omitempty"`
	Image            IDOrName                 `json:"image"`
	Labels           *[]string                `json:"labels,omitempty"`
	MachineType      IDOrName                 `json:"machine_type"`
	Name             string                   `json:"name"`
	Network          *CreateParametersNetwork `json:"network,omitempty"`
	SshKeyName       *string                  `json:"ssh_key_name,omitempty"`
	UserData         *string                  `json:"user_data,omitempty"`
}

// CreateParametersNetwork represents network configuration for instance creation.
type CreateParametersNetwork struct {
	AssociatePublicIp *bool                             `json:"associate_public_ip,omitempty"`
	Interface         *CreateParametersNetworkInterface `json:"interface,omitempty"`
	Vpc               *IDOrName                         `json:"vpc,omitempty"`
}

// CreateParametersNetworkInterface represents network interface configuration.
type CreateParametersNetworkInterface struct {
	ID             *string                                   `json:"id,omitempty"`
	SecurityGroups *[]CreateParametersNetworkInterfaceWithID `json:"security_groups,omitempty"`
}

// CreateParametersNetworkInterfaceWithID represents a security group item.
type CreateParametersNetworkInterfaceWithID struct {
	ID string `json:"id"`
}

// IDOrName represents a resource that can be identified by ID or name.
type IDOrName struct {
	ID   *string `json:"id,omitempty,omitzero"`
	Name *string `json:"name,omitempty,omitzero"`
}

// UpdateNameRequest represents the request to update an instance name.
type UpdateNameRequest struct {
	Name string `json:"name"`
}

// RetypeRequest represents the request to change an instance's machine type.
type RetypeRequest struct {
	MachineType IDOrName `json:"machine_type"`
}

// WindowsPasswordResponse represents the response from getting Windows password.
type WindowsPasswordResponse struct {
	Instance WindowsPasswordInstance `json:"instance"`
}

// WindowsPasswordInstance represents Windows password information.
type WindowsPasswordInstance struct {
	ID        string    `json:"id"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	User      string    `json:"user,omitempty"`
}

// NICRequest represents the request to attach or detach a network interface.
type NICRequest struct {
	Instance IDOrName            `json:"instance"`
	Network  NICRequestInterface `json:"network"`
}

// NICRequestInterface represents network interface configuration for NIC operations.
type NICRequestInterface struct {
	Interface IDOrName `json:"interface"`
}

// IpAddressNewExpand represents IP address information for network interfaces.
type IpAddressNewExpand struct {
	PrivateIpv4 string `json:"private_ipv4"`
	PublicIpv6  string `json:"public_ipv6,omitempty"`
}

// NetworkInterface represents a network interface attached to an instance.
type NetworkInterface struct {
	ID                   string             `json:"id"`
	Name                 string             `json:"name"`
	SecurityGroups       *[]string          `json:"security_groups"`
	Primary              *bool              `json:"primary"`
	AssociatedPublicIpv4 *string            `json:"associated_public_ipv4,omitempty"`
	IpAddresses          IpAddressNewExpand `json:"ip_addresses"`
}

// Network represents the network configuration of an instance.
type Network struct {
	Vpc        *IDOrName           `json:"vpc,omitempty"`
	Interfaces *[]NetworkInterface `json:"interfaces,omitempty"`
}

// InitLogResponse represents the response from getting instance initialization logs.
type InitLogResponse struct {
	Logs []string `json:"logs"`
}

// InstanceService provides operations for managing virtual machine instances.
type InstanceService interface {
	List(ctx context.Context, opts ListOptions) ([]Instance, error)
	Create(ctx context.Context, req CreateRequest) (string, error)
	Get(ctx context.Context, id string, expand []string) (*Instance, error)
	Delete(ctx context.Context, id string, deletePublicIP bool) error
	Rename(ctx context.Context, id string, newName string) error
	Retype(ctx context.Context, id string, req RetypeRequest) error
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Suspend(ctx context.Context, id string) error
	GetFirstWindowsPassword(ctx context.Context, id string) (*WindowsPasswordResponse, error)
	AttachNetworkInterface(ctx context.Context, req NICRequest) error
	DetachNetworkInterface(ctx context.Context, req NICRequest) error
	InitLog(ctx context.Context, id string, maxLines *int) (*InitLogResponse, error)
}

// instanceService implements the InstanceService interface.
type instanceService struct {
	client *VirtualMachineClient
}

// ListOptions defines the parameters for filtering and pagination of instance lists.
type ListOptions struct {
	Limit  *int
	Offset *int
	Sort   *string
	Expand []string
	Name   *string
}

// List retrieves all instances.
// This method makes an HTTP request to get the list of instances
// and applies the filters specified in the options.
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
func (s *instanceService) Start(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "start")
}

// Stop stops the instance.
// This method makes an HTTP request to gracefully power off a running instance.
// Returns an error if the instance is already stopped or if the operation fails.
func (s *instanceService) Stop(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "stop")
}

// Suspend suspends the instance.
// This method makes an HTTP request to pause the execution of an instance
// while maintaining its state in memory.
// Returns an error if the instance cannot be suspended or if the operation fails.
func (s *instanceService) Suspend(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "suspend")
}

// executeInstanceAction handles common instance state change operations.
// This is an internal method that should not be called directly by SDK users.
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
