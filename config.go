package togglr

import (
	"time"
)

// Config represents the SDK configuration
type Config struct {
	BaseURL      string        // default: "http://localhost:8090"
	APIKey       string        // required
	Timeout      time.Duration // default: 800*time.Millisecond
	Retries      int           // default: 2
	Backoff      Backoff       // backoff policy struct, default exponential
	CacheEnabled bool          // default: false
	CacheSize    int           // default: 100
	CacheTTL     time.Duration // default: 5*time.Second
	Logger       Logger        // optional (interface)
	Metrics      Metrics       // optional (interface)
	MaxConns     int           // optional transport tuning
	Insecure     bool
	// TLS configuration
	ClientCert string // path to client certificate file
	ClientKey  string // path to client private key file
	CACert     string // path to CA certificate file
}

// Backoff represents the backoff policy
type Backoff struct {
	BaseDelay time.Duration // default: 100ms
	MaxDelay  time.Duration // default: 2s
	Factor    float64       // default: 2.0
}

// DefaultBackoff returns the default backoff configuration
func DefaultBackoff() Backoff {
	return Backoff{
		BaseDelay: 100 * time.Millisecond,
		MaxDelay:  2 * time.Second,
		Factor:    2.0,
	}
}

// DefaultConfig creates a default configuration with the provided API key
func DefaultConfig(apiKey string) *Config {
	return &Config{
		BaseURL:      "http://localhost:8090",
		APIKey:       apiKey,
		Timeout:      800 * time.Millisecond,
		Retries:      2,
		Backoff:      DefaultBackoff(),
		CacheEnabled: false,
		CacheSize:    100,
		CacheTTL:     5 * time.Second,
		MaxConns:     100,
	}
}
