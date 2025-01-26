package client

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestWithAPIKey(t *testing.T) {
	config := &Config{}
	apiKey := "test-api-key"
	
	WithAPIKey(apiKey)(config)
	
	if config.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, config.APIKey)
	}
}

func TestWithBaseURL(t *testing.T) {
	config := &Config{}
	url := MgcUrl("https://api.example.com")
	
	WithBaseURL(url)(config)
	
	if config.BaseURL != url {
		t.Errorf("Expected BaseURL to be %s, got %s", url, config.BaseURL)
	}
}

func TestWithUserAgent(t *testing.T) {
	config := &Config{}
	userAgent := "test-user-agent"
	
	WithUserAgent(userAgent)(config)
	
	if config.UserAgent != userAgent {
		t.Errorf("Expected UserAgent to be %s, got %s", userAgent, config.UserAgent)
	}
}

func TestWithLogger(t *testing.T) {
	config := &Config{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	WithLogger(logger)(config)
	
	if config.Logger != logger {
		t.Errorf("Expected Logger to be %v, got %v", logger, config.Logger)
	}
}

func TestWithHTTPClient(t *testing.T) {
	config := &Config{}
	client := &http.Client{}
	
	WithHTTPClient(client)(config)
	
	if config.HTTPClient != client {
		t.Errorf("Expected HTTPClient to be %v, got %v", client, config.HTTPClient)
	}
}

func TestWithTimeout(t *testing.T) {
	config := &Config{}
	timeout := 30 * time.Second
	
	WithTimeout(timeout)(config)
	
	if config.Timeout != timeout {
		t.Errorf("Expected Timeout to be %v, got %v", timeout, config.Timeout)
	}
}

func TestWithRetryConfig(t *testing.T) {
	config := &Config{}
	maxAttempts := 3
	initialInterval := 1 * time.Second
	maxInterval := 10 * time.Second
	backoffFactor := 2.0
	
	WithRetryConfig(maxAttempts, initialInterval, maxInterval, backoffFactor)(config)
	
	expected := RetryConfig{
		MaxAttempts:     maxAttempts,
		InitialInterval: initialInterval,
		MaxInterval:     maxInterval,
		BackoffFactor:   backoffFactor,
	}
	
	if config.RetryConfig != expected {
		t.Errorf("Expected RetryConfig to be %+v, got %+v", expected, config.RetryConfig)
	}
}

func TestMultipleOptions(t *testing.T) {
	config := &Config{}
	apiKey := "test-api-key"
	url := MgcUrl("https://api.example.com")
	timeout := 30 * time.Second
	
	options := []Option{
		WithAPIKey(apiKey),
		WithBaseURL(url),
		WithTimeout(timeout),
	}
	
	for _, opt := range options {
		opt(config)
	}
	
	if config.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, config.APIKey)
	}
	if config.BaseURL != url {
		t.Errorf("Expected BaseURL to be %s, got %s", url, config.BaseURL)
	}
	if config.Timeout != timeout {
		t.Errorf("Expected Timeout to be %v, got %v", timeout, config.Timeout)
	}
}

func TestEmptyValues(t *testing.T) {
    config := &Config{}
    
    WithAPIKey("")(config)
    WithUserAgent("")(config)
    WithBaseURL("")(config)
    
    if config.APIKey != "" {
        t.Error("Expected empty APIKey")
    }
    if config.UserAgent != "" {
        t.Error("Expected empty UserAgent")
    }
    if config.BaseURL != "" {
        t.Error("Expected empty BaseURL")
    }
}

func TestNilValues(t *testing.T) {
    config := &Config{}
    
    WithLogger(nil)(config)
    WithHTTPClient(nil)(config)
    
    if config.Logger != nil {
        t.Error("Expected nil Logger")
    }
    if config.HTTPClient != nil {
        t.Error("Expected nil HTTPClient")
    }
}

func TestZeroDurationTimeout(t *testing.T) {
    config := &Config{}
    WithTimeout(0)(config)
    
    if config.Timeout != 0 {
        t.Errorf("Expected zero Timeout, got %v", config.Timeout)
    }
}

func TestNegativeRetryValues(t *testing.T) {
    config := &Config{}
    WithRetryConfig(-1, -1*time.Second, -10*time.Second, -1.0)(config)
    
    if config.RetryConfig.MaxAttempts > 0 {
        t.Error("Expected non-positive MaxAttempts")
    }
    if config.RetryConfig.InitialInterval > 0 {
        t.Error("Expected non-positive InitialInterval")
    }
    if config.RetryConfig.MaxInterval > 0 {
        t.Error("Expected non-positive MaxInterval")
    }
    if config.RetryConfig.BackoffFactor > 0 {
        t.Error("Expected non-positive BackoffFactor")
    }
}

func TestComplexConfiguration(t *testing.T) {
    config := &Config{}
    
    WithAPIKey("key1")(config)
    WithAPIKey("key2")(config)
    WithTimeout(1 * time.Second)(config)
    WithTimeout(2 * time.Second)(config)
    WithRetryConfig(3, 1*time.Second, 5*time.Second, 2.0)(config)
    WithRetryConfig(5, 2*time.Second, 10*time.Second, 3.0)(config)
    
    expected := Config{
        APIKey:  "key2",
        Timeout: 2 * time.Second,
        RetryConfig: RetryConfig{
            MaxAttempts:     5,
            InitialInterval: 2 * time.Second,
            MaxInterval:     10 * time.Second,
            BackoffFactor:   3.0,
        },
    }
    
    if config.APIKey != expected.APIKey {
        t.Errorf("APIKey mismatch: got %s, want %s", config.APIKey, expected.APIKey)
    }
    if config.Timeout != expected.Timeout {
        t.Errorf("Timeout mismatch: got %v, want %v", config.Timeout, expected.Timeout)
    }
    if config.RetryConfig != expected.RetryConfig {
        t.Errorf("RetryConfig mismatch: got %+v, want %+v", config.RetryConfig, expected.RetryConfig)
    }
}

func TestConfigurationChaining(t *testing.T) {
    config := &Config{}
    options := []Option{
        WithAPIKey("test-key"),
        WithBaseURL("https://api.test.com"),
        WithTimeout(5 * time.Second),
        WithUserAgent("test-agent"),
        WithRetryConfig(3, 1*time.Second, 5*time.Second, 2.0),
        WithLogger(slog.New(slog.NewTextHandler(os.Stdout, nil))),
        WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
    }
    
    for _, opt := range options {
        opt(config)
    }
    
    if config.APIKey == "" || config.BaseURL == "" || config.UserAgent == "" ||
        config.Logger == nil || config.HTTPClient == nil || config.Timeout == 0 ||
        config.RetryConfig.MaxAttempts == 0 {
        t.Error("One or more configuration options were not properly set during chaining")
    }
}

func TestOverrideRetryConfig(t *testing.T) {
    config := &Config{}
    initial := RetryConfig{
        MaxAttempts:     3,
        InitialInterval: 1 * time.Second,
        MaxInterval:     5 * time.Second,
        BackoffFactor:   2.0,
    }
    
    WithRetryConfig(
        initial.MaxAttempts,
        initial.InitialInterval,
        initial.MaxInterval,
        initial.BackoffFactor,
    )(config)
    
    override := RetryConfig{
        MaxAttempts:     5,
        InitialInterval: 2 * time.Second,
        MaxInterval:     10 * time.Second,
        BackoffFactor:   3.0,
    }
    
    WithRetryConfig(
        override.MaxAttempts,
        override.InitialInterval,
        override.MaxInterval,
        override.BackoffFactor,
    )(config)
    
    if config.RetryConfig != override {
        t.Errorf("RetryConfig was not properly overridden. Got %+v, want %+v",
            config.RetryConfig, override)
    }
}
