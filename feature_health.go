package togglr

import (
	"context"
	"fmt"
	"time"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

// FeatureHealth represents the health status of a feature
type FeatureHealth struct {
	FeatureKey     string
	EnvironmentKey string
	Enabled        bool
	AutoDisabled   bool
	ErrorRate      *float32
	Threshold      *float32
	LastErrorAt    *time.Time
}

// GetFeatureHealth retrieves the health status of a feature
func (c *Client) GetFeatureHealth(ctx context.Context, featureKey string) (*FeatureHealth, error) {
	start := time.Now()
	c.metrics.IncFeatureHealthRequest()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	// Make API call with retries
	health, err := c.getFeatureHealthWithRetries(ctx, featureKey)

	// Record metrics
	c.metrics.ObserveFeatureHealthLatency(time.Since(start))
	if err != nil {
		c.metrics.IncFeatureHealthError("get_health_failed")
	}

	return health, err
}

// getFeatureHealthWithRetries performs the actual API call with retry logic
func (c *Client) getFeatureHealthWithRetries(
	ctx context.Context,
	featureKey string,
) (*FeatureHealth, error) {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying get health after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Make API call
		params := api.GetFeatureHealthParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.GetFeatureHealth(ctx, params)
		if err == nil {
			// Handle response
			switch r := resp.(type) {
			case *api.FeatureHealth:
				return convertFeatureHealth(r), nil
			case *api.ErrorBadRequest:
				return nil, ErrBadRequest
			case *api.ErrorUnauthorized:
				return nil, ErrUnauthorized
			case *api.ErrorNotFound:
				return nil, ErrFeatureNotFound
			case *api.ErrorInternalServerError:
				return nil, ErrInternalServerError
			default:
				return nil, fmt.Errorf("unexpected response type: %T", resp)
			}
		}

		lastErr = err

		// Check if we should retry
		if !shouldRetry(err) {
			c.logger.Debug("not retrying get health due to error type", "error", err)
			break
		}

		c.logger.Debug("retrying get health due to error", "attempt", attempt, "error", err)
	}

	return nil, lastErr
}

// IsFeatureHealthy checks if a feature is healthy (enabled and not auto-disabled)
func (c *Client) IsFeatureHealthy(ctx context.Context, featureKey string) (bool, error) {
	health, err := c.GetFeatureHealth(ctx, featureKey)
	if err != nil {
		return false, err
	}

	return health.Enabled && !health.AutoDisabled, nil
}

// convertFeatureHealth converts API FeatureHealth to SDK FeatureHealth
func convertFeatureHealth(apiHealth *api.FeatureHealth) *FeatureHealth {
	health := &FeatureHealth{
		FeatureKey:     apiHealth.FeatureKey,
		EnvironmentKey: apiHealth.EnvironmentKey,
		Enabled:        apiHealth.Enabled,
		AutoDisabled:   apiHealth.AutoDisabled,
	}

	if apiHealth.ErrorRate.IsSet() {
		health.ErrorRate = &apiHealth.ErrorRate.Value
	}

	if apiHealth.Threshold.IsSet() {
		health.Threshold = &apiHealth.Threshold.Value
	}

	if apiHealth.LastErrorAt.IsSet() {
		health.LastErrorAt = &apiHealth.LastErrorAt.Value
	}

	return health
}
