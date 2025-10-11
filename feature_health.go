package togglr

import (
	"context"
	"fmt"
	"time"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

type FeatureHealth struct {
	FeatureKey     string
	EnvironmentKey string
	Enabled        bool
	AutoDisabled   bool
	ErrorRate      *float32
	Threshold      *float32
	LastErrorAt    *time.Time
}

func (c *Client) GetFeatureHealth(ctx context.Context, featureKey string) (*FeatureHealth, error) {
	start := time.Now()
	c.metrics.IncFeatureHealthRequest()

	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	health, err := c.getFeatureHealthWithRetries(ctx, featureKey)

	c.metrics.ObserveFeatureHealthLatency(time.Since(start))
	if err != nil {
		c.metrics.IncFeatureHealthError("get_health_failed")
	}

	return health, err
}

func (c *Client) getFeatureHealthWithRetries(
	ctx context.Context,
	featureKey string,
) (*FeatureHealth, error) {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying get health after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		params := api.GetFeatureHealthParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.GetFeatureHealth(ctx, params)
		if err == nil {
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

		if !shouldRetry(err) {
			c.logger.Debug("not retrying get health due to error type", "error", err)
			break
		}

		c.logger.Debug("retrying get health due to error", "attempt", attempt, "error", err)
	}

	return nil, lastErr
}

func (c *Client) IsFeatureHealthy(ctx context.Context, featureKey string) (bool, error) {
	health, err := c.GetFeatureHealth(ctx, featureKey)
	if err != nil {
		return false, err
	}

	return health.Enabled && !health.AutoDisabled, nil
}

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
