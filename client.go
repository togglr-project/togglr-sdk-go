package togglr

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

// Client represents the Togglr SDK client
type Client struct {
	cfg        *Config
	httpClient *http.Client
	apiClient  *api.Client
	cache      *LRUCache
	logger     Logger
	metrics    Metrics
}

// NewClient creates a new Togglr client with the given configuration
func NewClient(cfg *Config, opts ...Option) (*Client, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set defaults for optional fields
	if cfg.Logger == nil {
		cfg.Logger = &NoOpLogger{}
	}
	if cfg.Metrics == nil {
		cfg.Metrics = &NoOpMetrics{}
	}

	// Create an HTTP client with custom transport
	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxConns,
		MaxIdleConnsPerHost: cfg.MaxConns,
		IdleConnTimeout:     90 * time.Second,
	}

	// Configure TLS
	tlsConfig := &tls.Config{}

	if cfg.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	// Load client certificate and key if provided
	if cfg.ClientCert != "" && cfg.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if provided
	if cfg.CACert != "" {
		caCert, err := os.ReadFile(cfg.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	transport.TLSClientConfig = tlsConfig

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	// Create API client using ogen (как в вашем примере)
	apiClient, err := api.NewClient(
		cfg.BaseURL,
		&apiKeySecuritySource{apiKey: cfg.APIKey},
		api.WithClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Create cache if enabled
	var cache *LRUCache
	if cfg.CacheEnabled {
		cache = NewLRUCache(cfg.CacheSize, cfg.CacheTTL)
	}

	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
		apiClient:  apiClient,
		cache:      cache,
		logger:     cfg.Logger,
		metrics:    cfg.Metrics,
	}, nil
}

// NewClientWithDefaults creates a new client with default configuration
func NewClientWithDefaults(apiKey string, opts ...Option) (*Client, error) {
	cfg := DefaultConfig(apiKey)

	return NewClient(cfg, opts...)
}

// Close closes the client and cleans up resources
func (c *Client) Close() error {
	if c.cache != nil {
		c.cache.Clear()
	}

	return nil
}

// HealthCheck performs a health check on the API
func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.apiClient.SdkV1HealthGet(ctx)

	return err
}

// apiKeySecuritySource implements the SecuritySource interface (как в вашем примере)
type apiKeySecuritySource struct {
	apiKey string
}

func (s *apiKeySecuritySource) ApiKeyAuth(context.Context, api.OperationName, *api.Client) (api.ApiKeyAuth, error) {
	return api.ApiKeyAuth{APIKey: s.apiKey}, nil
}
