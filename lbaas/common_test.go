package lbaas

import (
	"testing"
)

func TestUrlNetworkLoadBalancer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		lbID      *string
		extraPath []string
		expected  string
	}{
		{
			name:      "nil lbID without extra path",
			lbID:      nil,
			extraPath: nil,
			expected:  "/v0beta1/network-load-balancers",
		},
		{
			name:      "nil lbID with extra path",
			lbID:      nil,
			extraPath: []string{"listeners", "123"},
			expected:  "/v0beta1/network-load-balancers",
		},
		{
			name:      "valid lbID without extra path",
			lbID:      new("lb-123"),
			extraPath: nil,
			expected:  "/v0beta1/network-load-balancers/lb-123",
		},
		{
			name:      "valid lbID with single extra path",
			lbID:      new("lb-456"),
			extraPath: []string{"listeners"},
			expected:  "/v0beta1/network-load-balancers/lb-456/listeners",
		},
		{
			name:      "valid lbID with multiple extra paths",
			lbID:      new("lb-789"),
			extraPath: []string{"backends", "backend-123", "targets"},
			expected:  "/v0beta1/network-load-balancers/lb-789/backends/backend-123/targets",
		},
		{
			name:      "empty lbID without extra path",
			lbID:      new(""),
			extraPath: nil,
			expected:  "/v0beta1/network-load-balancers",
		},
		{
			name:      "empty lbID with extra path",
			lbID:      new(""),
			extraPath: []string{"health-checks"},
			expected:  "/v0beta1/network-load-balancers",
		},
		{
			name:      "valid lbID with empty extra path slice",
			lbID:      new("lb-abc"),
			extraPath: []string{},
			expected:  "/v0beta1/network-load-balancers/lb-abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := urlNetworkLoadBalancer(tt.lbID, tt.extraPath...)
			assertEqual(t, tt.expected, result)
		})
	}
}
