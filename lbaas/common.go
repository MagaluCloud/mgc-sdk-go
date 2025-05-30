package lbaas

import (
	"encoding/json"
	"strings"
)

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

func (t *TargetsRawOrInstancesRequest) UnmarshalJSON(data []byte) error {
	var targets any
	if err := json.Unmarshal(data, &targets); err != nil {
		return err
	}
	return nil
}

func (t *TargetsRawOrInstancesRequest) MarshalJSON() ([]byte, error) {
	if len(t.TargetsInstances) > 0 {
		return json.Marshal(t.TargetsInstances)
	}
	if len(t.TargetsRaw) > 0 {
		return json.Marshal(t.TargetsRaw)
	}
	return nil, nil
}

func (t *TargetsRawOrInstancesUpdateRequest) UnmarshalJSON(data []byte) error {
	var targets any
	if err := json.Unmarshal(data, &targets); err != nil {
		return err
	}
	return nil
}

func (t *TargetsRawOrInstancesUpdateRequest) MarshalJSON() ([]byte, error) {
	if len(t.TargetsInstances) > 0 {
		return json.Marshal(t.TargetsInstances)
	}
	if len(t.TargetsRaw) > 0 {
		return json.Marshal(t.TargetsRaw)
	}
	return nil, nil
}
