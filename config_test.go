package togglr

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig("test-api-key")

	if cfg.APIKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got %s", cfg.APIKey)
	}

	if cfg.BaseURL != "http://localhost:8090" {
		t.Errorf("Expected base URL 'http://localhost:8090', got %s", cfg.BaseURL)
	}

	if cfg.Timeout != 800*time.Millisecond {
		t.Errorf("Expected timeout 800ms, got %v", cfg.Timeout)
	}

	if cfg.Retries != 2 {
		t.Errorf("Expected retries 2, got %d", cfg.Retries)
	}

	if cfg.CacheEnabled != false {
		t.Errorf("Expected cache disabled by default, got %v", cfg.CacheEnabled)
	}
}

func TestOptions(t *testing.T) {
	cfg := DefaultConfig("test-api-key")

	// Test WithBaseURL
	WithBaseURL("https://api.example.com")(cfg)
	if cfg.BaseURL != "https://api.example.com" {
		t.Errorf("Expected base URL 'https://api.example.com', got %s", cfg.BaseURL)
	}

	// Test WithTimeout
	WithTimeout(2 * time.Second)(cfg)
	if cfg.Timeout != 2*time.Second {
		t.Errorf("Expected timeout 2s, got %v", cfg.Timeout)
	}

	// Test WithRetries
	WithRetries(5)(cfg)
	if cfg.Retries != 5 {
		t.Errorf("Expected retries 5, got %d", cfg.Retries)
	}

	// Test WithCache
	WithCache(1000, 30*time.Second)(cfg)
	if !cfg.CacheEnabled {
		t.Errorf("Expected cache enabled")
	}
	if cfg.CacheSize != 1000 {
		t.Errorf("Expected cache size 1000, got %d", cfg.CacheSize)
	}
	if cfg.CacheTTL != 30*time.Second {
		t.Errorf("Expected cache TTL 30s, got %v", cfg.CacheTTL)
	}
}
