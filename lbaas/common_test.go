package lbaas

import (
	"encoding/json"
	"testing"
)

func TestTargetsRawOrInstancesRequest_MarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		targets  TargetsRawOrInstancesRequest
		expected string
		wantErr  bool
	}{
		{
			name: "marshal targets instances",
			targets: TargetsRawOrInstancesRequest{
				TargetsInstances: []NetworkBackendInstanceRequest{
					{NicID: "nic-123", Port: 80},
					{NicID: "nic-456", Port: 443},
				},
				TargetsRaw: nil,
			},
			expected: `[{"nic_id":"nic-123","port":80},{"nic_id":"nic-456","port":443}]`,
			wantErr:  false,
		},
		{
			name: "marshal targets raw",
			targets: TargetsRawOrInstancesRequest{
				TargetsInstances: nil,
				TargetsRaw: []NetworkBackendRawTargetRequest{
					{IPAddress: "192.168.1.1", Port: 80},
					{IPAddress: "192.168.1.2", Port: 443},
				},
			},
			expected: `[{"ip_address":"192.168.1.1","port":80},{"ip_address":"192.168.1.2","port":443}]`,
			wantErr:  false,
		},
		{
			name: "marshal empty targets",
			targets: TargetsRawOrInstancesRequest{
				TargetsInstances: nil,
				TargetsRaw:       nil,
			},
			expected: ``,
			wantErr:  false,
		},
		{
			name: "marshal both targets instances and raw - instances takes priority",
			targets: TargetsRawOrInstancesRequest{
				TargetsInstances: []NetworkBackendInstanceRequest{
					{NicID: "nic-123", Port: 80},
				},
				TargetsRaw: []NetworkBackendRawTargetRequest{
					{IPAddress: "192.168.1.1", Port: 80},
				},
			},
			expected: `[{"nic_id":"nic-123","port":80}]`,
			wantErr:  false,
		},
		{
			name: "marshal empty slices",
			targets: TargetsRawOrInstancesRequest{
				TargetsInstances: []NetworkBackendInstanceRequest{},
				TargetsRaw:       []NetworkBackendRawTargetRequest{},
			},
			expected: ``,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.targets.MarshalJSON()

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.expected, string(result))
		})
	}
}

func TestTargetsRawOrInstancesRequest_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "unmarshal valid json array",
			data:    `[{"nic_id":"nic-123","port":80}]`,
			wantErr: false,
		},
		{
			name:    "unmarshal valid json object",
			data:    `{"nic_id":"nic-123","port":80}`,
			wantErr: false,
		},
		{
			name:    "unmarshal null",
			data:    `null`,
			wantErr: false,
		},
		{
			name:    "unmarshal empty array",
			data:    `[]`,
			wantErr: false,
		},
		{
			name:    "unmarshal string",
			data:    `"test"`,
			wantErr: false,
		},
		{
			name:    "unmarshal number",
			data:    `123`,
			wantErr: false,
		},
		{
			name:    "unmarshal invalid json",
			data:    `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var targets TargetsRawOrInstancesRequest
			err := targets.UnmarshalJSON([]byte(tt.data))

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

func TestTargetsRawOrInstancesUpdateRequest_MarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		targets  TargetsRawOrInstancesUpdateRequest
		expected string
		wantErr  bool
	}{
		{
			name: "marshal targets instances",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsInstances: []NetworkBackendInstanceUpdateRequest{
					{NicID: "nic-123", Port: 80},
					{NicID: "nic-456", Port: 443},
				},
				TargetsRaw: nil,
			},
			expected: `[{"nic_id":"nic-123","port":80},{"nic_id":"nic-456","port":443}]`,
			wantErr:  false,
		},
		{
			name: "marshal targets raw",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsInstances: nil,
				TargetsRaw: []NetworkBackendRawTargetUpdateRequest{
					{IPAddress: "192.168.1.1", Port: 80},
					{IPAddress: "192.168.1.2", Port: 443},
				},
			},
			expected: `[{"ip_address":"192.168.1.1","port":80},{"ip_address":"192.168.1.2","port":443}]`,
			wantErr:  false,
		},
		{
			name: "marshal empty targets",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsInstances: nil,
				TargetsRaw:       nil,
			},
			expected: ``,
			wantErr:  false,
		},
		{
			name: "marshal both targets instances and raw - instances takes priority",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsInstances: []NetworkBackendInstanceUpdateRequest{
					{NicID: "nic-123", Port: 80},
				},
				TargetsRaw: []NetworkBackendRawTargetUpdateRequest{
					{IPAddress: "192.168.1.1", Port: 80},
				},
			},
			expected: `[{"nic_id":"nic-123","port":80}]`,
			wantErr:  false,
		},
		{
			name: "marshal empty slices",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsInstances: []NetworkBackendInstanceUpdateRequest{},
				TargetsRaw:       []NetworkBackendRawTargetUpdateRequest{},
			},
			expected: ``,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.targets.MarshalJSON()

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.expected, string(result))
		})
	}
}

func TestTargetsRawOrInstancesUpdateRequest_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "unmarshal valid json array",
			data:    `[{"nic_id":"nic-123","port":80}]`,
			wantErr: false,
		},
		{
			name:    "unmarshal valid json object",
			data:    `{"nic_id":"nic-123","port":80}`,
			wantErr: false,
		},
		{
			name:    "unmarshal null",
			data:    `null`,
			wantErr: false,
		},
		{
			name:    "unmarshal empty array",
			data:    `[]`,
			wantErr: false,
		},
		{
			name:    "unmarshal string",
			data:    `"test"`,
			wantErr: false,
		},
		{
			name:    "unmarshal number",
			data:    `123`,
			wantErr: false,
		},
		{
			name:    "unmarshal invalid json",
			data:    `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var targets TargetsRawOrInstancesUpdateRequest
			err := targets.UnmarshalJSON([]byte(tt.data))

			if tt.wantErr {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
		})
	}
}

// Testes de integração para verificar se o marshal/unmarshal funciona corretamente em conjunto
func TestTargetsRawOrInstancesRequest_MarshalUnmarshal_Integration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		targets TargetsRawOrInstancesRequest
	}{
		{
			name: "round trip with instances",
			targets: TargetsRawOrInstancesRequest{
				TargetsInstances: []NetworkBackendInstanceRequest{
					{NicID: "nic-123", Port: 80},
					{NicID: "nic-456", Port: 443},
				},
			},
		},
		{
			name: "round trip with raw targets",
			targets: TargetsRawOrInstancesRequest{
				TargetsRaw: []NetworkBackendRawTargetRequest{
					{IPAddress: "192.168.1.1", Port: 80},
					{IPAddress: "192.168.1.2", Port: 443},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Marshal
			data, err := json.Marshal(tt.targets)
			assertNoError(t, err)

			// Unmarshal
			var result TargetsRawOrInstancesRequest
			err = json.Unmarshal(data, &result)
			assertNoError(t, err)
		})
	}
}

func TestTargetsRawOrInstancesUpdateRequest_MarshalUnmarshal_Integration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		targets TargetsRawOrInstancesUpdateRequest
	}{
		{
			name: "round trip with instances",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsInstances: []NetworkBackendInstanceUpdateRequest{
					{NicID: "nic-123", Port: 80},
					{NicID: "nic-456", Port: 443},
				},
			},
		},
		{
			name: "round trip with raw targets",
			targets: TargetsRawOrInstancesUpdateRequest{
				TargetsRaw: []NetworkBackendRawTargetUpdateRequest{
					{IPAddress: "192.168.1.1", Port: 80},
					{IPAddress: "192.168.1.2", Port: 443},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Marshal
			data, err := json.Marshal(tt.targets)
			assertNoError(t, err)

			// Unmarshal
			var result TargetsRawOrInstancesUpdateRequest
			err = json.Unmarshal(data, &result)
			assertNoError(t, err)
		})
	}
}
