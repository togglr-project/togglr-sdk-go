package togglr

import (
	"time"
)

type Config struct {
	BaseURL      string
	APIKey       string
	Timeout      time.Duration
	Retries      int
	Backoff      Backoff
	CacheEnabled bool
	CacheSize    int
	CacheTTL     time.Duration
	Logger       Logger
	Metrics      Metrics
	MaxConns     int
	Insecure     bool
	ClientCert   string
	ClientKey    string
	CACert       string
}

type Backoff struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
	Factor    float64
}

func DefaultBackoff() Backoff {
	return Backoff{
		BaseDelay: 100 * time.Millisecond,
		MaxDelay:  2 * time.Second,
		Factor:    2.0,
	}
}

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
