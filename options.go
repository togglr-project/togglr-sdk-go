package togglr

import (
	"time"
)

type Option func(*Config)

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

func WithTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = d
	}
}

func WithRetries(n int) Option {
	return func(cfg *Config) {
		cfg.Retries = n
	}
}

func WithBackoff(b Backoff) Option {
	return func(cfg *Config) {
		cfg.Backoff = b
	}
}

func WithCache(size int, ttl time.Duration) Option {
	return func(cfg *Config) {
		cfg.CacheEnabled = true
		cfg.CacheSize = size
		cfg.CacheTTL = ttl
	}
}

func WithLogger(l Logger) Option {
	return func(cfg *Config) {
		cfg.Logger = l
	}
}

func WithMetrics(m Metrics) Option {
	return func(cfg *Config) {
		cfg.Metrics = m
	}
}

func WithMaxConns(maxConns int) Option {
	return func(cfg *Config) {
		cfg.MaxConns = maxConns
	}
}

func WithClientCert(certPath string) Option {
	return func(cfg *Config) {
		cfg.ClientCert = certPath
	}
}

func WithClientKey(keyPath string) Option {
	return func(cfg *Config) {
		cfg.ClientKey = keyPath
	}
}

func WithCACert(caPath string) Option {
	return func(cfg *Config) {
		cfg.CACert = caPath
	}
}

func WithClientCertAndKey(certPath, keyPath string) Option {
	return func(cfg *Config) {
		cfg.ClientCert = certPath
		cfg.ClientKey = keyPath
	}
}
