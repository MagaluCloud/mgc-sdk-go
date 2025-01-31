package kubernetes

import (
	"net/http"
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestNewKubernetesClient(t *testing.T) {
	tests := []struct {
		name    string
		core    *client.CoreClient
		wantNil bool
	}{
		{
			name:    "valid core client",
			core:    client.NewMgcClient("test-token"),
			wantNil: false,
		},
		{
			name:    "nil core client",
			core:    nil,
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

func TestKubernetesClient_Services(t *testing.T) {
	core := client.NewMgcClient("test-token",
		client.WithHTTPClient(&http.Client{}),
		client.WithBaseURL("https://api.test.com"))

	k8sClient := New(core)

	t.Run("Clusters service", func(t *testing.T) {
		service := k8sClient.Clusters()
		if service == nil {
			t.Error("Clusters() returned nil")
		}
	})

	t.Run("Flavors service", func(t *testing.T) {
		service := k8sClient.Flavors()
		if service == nil {
			t.Error("Flavors() returned nil")
		}
	})

	t.Run("Info service", func(t *testing.T) {
		service := k8sClient.Info()
		if service == nil {
			t.Error("Info() returned nil")
		}
	})

	t.Run("Nodepools service", func(t *testing.T) {
		service := k8sClient.Nodepools()
		if service == nil {
			t.Error("Nodepools() returned nil")
		}
	})

	t.Run("Versions service", func(t *testing.T) {
		service := k8sClient.Versions()
		if service == nil {
			t.Error("Versions() returned nil")
		}
	})
}
