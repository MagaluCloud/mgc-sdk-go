# Client

Package client provides the core client functionality for the MagaluCloud SDK.
This package contains the main client implementation, configuration options,
and error handling.
```
const (
RequestIDKey           XRequestID = "x-request-id"
DefaultUserAgent                  = "mgc-sdk-go"
DefaultMaxAttempts                = 3
DefaultInitialInterval            = 1 * time.Second
DefaultMaxInterval                = 30 * time.Second
DefaultBackoffFactor              = 2.0
DefaultTimeout                    = 15 * time.Minute
)
Default configuration constants for the client.



```
```
type Config struct {
APIKey        string
BaseURL       MgcUrl
UserAgent     string
Logger        *slog.Logger
HTTPClient    *http.Client
Timeout       time.Duration
RetryConfig   RetryConfig
ContentType   string
CustomHeaders map[string]string
}
Config contains all configuration options for the client.

```
```
type CoreClient struct {
// Has unexported fields.
}
CoreClient represents the main client for interacting with MagaluCloud APIs.
It encapsulates the configuration and provides methods for making HTTP
requests.

```
```
func NewMgcClient(apiKey string, opts ...Option) *CoreClient
NewMgcClient creates a new instance of CoreClient with the specified API
key and options. The client is configured with sensible defaults and can be
customized using the provided options.

```
```
func (c *CoreClient) GetConfig() *Config
GetConfig returns a pointer to the client's configuration. This method
allows access to the current configuration for inspection or modification.

```
```
type HTTPError struct {
StatusCode int
Status     string
Body       []byte
Response   *http.Response
}
HTTPError represents an error that occurred during an HTTP request. This
error type includes the HTTP status code, status message, and response body.

```
```
func NewHTTPError(resp *http.Response) *HTTPError
NewHTTPError creates a new HTTPError from an HTTP response. This function
reads the response body and creates an error with all available information.

```
```
func (e *HTTPError) Error() string
Error returns a string representation of the HTTP error. This method
implements the error interface.

```
```
type MgcUrl string
MgcUrl represents a MagaluCloud API URL. This type is used to ensure type
safety when working with API endpoints.

```
```
const (
// BrNe1 is the URL for the Brazil Northeast 1 region
BrNe1 MgcUrl = "https://api.magalu.cloud/br-ne1"
// BrSe1 is the URL for the Brazil Southeast 1 region
BrSe1 MgcUrl = "https://api.magalu.cloud/br-se1"
// BrMgl1 is the URL for the Brazil Magalu region
BrMgl1 MgcUrl = "https://api.magalu.cloud/br-se-1"
// Global is the default URL for products that don't have a specific region
Global MgcUrl = "https://api.magalu.cloud"
)
Predefined MagaluCloud API endpoints for different regions.

```
```
func (m MgcUrl) String() string
String returns the string representation of the MgcUrl. This method
implements the Stringer interface.

```
```
type Option func(*Config)
Option is a function type that modifies the client configuration. Options
are used to customize the client behavior during initialization.

```
```
func WithAPIKey(key string) Option
WithAPIKey sets the API key for authentication. This option is required for
all API operations.

```
```
func WithBaseURL(url MgcUrl) Option
WithBaseURL sets the base URL for API requests. This option allows
specifying a custom endpoint for the API.

```
```
func WithCustomHeader(key, value string) Option
WithCustomHeader adds a custom HTTP header to all requests. This option
allows adding additional headers for specific requirements.

```
```
func WithHTTPClient(client *http.Client) Option
WithHTTPClient sets the HTTP client for making requests. This option allows
using a custom HTTP client with specific settings.

```
```
func WithLogger(logger *slog.Logger) Option
WithLogger sets the logger instance for client operations. This option
allows customizing logging behavior.

```
```
func WithRetryConfig(maxAttempts int, initialInterval, maxInterval time.Duration, backoffFactor float64) Option
WithRetryConfig sets the retry configuration for failed requests. This
option allows customizing retry behavior with exponential backoff.

```
```
func WithTimeout(timeout time.Duration) Option
WithTimeout sets the timeout for HTTP requests. This option controls how
long to wait for responses.

```
```
func WithUserAgent(ua string) Option
WithUserAgent sets the user agent string for HTTP requests. This option
allows customizing the user agent header.

```
```
type RetryConfig struct {
MaxAttempts     int
InitialInterval time.Duration
MaxInterval     time.Duration
BackoffFactor   float64
}
RetryConfig contains configuration for retry behavior.

```
```
type RetryError struct {
LastError error
Retries   int
}
RetryError represents an error that occurred after exhausting all retry
attempts. This error type includes the last error encountered and the number
of retries attempted.

```
```
func (e *RetryError) Error() string
Error returns a string representation of the retry error. This method
implements the error interface.

```
```
type ValidationError struct {
Field   string
Message string
}
ValidationError represents an error that occurred during input validation.
This error type includes the field that failed validation and a descriptive
message.

```
```
func (e *ValidationError) Error() string
Error returns a string representation of the validation error. This method
implements the error interface.

```
```
type XRequestID string
XRequestID represents a request ID type for tracking requests.


```

