# Containerregistry

Package containerregistry provides a client for interacting with the Magalu
Cloud Container Registry API. This package allows you to manage container
registries, repositories, images, and credentials.
```
const (
DefaultBasePath = "/container-registry"
)


```
```
type AmountRepositoryResponse struct {
Total int `json:"total"`
}
AmountRepositoryResponse represents the total count of repositories

```
```
type ClientOption func(*ContainerRegistryClient)
ClientOption is a function type for configuring ContainerRegistryClient
options

```
```
type ContainerRegistryClient struct {
*client.CoreClient
}
ContainerRegistryClient represents a client for the Container Registry
service

```
```
func New(core *client.CoreClient, opts ...ClientOption) *ContainerRegistryClient
New creates a new ContainerRegistryClient instance with the provided core
client and options

```
```
func (c *ContainerRegistryClient) Credentials() CredentialsService
Credentials returns a service for managing container registry credentials

```
```
func (c *ContainerRegistryClient) Images() ImagesService
Images returns a service for managing images within repositories

```
```
func (c *ContainerRegistryClient) Registries() RegistriesService
Registries returns a service for managing container registries

```
```
func (c *ContainerRegistryClient) Repositories() RepositoriesService
Repositories returns a service for managing repositories within registries

```
```
type CredentialsResponse struct {
Username string `json:"username"`
Password string `json:"password"`
Email    string `json:"email"`
}
CredentialsResponse represents the response containing registry credentials

```
```
type CredentialsService interface {
Get(ctx context.Context) (*CredentialsResponse, error)
ResetPassword(ctx context.Context) (*CredentialsResponse, error)
}
CredentialsService provides methods for managing container registry
credentials

```
```
type ImageResponse struct {
Digest            string             `json:"digest"`
SizeBytes         int                `json:"size_bytes"`
PushedAt          string             `json:"pushed_at"`
PulledAt          string             `json:"pulled_at"`
ManifestMediaType string             `json:"manifest_media_type"`
MediaType         string             `json:"media_type"`
Tags              []string           `json:"tags"`
TagsDetails       []ImageTagResponse `json:"tags_details"`
ExtraAttr         string             `json:"extra_attr"`
}
ImageResponse represents a container image

```
```
type ImageTagResponse struct {
Name     string `json:"name"`
PushedAt string `json:"pushed_at"`
PulledAt string `json:"pulled_at"`
Signed   bool   `json:"signed"`
}
ImageTagResponse represents detailed information about an image tag

```
```
type ImagesResponse struct {
Results []ImageResponse `json:"results"`
}
ImagesResponse represents the response when listing images

```
```
type ImagesService interface {
List(ctx context.Context, registryID, repositoryName string, opts ListOptions) (*ImagesResponse, error)
Delete(ctx context.Context, registryID, repositoryName, digestOrTag string) error
Get(ctx context.Context, registryID, repositoryName, digestOrTag string) (*ImageResponse, error)
}
ImagesService provides methods for managing images within repositories

```
```
type ListOptions struct {
Limit  *int
Offset *int
Sort   *string
Expand []string
}
ListOptions provides options for listing registries

```
```
type ListRegistriesResponse struct {
Registries []RegistryResponse `json:"results"`
}
ListRegistriesResponse represents the response when listing registries

```
```
type RegistriesService interface {
Create(ctx context.Context, request *RegistryRequest) (*RegistryResponse, error)
List(ctx context.Context, opts ListOptions) (*ListRegistriesResponse, error)
Get(ctx context.Context, registryID string) (*RegistryResponse, error)
Delete(ctx context.Context, registryID string) error
}
RegistriesService provides methods for managing container registries

```
```
type RegistryRequest struct {
Name string `json:"name"`
}
RegistryRequest represents the request payload for creating a registry

```
```
type RegistryResponse struct {
ID        string `json:"id"`
Name      string `json:"name"`
Storage   int    `json:"storage_usage_bytes"`
CreatedAt string `json:"created_at"`
UpdatedAt string `json:"updated_at"`
}
RegistryResponse represents a container registry

```
```
type RepositoriesResponse struct {
Goal    AmountRepositoryResponse `json:"goal"`
Results []RepositoryResponse     `json:"results"`
}
RepositoriesResponse represents the response when listing repositories

```
```
type RepositoriesService interface {
List(ctx context.Context, registryID string, opts ListOptions) (*RepositoriesResponse, error)
Get(ctx context.Context, registryID, repositoryName string) (*RepositoryResponse, error)
Delete(ctx context.Context, registryID, repositoryName string) error
}
RepositoriesService provides methods for managing repositories within
container registries

```
```
type RepositoryResponse struct {
RegistryName string `json:"registry_name"`
Name         string `json:"name"`
ImageCount   int    `json:"image_count"`
CreatedAt    string `json:"created_at"`
UpdatedAt    string `json:"updated_at"`
}
RepositoryResponse represents a repository within a container registry


```

