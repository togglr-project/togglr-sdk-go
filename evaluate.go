package togglr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-faster/jx"

	"github.com/rom8726/togglr-sdk-go/internal/fingerprint"
	api "github.com/rom8726/togglr-sdk-go/internal/generated/client"
)

// EvaluateWithContext evaluates a feature with the given context
func (c *Client) EvaluateWithContext(
	ctx context.Context,
	featureKey string,
	req RequestContext,
) (value string, enabled bool, found bool, err error) {
	start := time.Now()
	c.metrics.IncEvaluateRequest()

	// Create a cache key if caching is enabled
	var cacheKey string
	if c.cache != nil {
		cacheKey = fmt.Sprintf("%s:%s", featureKey, fingerprint.Fingerprint(req))

		// Check cache first
		if entry, hit := c.cache.Get(cacheKey); hit {
			c.metrics.IncCacheHit()
			c.logger.Debug("cache hit", "feature_key", featureKey, "cache_key", cacheKey)

			return entry.Value, entry.Enabled, entry.Found, nil
		}
		c.metrics.IncCacheMiss()
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	// Convert RequestContext to EvaluateRequest for API call
	reqMap := make(map[string]any)
	for k, v := range req {
		reqMap[k] = v
	}

	// Make API call with retries
	value, enabled, found, err = c.evaluateWithRetries(ctx, featureKey, reqMap)

	// Record metrics
	c.metrics.ObserveEvaluateLatency(time.Since(start))
	if err != nil {
		c.metrics.IncEvaluateError(getErrorCode(err))
	}

	// Cache result if successful and caching is enabled
	if err == nil && c.cache != nil {
		c.cache.Set(cacheKey, value, enabled, found)
	}

	return value, enabled, found, err
}

// Evaluate evaluates a feature using the configured API key
func (c *Client) Evaluate(featureKey string, req RequestContext) (value string, enabled bool, found bool, err error) {
	return c.EvaluateWithContext(context.Background(), featureKey, req)
}

// IsEnabled checks if a feature is enabled
func (c *Client) IsEnabled(featureKey string, req RequestContext) (bool, error) {
	_, enabled, found, err := c.Evaluate(featureKey, req)
	if err != nil {
		return false, err
	}

	if !found {
		return false, ErrFeatureNotFound
	}

	return enabled, nil
}

// IsEnabledOrDefault evaluates a feature and returns a default value on error
func (c *Client) IsEnabledOrDefault(featureKey string, req RequestContext, def bool) bool {
	enabled, err := c.IsEnabled(featureKey, req)
	if err != nil {
		c.logger.Warn("evaluation failed, using default",
			"feature_key", featureKey, "error", err, "default", def)

		return def
	}

	return enabled
}

// evaluateWithRetries performs the actual API call with retry logic
func (c *Client) evaluateWithRetries(
	ctx context.Context,
	featureKey string,
	req map[string]any,
) (string, bool, bool, error) {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return "", false, false, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Convert map to api.EvaluateRequest
		evalReq := make(api.EvaluateRequest, len(req))
		for k, v := range req {
			// Convert value to jx.Raw
			if raw, err := json.Marshal(v); err == nil {
				evalReq[k] = jx.Raw(raw)
			}
		}

		// Make API call using ogen client
		params := api.SdkV1FeaturesFeatureKeyEvaluatePostParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.SdkV1FeaturesFeatureKeyEvaluatePost(ctx, evalReq, params)
		if err == nil {
			// Check response type
			switch r := resp.(type) {
			case *api.EvaluateResponse:
				return r.Value, r.Enabled, true, nil
			case *api.ErrorNotFound:
				return "", false, false, nil // Feature not found, not an error
			case *api.ErrorUnauthorized:
				return "", false, false, ErrUnauthorized
			case *api.ErrorBadRequest:
				return "", false, false, ErrBadRequest
			case *api.ErrorInternalServerError:
				return "", false, false, ErrInternalServerError
			default:
				return "", false, false, fmt.Errorf("unexpected response type: %T", resp)
			}
		}

		lastErr = err

		// Check if we should retry
		if !shouldRetry(err) {
			c.logger.Debug("not retrying due to error type", "error", err)

			break
		}

		c.logger.Debug("retrying due to error", "attempt", attempt, "error", err)
	}

	// Handle specific error types
	var apiErr *APIError
	if errors.As(lastErr, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return "", false, false, nil // Feature not found, not an error
		case http.StatusUnauthorized:
			return "", false, false, ErrUnauthorized
		case http.StatusForbidden:
			return "", false, false, ErrForbidden
		case http.StatusTooManyRequests:
			return "", false, false, ErrTooManyRequests
		}
	}

	return "", false, false, lastErr
}

// calculateBackoffDelay calculates the delay for the given attempt
func (c *Client) calculateBackoffDelay(attempt int) time.Duration {
	delay := c.cfg.Backoff.BaseDelay
	for i := 1; i < attempt; i++ {
		delay = time.Duration(float64(delay) * c.cfg.Backoff.Factor)
		if delay > c.cfg.Backoff.MaxDelay {
			delay = c.cfg.Backoff.MaxDelay

			break
		}
	}

	return delay
}

// shouldRetry determines if an error should trigger a retry
func shouldRetry(err error) bool {
	// For now, retry on any error except context cancellation
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Retry on network errors and timeouts
	return err != nil
}

// getErrorCode extracts error code for metrics
func getErrorCode(err error) string {
	// For now, return a generic error code
	return "unknown"
}
