package objectstorage

import (
	"testing"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestNewObjectStorageClient(t *testing.T) {
	tests := []struct {
		name      string
		core      *client.CoreClient
		accessKey string
		secretKey string
		opts      []ClientOption
		wantErr   bool
		errField  string
		wantEp    Endpoint
	}{
		{
			name:      "invalid - nil core client",
			core:      nil,
			accessKey: "minioadmin",
			secretKey: "minioadmin",
			wantErr:   true,
			errField:  "core",
		},
		{
			name:      "invalid - empty access key",
			core:      createMockCoreClient(),
			accessKey: "",
			secretKey: "minioadmin",
			wantErr:   true,
			errField:  "accessKey",
		},
		{
			name:      "invalid - empty secret key",
			core:      createMockCoreClient(),
			accessKey: "minioadmin",
			secretKey: "",
			wantErr:   true,
			errField:  "secretKey",
		},
		{
			name:      "valid - default BR-SE1 endpoint",
			core:      createMockCoreClient(),
			accessKey: "minioadmin",
			secretKey: "minioadmin",
			wantErr:   false,
			wantEp:    BrSe1,
		},
		{
			name:      "valid - custom BR-NE1 endpoint",
			core:      createMockCoreClient(),
			accessKey: "minioadmin",
			secretKey: "minioadmin",
			opts:      []ClientOption{WithEndpoint(BrNe1)},
			wantErr:   false,
			wantEp:    BrNe1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osClient, err := New(tt.core, tt.accessKey, tt.secretKey, tt.opts...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("New() error = nil, wantErr %v", tt.wantErr)
					return
				}

				validErr, ok := err.(*client.ValidationError)
				if !ok {
					t.Errorf("New() expected ValidationError, got %T", err)
					return
				}

				if validErr.Field != tt.errField {
					t.Errorf("New() error field = %s, want %s", validErr.Field, tt.errField)
				}
				return
			}

			if err != nil {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if osClient == nil {
				t.Error("New() returned nil client")
				return
			}

			if osClient.endpoint != tt.wantEp {
				t.Errorf("New() endpoint = %v, want %v", osClient.endpoint, tt.wantEp)
			}
		})
	}
}

func TestDefaultEndpointIsBRSE1(t *testing.T) {
	t.Parallel()

	core := createMockCoreClient()
	osClient, err := New(core, "minioadmin", "minioadmin")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if osClient.endpoint != BrSe1 {
		t.Errorf("expected default endpoint to be BrSe1, got %v", osClient.endpoint)
	}
}

func TestObjectStorageClientBuckets(t *testing.T) {
	t.Parallel()

	core := createMockCoreClient()
	osClient, err := New(core, "minioadmin", "minioadmin")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	buckets := osClient.Buckets()
	if buckets == nil {
		t.Errorf("expected buckets service, got nil")
	}

	bucketService, ok := buckets.(*bucketService)
	if !ok {
		t.Errorf("expected *bucketService, got %T", buckets)
	}

	if bucketService.client != osClient {
		t.Errorf("expected bucketService to reference the ObjectStorageClient")
	}
}

func TestObjectStorageClientObjects(t *testing.T) {
	t.Parallel()

	core := createMockCoreClient()
	osClient, err := New(core, "minioadmin", "minioadmin")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	objects := osClient.Objects()
	if objects == nil {
		t.Errorf("expected objects service, got nil")
	}

	objectService, ok := objects.(*objectService)
	if !ok {
		t.Errorf("expected *objectService, got %T", objects)
	}

	if objectService.client != osClient {
		t.Errorf("expected objectService to reference the ObjectStorageClient")
	}
}

func TestWithEndpointOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint Endpoint
	}{
		{
			name:     "set endpoint to BR-SE1",
			endpoint: BrSe1,
		},
		{
			name:     "set endpoint to BR-NE1",
			endpoint: BrNe1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := createMockCoreClient()
			osClient, err := New(core, "minioadmin", "minioadmin", WithEndpoint(tt.endpoint))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			if osClient.endpoint != tt.endpoint {
				t.Errorf("expected endpoint %v, got %v", tt.endpoint, osClient.endpoint)
			}
		})
	}
}

func TestWithMinioClientOption(t *testing.T) {
	t.Parallel()

	core := createMockCoreClient()

	client1, err := New(core, "minioadmin", "minioadmin")
	if err != nil {
		t.Fatalf("failed to create first client: %v", err)
	}

	mockMinioClient := client1.minioClient

	client2, err := New(core, "admin", "admin", WithEndpoint(BrNe1), WithMinioClientInterface(mockMinioClient))
	if err != nil {
		t.Fatalf("failed to create second client: %v", err)
	}

	if client2.minioClient != mockMinioClient {
		t.Errorf("expected minioClient to be set via WithMinioClient option")
	}

	if client2.endpoint != BrNe1 {
		t.Errorf("expected endpoint to be BrNe1, got %v", client2.endpoint)
	}
}

func TestNewWithEndpointDeprecated(t *testing.T) {
	t.Parallel()

	core := createMockCoreClient()
	osClient, err := NewWithEndpoint(core, BrNe1, "minioadmin", "minioadmin")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if osClient.endpoint != BrNe1 {
		t.Errorf("expected endpoint BrNe1, got %v", osClient.endpoint)
	}
}

func TestNewWithInvalidEndpoint(t *testing.T) {
	t.Parallel()

	core := createMockCoreClient()
	invalidEndpoint := Endpoint("invalid-endpoint")

	_, err := New(core, "minioadmin", "minioadmin", WithEndpoint(invalidEndpoint))

	if err == nil {
		t.Error("New() expected error for invalid endpoint, got nil")
	}

	validErr, ok := err.(*client.ValidationError)
	if !ok {
		t.Errorf("New() expected ValidationError, got %T", err)
		return
	}

	if validErr.Field != "endpoint" {
		t.Errorf("New() error field = %s, want endpoint", validErr.Field)
	}
}

func createMockCoreClient() *client.CoreClient {
	return client.NewMgcClient("mock-api-token")
}

func TestParseEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint Endpoint
		want     string
	}{
		{
			name:     "https prefix removed",
			endpoint: BrSe1,
			want:     "br-se1.magaluobjects.com",
		},
		{
			name:     "http prefix removed",
			endpoint: Endpoint("http://br-ne1.magaluobjects.com"),
			want:     "br-ne1.magaluobjects.com",
		},
		{
			name:     "empty endpoint",
			endpoint: Endpoint(""),
			want:     "",
		},
		{
			name:     "no prefix",
			endpoint: Endpoint("br-se1.magaluobjects.com"),
			want:     "br-se1.magaluobjects.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseEndpoint(tt.endpoint)
			if got != tt.want {
				t.Errorf("parseEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
