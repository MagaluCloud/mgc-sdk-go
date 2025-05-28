package lbaas

import "encoding/json"

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
