# Sshkeys

Package sshkeys provides client implementation for managing SSH keys in the
Magalu Cloud platform. SSH keys are managed as a global service, meaning they
are not bound to any specific region. By default, the service uses the global
endpoint, but this can be overridden if needed.
```
const (
// DefaultBasePath is the default API base path for SSH key operations
DefaultBasePath = "/profile"
)


```
```
type ClientOption func(*SSHKeyClient)
ClientOption allows customizing the SSH key client configuration.

```
```
func WithGlobalBasePath(basePath client.MgcUrl) ClientOption
WithGlobalBasePath allows overriding the default global endpoint for SSH
keys service. This is rarely needed as SSH keys are managed globally,
but provided for flexibility.

Example:

client := sshkeys.New(core, sshkeys.WithGlobalBasePath("custom-endpoint"))

```
```
type CreateSSHKeyRequest struct {
Name string `json:"name"`
Key  string `json:"key"`
}
CreateSSHKeyRequest represents the parameters for creating a new SSH key

```
```
type KeyService interface {
List(ctx context.Context, opts ListOptions) ([]SSHKey, error)
Create(ctx context.Context, req CreateSSHKeyRequest) (*SSHKey, error)
Get(ctx context.Context, keyID string) (*SSHKey, error)
Delete(ctx context.Context, keyID string) (*SSHKey, error)
}
KeyService provides methods for managing SSH keys. All operations in this
service are performed against the global endpoint, as SSH keys are not
region-specific resources.

```
```
type ListOptions struct {
Limit  *int
Offset *int
Sort   *string
}
ListOptions defines parameters for filtering and paginating SSH key lists

```
```
type ListSSHKeysResponse struct {
Results []SSHKey `json:"results"`
}
ListSSHKeysResponse represents a list of SSH keys response

```
```
type SSHKey struct {
ID      string `json:"id"`
Name    string `json:"name"`
Key     string `json:"key"`
KeyType string `json:"key_type"`
}
SSHKey represents an SSH key resource

```
```
type SSHKeyClient struct {
*client.CoreClient
}
SSHKeyClient represents a client for interacting with the SSH keys service

```
```
func New(core *client.CoreClient, opts ...ClientOption) *SSHKeyClient
New creates a new SSH key client using the provided core client. The SSH
keys service operates globally and is not region-specific. By default,
it uses the global endpoint (api.magalu.cloud).

To customize the endpoint, use WithGlobalBasePath option.

```
```
func (c *SSHKeyClient) Keys() KeyService
Keys returns a service for managing SSH key resources


```

