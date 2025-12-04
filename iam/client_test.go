package iam

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestNewIAMClient(t *testing.T) {
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
			iamClient := New(tt.core)
			if (iamClient == nil) != tt.wantNil {
				t.Errorf("New() returned nil = %v, want nil = %v", iamClient == nil, tt.wantNil)
			}
		})
	}
}

func TestIAMClient_Services(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithHTTPClient(&http.Client{}),
		client.WithBaseURL("https://api.test.com"))

	iamClient := New(core)

	t.Run("Members service", func(t *testing.T) {
		service := iamClient.Members()
		if service == nil {
			t.Error("Members() returned nil")
		}
	})

	t.Run("Roles service", func(t *testing.T) {
		service := iamClient.Roles()
		if service == nil {
			t.Error("Roles() returned nil")
		}
	})

	t.Run("Permissions service", func(t *testing.T) {
		service := iamClient.Permissions()
		if service == nil {
			t.Error("Permissions() returned nil")
		}
	})

	t.Run("AccessControl service", func(t *testing.T) {
		service := iamClient.AccessControl()
		if service == nil {
			t.Error("AccessControl() returned nil")
		}
	})

	t.Run("ServiceAccounts service", func(t *testing.T) {
		service := iamClient.ServiceAccounts()
		if service == nil {
			t.Error("ServiceAccounts() returned nil")
		}
	})

	t.Run("Scopes service", func(t *testing.T) {
		service := iamClient.Scopes()
		if service == nil {
			t.Error("Scopes() returned nil")
		}
	})
}

func TestIAMClient_NewRequest(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithHTTPClient(&http.Client{}),
	)
	iamClient := New(core, WithGlobalBasePath("https://api.test.com"))

	t.Run("create valid request", func(t *testing.T) {
		req, err := iamClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		expectedURL := "https://api.test.com/iam/api/v1/test"
		if req.URL.String() != expectedURL {
			t.Errorf("URL esperada: %s, obtida: %s", expectedURL, req.URL.String())
		}

		if req.Header.Get("X-API-Key") != "test-api-key" {
			t.Error("Header X-API-Key ausente ou incorreto")
		}
	})

	t.Run("handle invalid input", func(t *testing.T) {
		_, err := iamClient.newRequest(context.Background(), "INVALID\nMETHOD", "/test", nil)
		if err == nil {
			t.Error("Esperado erro com método inválido")
		}
	})
}

func TestIAMClient_RequestIDPropagation(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"))
	iamClient := New(core)

	ctx := context.WithValue(context.Background(), client.RequestIDKey, "test-request-123")

	req, err := iamClient.newRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if req.Header.Get("X-Request-ID") != "test-request-123" {
		t.Error("X-Request-ID não propagado corretamente")
	}
}

func TestIAMClient_RetryConfiguration(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithRetryConfig(3, 1*time.Second, 10*time.Second, 2.0),
	)
	iamClient := New(core)

	if iamClient.GetConfig().RetryConfig.MaxAttempts != 3 {
		t.Error("Configuração de retry não aplicada corretamente")
	}
}

func TestIAMClient_DefaultBasePath(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL("https://api.test.com"),
	)
	iamClient := New(core)

	req, err := iamClient.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedPath := "/api/v1/test"
	if !strings.Contains(req.URL.Path, expectedPath) {
		t.Errorf("Caminho base padrão incorreto. Esperado: %s, Obtido: %s", expectedPath, req.URL.Path)
	}
}

func TestNewIAMClient_WithOptions(t *testing.T) {
	core := client.NewMgcClient(client.WithAPIKey("test-api-key"),
		client.WithBaseURL("https://api.test.com"),
		client.WithHTTPClient(&http.Client{}),
	)

	var calledOpt bool
	testOpt := func(c *IAMClient) {
		calledOpt = true
	}

	iamClient := New(core, testOpt)

	if !calledOpt {
		t.Error("expected option to be called")
	}

	if iamClient == nil {
		t.Error("expected client to not be nil")
	}
}
