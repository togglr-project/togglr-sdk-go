package togglr

import (
	"time"
)

// Option represents a functional option for configuring the client
type Option func(*Config)

// WithBaseURL sets the base URL for the API
func WithBaseURL(url string) Option {
	return func(cfg *Config) {
		cfg.BaseURL = url
	}
}

func WithInsecure() Option {
	return func(config *Config) {
		config.Insecure = true
	}
}

// WithTimeout sets the request timeout
func WithTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = d
	}
}

// WithRetries sets the number of retries
func WithRetries(n int) Option {
	return func(cfg *Config) {
		cfg.Retries = n
	}
}

// WithBackoff sets the backoff policy
func WithBackoff(b Backoff) Option {
	return func(cfg *Config) {
		cfg.Backoff = b
	}
}

// WithCache enables caching with the specified size and TTL
func WithCache(size int, ttl time.Duration) Option {
	return func(cfg *Config) {
		cfg.CacheEnabled = true
		cfg.CacheSize = size
		cfg.CacheTTL = ttl
	}
}

// WithLogger sets the logger
func WithLogger(l Logger) Option {
	return func(cfg *Config) {
		cfg.Logger = l
	}
}

// WithMetrics sets the metrics collector
func WithMetrics(m Metrics) Option {
	return func(cfg *Config) {
		cfg.Metrics = m
	}
}

// WithMaxConns sets the maximum number of connections
func WithMaxConns(maxConns int) Option {
	return func(cfg *Config) {
		cfg.MaxConns = maxConns
	}
}

// WithClientCert sets the client certificate file path
func WithClientCert(certPath string) Option {
	return func(cfg *Config) {
		cfg.ClientCert = certPath
	}
}

// WithClientKey sets the client private key file path
func WithClientKey(keyPath string) Option {
	return func(cfg *Config) {
		cfg.ClientKey = keyPath
	}
}

// WithCACert sets the CA certificate file path
func WithCACert(caPath string) Option {
	return func(cfg *Config) {
		cfg.CACert = caPath
	}
}

// WithClientCertAndKey sets both client certificate and key file paths
func WithClientCertAndKey(certPath, keyPath string) Option {
	return func(cfg *Config) {
		cfg.ClientCert = certPath
		cfg.ClientKey = keyPath
	}
}
