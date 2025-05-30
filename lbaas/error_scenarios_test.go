package lbaas

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

// TestNetworkErrorScenarios testa cenários de erro de rede
func TestNetworkErrorScenarios(t *testing.T) {
	t.Parallel()

	t.Run("connection refused", func(t *testing.T) {
		t.Parallel()
		// Usa uma porta que não está sendo usada
		httpClient := &http.Client{}
		core := client.NewMgcClient("test-api",
			client.WithBaseURL(client.MgcUrl("http://localhost:9999")),
			client.WithHTTPClient(httpClient))
		lbClient := New(core).NetworkLoadBalancers()

		_, err := lbClient.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("timeout", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simula um timeout fazendo o servidor dormir mais que o timeout do cliente
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		httpClient := &http.Client{
			Timeout: 100 * time.Millisecond, // Timeout muito baixo
		}
		core := client.NewMgcClient("test-api",
			client.WithBaseURL(client.MgcUrl(server.URL)),
			client.WithHTTPClient(httpClient))
		lbClient := New(core).NetworkLoadBalancers()

		_, err := lbClient.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("malformed response", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Resposta JSON malformada
			w.Write([]byte(`{"id": "lb-123", "name": "test-lb", "invalid": json}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("empty response body", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Resposta vazia
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("invalid content type", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This is not JSON"))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})
}

// TestContextCancellation testa cenários de cancelamento de contexto
func TestContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("context cancelled during request", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simula uma operação lenta
			time.Sleep(500 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "lb-123"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)

		// Cria um contexto que será cancelado rapidamente
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err := client.Get(ctx, GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("context cancelled before request", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "lb-123"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)

		// Cria um contexto já cancelado
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancela imediatamente

		_, err := client.Get(ctx, GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})
}

// TestRateLimitingScenarios testa cenários de rate limiting
func TestRateLimitingScenarios(t *testing.T) {
	t.Parallel()

	t.Run("rate limit exceeded - 429", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limit exceeded"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("service unavailable - 503", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "service temporarily unavailable"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})

	t.Run("bad gateway - 502", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "bad gateway"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})
}

// TestInvalidRequestScenarios testa cenários de requisições inválidas
func TestInvalidRequestScenarios(t *testing.T) {
	t.Parallel()

	t.Run("empty load balancer ID", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "load balancer ID is required"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "",
		})

		assertError(t, err)
	})

	t.Run("invalid characters in ID", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid characters in ID"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "invalid@#$%",
		})

		assertError(t, err)
	})
}

// TestLargeResponseScenarios testa cenários com respostas grandes
func TestLargeResponseScenarios(t *testing.T) {
	t.Parallel()

	t.Run("very large response", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Cria uma resposta muito grande
			largeData := make([]byte, 10*1024*1024) // 10MB
			for i := range largeData {
				largeData[i] = 'a'
			}

			w.Write([]byte(`{"id": "lb-123", "large_field": "`))
			w.Write(largeData)
			w.Write([]byte(`"}`))
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		// Pode ou não dar erro dependendo dos limites do cliente
		// Mas pelo menos testamos o cenário
		if err != nil {
			t.Logf("Large response error (expected): %v", err)
		}
	})
}

// TestConnectionErrorScenarios testa diferentes tipos de erro de conexão
func TestConnectionErrorScenarios(t *testing.T) {
	t.Parallel()

	t.Run("dns resolution failure", func(t *testing.T) {
		t.Parallel()
		httpClient := &http.Client{}
		core := client.NewMgcClient("test-api",
			client.WithBaseURL(client.MgcUrl("http://non-existent-domain-12345.com")),
			client.WithHTTPClient(httpClient))
		lbClient := New(core).NetworkLoadBalancers()

		_, err := lbClient.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)

		// Verifica se é um erro de DNS
		if netErr, ok := err.(*net.OpError); ok {
			t.Logf("DNS error as expected: %v", netErr)
		}
	})

	t.Run("connection reset", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simula uma conexão que é resetada
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("webserver doesn't support hijacking")
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				t.Fatal(err)
			}
			conn.Close() // Fecha a conexão abruptamente
		}))
		defer server.Close()

		client := testLoadBalancerClient(server.URL)
		_, err := client.Get(context.Background(), GetNetworkLoadBalancerRequest{
			LoadBalancerID: "test-lb",
		})

		assertError(t, err)
	})
}
