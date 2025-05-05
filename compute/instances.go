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

const (
	InstanceImageExpand       = "image"
	InstanceMachineTypeExpand = "machine-type"
	InstanceNetworkExpand     = "network"
)

const (
	VmInstanceHeaderVersionName = "x-api-version"
	VmInstanceHeaderVersion     = "1.1"
)

type (
	ListInstancesResponse struct {
		Instances []Instance `json:"instances"`
	}

	InstanceTypes struct {
		ID    string  `json:"id"`
		Name  *string `json:"name"`
		Vcpus *int    `json:"vcpus"`
		Ram   *int    `json:"ram"`
		Disk  *int    `json:"disk"`
	}

	VmImage struct {
		ID       string  `json:"id"`
		Name     *string `json:"name"`
		Platform *string `json:"platform,omitempty,omitzero"`
	}

	Instance struct {
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

	Error struct {
		Message string `json:"message"`
		Slug    string `json:"slug"`
	}

	CreateRequest struct {
		AvailabilityZone *string                  `json:"availability_zone,omitempty"`
		Image            IDOrName                 `json:"image"`
		Labels           *[]string                `json:"labels,omitempty"`
		MachineType      IDOrName                 `json:"machine_type"`
		Name             string                   `json:"name"`
		Network          *CreateParametersNetwork `json:"network,omitempty"`
		SshKeyName       *string                  `json:"ssh_key_name,omitempty"`
		UserData         *string                  `json:"user_data,omitempty"`
	}

	CreateParametersNetwork struct {
		AssociatePublicIp *bool                             `json:"associate_public_ip,omitempty"`
		Interface         *CreateParametersNetworkInterface `json:"interface,omitempty"`
		Vpc               *IDOrName                         `json:"vpc,omitempty"`
	}

	CreateParametersNetworkInterface struct {
		Interface      IDOrName                                              `json:"interface"`
		SecurityGroups *[]CreateParametersNetworkInterfaceSecurityGroupsItem `json:"security_groups,omitempty"`
	}

	CreateParametersNetworkInterfaceSecurityGroupsItem struct {
		Id string `json:"id"`
	}

	IDOrName struct {
		ID   *string `json:"id,omitempty,omitzero"`
		Name *string `json:"name,omitempty,omitzero"`
	}

	UpdateNameRequest struct {
		Name string `json:"name"`
	}

	RetypeRequest struct {
		MachineType IDOrName `json:"machine_type"`
	}

	WindowsPasswordResponse struct {
		Instance WindowsPasswordInstance `json:"instance"`
	}

	WindowsPasswordInstance struct {
		ID        string    `json:"id"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
		User      string    `json:"user,omitempty"`
	}

	NICRequest struct {
		Instance IDOrName            `json:"instance"`
		Network  NICRequestInterface `json:"network"`
	}

	NICRequestInterface struct {
		Interface IDOrName `json:"interface"`
	}

	IpAddressNewExpand struct {
		PrivateIpv4 string `json:"private_ipv4"`
		PublicIpv6  string `json:"public_ipv6,omitempty"`
	}

	NetworkInterface struct {
		ID                   string             `json:"id"`
		Name                 string             `json:"name"`
		SecurityGroups       *[]string          `json:"security_groups"`
		Primary              *bool              `json:"primary"`
		AssociatedPublicIpv4 *string            `json:"associated_public_ipv4,omitempty"`
		IpAddresses          IpAddressNewExpand `json:"ip_addresses"`
	}

	Network struct {
		Vpc        *IDOrName           `json:"vpc,omitempty"`
		Interfaces *[]NetworkInterface `json:"interfaces,omitempty"`
	}
)

// InstanceService provides operations for managing virtual machine instances
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
}

type instanceService struct {
	client *VirtualMachineClient
}

// ListOptions defines the parameters for filtering and pagination of instance lists
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

// List retrieves all instances
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

// Create creates a new instance
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

// Get retrieves a specific instance
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

// Delete removes an instance
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

// Rename changes the instance name
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

// Retype changes the instance machine type
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

// Start the instance
func (s *instanceService) Start(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "start")
}

// Stop the instance
func (s *instanceService) Stop(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "stop")
}

// Suspend the instance
func (s *instanceService) Suspend(ctx context.Context, id string) error {
	return s.executeInstanceAction(ctx, id, "suspend")
}

// executeInstanceAction handles common instance state change operations
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
