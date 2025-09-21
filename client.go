package togglr

import (
	"context"
	"fmt"
	"net/http"
	"time"

	api "github.com/rom8726/togglr-sdk-go/internal/generated/client"
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
