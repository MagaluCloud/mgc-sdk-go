package lbaas

import (
	"encoding/json"
	"strings"
)

// urlNetworkLoadBalancer constructs the URL path for network load balancer operations
// If lbID is provided, it appends the ID to the base path
// Additional path segments can be provided via extraPath parameter
func urlNetworkLoadBalancer(lbID *string, extraPath ...string) string {
	result := "/v0beta1/network-load-balancers"
	if lbID == nil || *lbID == "" {
		return result
	}
	result += "/" + *lbID
	if len(extraPath) > 0 {
		result += "/" + strings.Join(extraPath, "/")
	}
	return result
}

// UnmarshalJSON implements custom JSON unmarshaling for TargetsRawOrInstancesRequest
// This allows flexible handling of different target types in the JSON payload
func (t *TargetsRawOrInstancesRequest) UnmarshalJSON(data []byte) error {
	var targets any
	if err := json.Unmarshal(data, &targets); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for TargetsRawOrInstancesRequest
// Prioritizes TargetsInstances over TargetsRaw when marshaling
func (t *TargetsRawOrInstancesRequest) MarshalJSON() ([]byte, error) {
	if len(t.TargetsInstances) > 0 {
		return json.Marshal(t.TargetsInstances)
	}
	if len(t.TargetsRaw) > 0 {
		return json.Marshal(t.TargetsRaw)
	}
	return nil, nil
}

// UnmarshalJSON implements custom JSON unmarshaling for TargetsRawOrInstancesUpdateRequest
// This allows flexible handling of different target types in the JSON payload
func (t *TargetsRawOrInstancesUpdateRequest) UnmarshalJSON(data []byte) error {
	var targets any
	if err := json.Unmarshal(data, &targets); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for TargetsRawOrInstancesUpdateRequest
// Prioritizes TargetsInstances over TargetsRaw when marshaling
func (t *TargetsRawOrInstancesUpdateRequest) MarshalJSON() ([]byte, error) {
	if len(t.TargetsInstances) > 0 {
		return json.Marshal(t.TargetsInstances)
	}
	if len(t.TargetsRaw) > 0 {
		return json.Marshal(t.TargetsRaw)
	}
	return nil, nil
}
