package compute

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

type (
	ListInstancesResponse struct {
		Instances []Instance `json:"instances"`
	}

	Instance struct {
		ID               string     `json:"id"`
		Name             string     `json:"name,omitempty"`
		MachineType      IDOrName   `json:"machine_type"`
		Image            IDOrName   `json:"image"`
		Status           string     `json:"status"`
		State            string     `json:"state"`
		CreatedAt        time.Time  `json:"created_at"`
		UpdatedAt        *time.Time `json:"updated_at,omitempty"`
		SSHKeyName       string     `json:"ssh_key_name,omitempty"`
		AvailabilityZone string     `json:"availability_zone,omitempty"`
	}

	CreateRequest struct {
		AvailabilityZone *string                  `json:"availability_zone,omitempty"`
		Image            IDOrName                 `json:"image"`
		Labels           *CreateParametersLabels  `json:"labels,omitempty"`
		MachineType      IDOrName                 `json:"machine_type"`
		Name             string                   `json:"name"`
		Network          *CreateParametersNetwork `json:"network,omitempty"`
		SshKeyName       *string                  `json:"ssh_key_name,omitempty"`
		UserData         *string                  `json:"user_data,omitempty"`
	}

	CreateParametersLabels struct {
		Values []string
	}

	CreateParametersNetwork struct {
		AssociatePublicIp *bool                             `json:"associate_public_ip,omitempty"`
		Interface         *CreateParametersNetworkInterface `json:"interface,omitempty"`
		Vpc               *CreateParametersNetworkVpc       `json:"vpc,omitempty"`
	}

	CreateParametersNetworkInterface struct {
		Interface      IDOrName                                        `json:"interface"`
		SecurityGroups *CreateParametersNetworkInterfaceSecurityGroups `json:"security_groups,omitempty"`
	}

	CreateParametersNetworkInterfaceSecurityGroupsItem struct {
		Id string `json:"id"`
	}

	CreateParametersNetworkInterfaceSecurityGroups struct {
		Items []CreateParametersNetworkInterfaceSecurityGroupsItem
	}

	CreateParametersNetworkVpc struct {
		Vpc            IDOrName                                        `json:"vpc"`
		SecurityGroups *CreateParametersNetworkInterfaceSecurityGroups `json:"security_groups,omitempty"`
	}

	IDOrName struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	UpdateNameRequest struct {
		Name string `json:"name"`
	}

	RetypeRequest struct {
		MachineType IDOrName `json:"machine_type"`
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
}

// List retrieves all instances
func (s *instanceService) List(ctx context.Context, opts ListOptions) ([]Instance, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "/v1/instances", nil)
	if err != nil {
		return nil, err
	}

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
	req.URL.RawQuery = q.Encode()

	var response struct {
		Instances []Instance `json:"instances"`
	}
	resp, err := s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}
	
	if resp == nil {
		return nil, fmt.Errorf("empty response")
	}
	
	r, ok := resp.(*struct {
		Instances []Instance `json:"instances"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", resp)
	}
	
	if r.Instances == nil {
		return nil, fmt.Errorf("invalid response format: missing instances field")
	}
	
	return r.Instances, nil
}

// Create creates a new instance
func (s *instanceService) Create(ctx context.Context, createReq CreateRequest) (string, error) {
	var result struct {
		ID string `json:"id"`
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, "/v1/instances", createReq)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return "", err
	}
	if r, ok := resp.(*struct {
		ID string `json:"id"`
	}); ok {
		return r.ID, nil
	}
	return "", fmt.Errorf("unexpected response type: %T", resp)
}

// Get retrieves a specific instance
func (s *instanceService) Get(ctx context.Context, id string, expand []string) (*Instance, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/instances/%s", id), nil)
	if err != nil {
		return nil, err
	}

	if len(expand) > 0 {
		q := req.URL.Query()
		q.Add("expand", strings.Join(expand, ","))
		req.URL.RawQuery = q.Encode()
	}

	var instance Instance
	resp, err := s.client.Do(ctx, req, &instance)
	if err != nil {
		return nil, err
	}
	if r, ok := resp.(*Instance); ok {
		return r, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", resp)
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

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for delete operation")
	}
	return nil
}

// Rename changes the instance name
func (s *instanceService) Rename(ctx context.Context, id string, newName string) error {
	if id == "" {
		return &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}

	req, err := s.client.newRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/v1/instances/%s/rename", id),
		UpdateNameRequest{Name: newName})
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for rename operation")
	}
	return nil
}

// Retype changes the instance machine type
func (s *instanceService) Retype(ctx context.Context, id string, retypeReq RetypeRequest) error {
	if id == "" {
		return &client.ValidationError{Field: "id", Message: "cannot be empty"}
	}

	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/instances/%s/retype", id),
		retypeReq)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for retype operation")
	}
	return nil
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

	req, err := s.client.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/v1/instances/%s/%s", id, action),
		nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		return fmt.Errorf("unexpected response data for %s operation", action)
	}
	return nil
}
