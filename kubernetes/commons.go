package kubernetes

import (
	"time"
)

type (
	// ListOptions provides options for listing resources
	ListOptions struct {
		Limit  *int
		Offset *int
		Sort   *string
		Expand []string
	}

	// MessageState represents a status message
	MessageState struct {
		State   string `json:"state"`
		Message string `json:"message"`
	}

	// Flavor represents a Kubernetes flavor (instance type)
	Flavor struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		VCPU int    `json:"vcpu"`
		RAM  int    `json:"ram"`
		Size int    `json:"size"`
	}

	// Status represents a status with messages
	Status struct {
		State    string   `json:"state"`
		Messages []string `json:"messages,omitempty"`
	}

	// Taint represents a node taint
	Taint struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Effect string `json:"effect"`
	}

	// AutoScale represents autoscaling configuration
	AutoScale struct {
		MinReplicas *int `json:"min_replicas"`
		MaxReplicas *int `json:"max_replicas"`
	}

	// InstanceTemplate represents the template for node instances
	InstanceTemplate struct {
		Flavor    Flavor `json:"flavor"`
		NodeImage string `json:"node_image"`
		DiskSize  int    `json:"disk_size"`
		DiskType  string `json:"disk_type"`
	}

	// NodePool represents a Kubernetes node pool
	NodePool struct {
		ID                string            `json:"id"`
		Name              string            `json:"name"`
		InstanceTemplate  InstanceTemplate  `json:"instance_template"`
		Replicas          int               `json:"replicas"`
		Labels            map[string]string `json:"labels,omitempty"`
		Taints            *[]Taint          `json:"taints,omitempty"`
		SecurityGroups    *[]string         `json:"security_groups,omitempty"`
		CreatedAt         *time.Time        `json:"created_at"`
		UpdatedAt         *time.Time        `json:"updated_at,omitempty"`
		AutoScale         *AutoScale        `json:"auto_scale,omitempty"`
		Status            Status            `json:"status"`
		Flavor            string            `json:"flavor"`
		MaxPodsPerNode    *int              `json:"max_pods_per_node,omitempty"`
		AvailabilityZones *[]string         `json:"availability_zones,omitempty"`
	}

	// CreateNodePoolRequest represents the request payload for creating a node pool
	CreateNodePoolRequest struct {
		Name              string     `json:"name"`
		Flavor            string     `json:"flavor"`
		Replicas          int        `json:"replicas"`
		Tags              *[]string  `json:"tags,omitempty"`
		Taints            *[]Taint   `json:"taints,omitempty"`
		AutoScale         *AutoScale `json:"auto_scale,omitempty"`
		MaxPodsPerNode    *int       `json:"max_pods_per_node,omitempty"`
		AvailabilityZones *[]string  `json:"availability_zones,omitempty"`
	}
)
