package virtualmachine

import "time"

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
