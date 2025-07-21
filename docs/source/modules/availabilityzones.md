# Availabilityzones

Package availabilityzones provides functionality to interact with the
MagaluCloud availability zones service. This package allows listing availability
zones across different regions.
```
const (
DefaultBasePath = "/profile"
)
DefaultBasePath defines the default base path for availability zones APIs.



```
```
type AvailabilityZone struct {
ID        string    `json:"az_id"`
BlockType BlockType `json:"block_type"`
}
AvailabilityZone represents a single availability zone within a region. Each
availability zone has a unique identifier and can have different blocking
states.

```
```
type BlockType string
BlockType represents the possible blocking states of an availability zone.
An availability zone can be in different states that affect its usability.

```
```
const (
BlockTypeNone     BlockType = "none"
BlockTypeTotal    BlockType = "total"
BlockTypeReadOnly BlockType = "read-only"
)
```
```
type Client struct {
*client.CoreClient
}
Client handles operations on availability zones in the Magalu Cloud
platform. Availability zones are managed as a global service, meaning they
are not bound to any specific region. By default, the service uses the
global endpoint.

```
```
func New(core *client.CoreClient, opts ...ClientOption) *Client
New creates a new availability zones client using the provided core client.
The availability zones service operates globally and is not region-specific.
By default, it uses the global endpoint (api.magalu.cloud).

To customize the endpoint, use WithGlobalBasePath option.

```
```
func (c *Client) AvailabilityZones() Service
AvailabilityZones returns a service to manage availability zones. This
method allows access to functionality such as listing availability zones.

```
```
type ClientOption func(*Client)
ClientOption allows customizing the availability zones client configuration

```
```
func WithGlobalBasePath(basePath client.MgcUrl) ClientOption
WithGlobalBasePath allows overriding the default global endpoint for
availability zones service. This is rarely needed as availability zones are
managed globally, but provided for flexibility.

Example:

client := availabilityzones.New(core, availabilityzones.WithGlobalBasePath("custom-endpoint"))

```
```
type ListOptions struct {
ShowBlocked bool
}
ListOptions contains the options for listing availability zones

```
```
type ListResponse struct {
Results []Region `json:"results"`
}
ListResponse represents the response from listing availability zones.
This structure encapsulates the API response format.

```
```
type Region struct {
ID                string             `json:"region_id"`
AvailabilityZones []AvailabilityZone `json:"availability_zones"`
}
Region represents a region and its associated availability zones.
A region contains multiple availability zones that can be used for resource
deployment.

```
```
type Service interface {
List(ctx context.Context, opts ListOptions) ([]Region, error)
}
Service defines the interface for availability zone operations


```

