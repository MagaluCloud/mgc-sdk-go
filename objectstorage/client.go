package objectstorage

import (
	"net/http"
	"strings"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ObjectStorageClient represents a client for the object storage service.
// It encapsulates functionality to access buckets and objects using MinIO as the backend.
type ObjectStorageClient struct {
	*client.CoreClient
	minioClient minioClientInterface
	endpoint    Endpoint
}

// ClientOption allows customizing the object storage client configuration.
type ClientOption func(*ObjectStorageClient)

// WithEndpoint sets a custom endpoint for the object storage client.
// If not specified, BR-SE1 is used as the default.
func WithEndpoint(endpoint Endpoint) ClientOption {
	return func(c *ObjectStorageClient) {
		c.endpoint = endpoint
	}
}

// WithMinioClient sets a custom MinIO client.
func WithMinioClient(minioClient *minio.Client) ClientOption {
	return func(c *ObjectStorageClient) {
		c.minioClient = minioClient
	}
}

// WithMinioClientInterface sets a custom MinIO client interface (for testing).
func WithMinioClientInterface(minioClient minioClientInterface) ClientOption {
	return func(c *ObjectStorageClient) {
		c.minioClient = minioClient
	}
}

// New creates a new instance of ObjectStorageClient.
// The default endpoint is BR-SE1. Use WithEndpoint option to specify a different region.
// If the core client is nil, returns an error.
func New(core *client.CoreClient, accessKey string, secretKey string, opts ...ClientOption) (*ObjectStorageClient, error) {
	if core == nil {
		return nil, &client.ValidationError{
			Field:   "core",
			Message: "core client cannot be nil",
		}
	}

	if accessKey == "" {
		return nil, &client.ValidationError{
			Field:   "accessKey",
			Message: "access key cannot be empty",
		}
	}

	if secretKey == "" {
		return nil, &client.ValidationError{
			Field:   "secretKey",
			Message: "secret key cannot be empty",
		}
	}

	osClient := &ObjectStorageClient{
		CoreClient: core,
		endpoint:   BrSe1,
	}

	for _, opt := range opts {
		opt(osClient)
	}

	if err := ValidateEndpoint(osClient.endpoint); err != nil {
		return nil, &client.ValidationError{
			Field:   "endpoint",
			Message: err.Error(),
		}
	}

	// Only create a new MinIO client if one wasn't provided via options
	if osClient.minioClient == nil {
		// MinIO requires just the hostname, not the full URL
		minioEndpoint := parseEndpoint(osClient.endpoint)

		minioClient, err := minio.New(minioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: true,
			Transport: &forceDeleteTransport{
				base: http.DefaultTransport,
			},
		})
		if err != nil {
			return nil, err
		}
		osClient.minioClient = minioClient
	}

	osClient.minioClient.SetAppInfo("wrapper", core.GetConfig().UserAgent)

	return osClient, nil
}

// NewWithEndpoint creates a new instance of ObjectStorageClient with a specific endpoint.
// Deprecated: Use New() with WithEndpoint() option instead.
func NewWithEndpoint(core *client.CoreClient, endpoint Endpoint, accessKey string, secretKey string, opts ...ClientOption) (*ObjectStorageClient, error) {
	return New(core, accessKey, secretKey, append(opts, WithEndpoint(endpoint))...)
}

// parseEndpoint extracts the host from a full endpoint URL.
// Example: "https://br-se1.magaluobjects.com" -> "br-se1.magaluobjects.com"
func parseEndpoint(endpoint Endpoint) string {
	endpointStr := endpoint.String()
	if endpointStr == "" {
		return ""
	}

	// Remove "https://" prefix if present
	endpointStr = strings.TrimPrefix(endpointStr, "https://")
	// Remove "http://" prefix if present
	endpointStr = strings.TrimPrefix(endpointStr, "http://")

	return endpointStr
}

// Buckets returns a service to manage buckets.
// This method allows access to functionality such as creating, listing, and managing buckets.
func (c *ObjectStorageClient) Buckets() BucketService {
	return &bucketService{client: c}
}

// Objects returns a service to manage objects.
// This method allows access to functionality such as uploading, downloading, and managing objects.
func (c *ObjectStorageClient) Objects() ObjectService {
	return &objectService{client: c}
}
