package kubernetes

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

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
			core:    client.NewMgcClient(client.WithAPIKey("test-api-key")),
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
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
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

func TestKubernetesClient_NewRequest(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL("https://api.test.com"),
		client.WithHTTPClient(&http.Client{}),
	)
	k8sClient := New(core)

	t.Run("create valid request", func(t *testing.T) {
		req, err := k8sClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		expectedURL := "https://api.test.com/kubernetes/test"
		if req.URL.String() != expectedURL {
			t.Errorf("URL esperada: %s, obtida: %s", expectedURL, req.URL.String())
		}

		if req.Header.Get("X-API-Key") != "test-api-key" {
			t.Error("Header X-API-Key ausente ou incorreto")
		}
	})

	t.Run("handle invalid input", func(t *testing.T) {
		_, err := k8sClient.newRequest(context.Background(), "INVALID\nMETHOD", "/test", nil)
		if err == nil {
			t.Error("Esperado erro com método inválido")
		}
	})
}

func TestKubernetesClient_RequestIDPropagation(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"))
	k8sClient := New(core)

	ctx := context.WithValue(context.Background(), client.RequestIDKey, "test-request-123")

	req, err := k8sClient.newRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if req.Header.Get("X-Request-ID") != "test-request-123" {
		t.Error("X-Request-ID não propagado corretamente")
	}
}

func TestKubernetesClient_RetryConfiguration(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithRetryConfig(3, 1*time.Second, 10*time.Second, 2.0),
	)
	k8sClient := New(core)

	if k8sClient.GetConfig().RetryConfig.MaxAttempts != 3 {
		t.Error("Configuração de retry não aplicada corretamente")
	}
}

func TestKubernetesClient_DefaultBasePath(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL("https://api.test.com"),
	)
	k8sClient := New(core)

	req, err := k8sClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedPath := "/kubernetes/test"
	if !strings.Contains(req.URL.Path, expectedPath) {
		t.Errorf("Caminho base padrão incorreto. Esperado: %s, Obtido: %s", expectedPath, req.URL.Path)
	}
}

func TestNewKubernetesClient_WithOptions(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL("https://api.test.com"),
		client.WithHTTPClient(&http.Client{}),
	)

	var calledOpt bool
	testOpt := func(c *KubernetesClient) {
		calledOpt = true
	}

	bsClient := New(core, testOpt)

	if !calledOpt {
		t.Error("expected option to be called")
	}

	if bsClient == nil {
		t.Error("expected client to not be nil")
	}
}
