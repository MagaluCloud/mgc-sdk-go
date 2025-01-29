package dbaas

import (
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		core     *client.CoreClient
		wantNil  bool
	}{
		{
			name: "valid core client",
			core: client.NewMgcClient("test-token"),
			wantNil: false,
		},
		{
			name: "nil core client",
			core: nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.core)
			if (client == nil) != tt.wantNil {
				t.Errorf("New() returned nil = %v, want nil = %v", client == nil, tt.wantNil)
			}
		})
	}
}

func TestDBaaSClient_Services(t *testing.T) {
	core := client.NewMgcClient("test-token",
		client.WithHTTPClient(&http.Client{}),
		client.WithBaseURL("https://api.test.com"))
	
	dbaas := New(core)

	t.Run("Engines service", func(t *testing.T) {
		service := dbaas.Engines()
		if service == nil {
			t.Error("Engines() returned nil")
		}
	})

	t.Run("InstanceTypes service", func(t *testing.T) {
		service := dbaas.InstanceTypes()
		if service == nil {
			t.Error("InstanceTypes() returned nil")
		}
	})

	t.Run("Instances service", func(t *testing.T) {
		service := dbaas.Instances()
		if service == nil {
			t.Error("Instances() returned nil")
		}
	})

	t.Run("Replicas service", func(t *testing.T) {
		service := dbaas.Replicas()
		if service == nil {
			t.Error("Replicas() returned nil")
		}
	})
}
