package togglr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-faster/jx"

	"github.com/togglr-project/togglr-sdk-go/internal/fingerprint"
	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

func (c *Client) EvaluateWithContext(
	ctx context.Context,
	featureKey string,
	req RequestContext,
) EvalResult {
	start := time.Now()
	c.metrics.IncEvaluateRequest()

	var cacheKey string
	if c.cache != nil {
		cacheKey = fmt.Sprintf("%s:%s", featureKey, fingerprint.Fingerprint(req))

		if entry, hit := c.cache.Get(cacheKey); hit {
			c.metrics.IncCacheHit()
			c.logger.Debug("cache hit", "feature_key", featureKey, "cache_key", cacheKey)

			return EvalResult{
				featureKey: featureKey,
				rawValue:   entry.Value,
				enabled:    entry.Enabled,
				found:      entry.Found,
				err:        nil,
			}
		}
		c.metrics.IncCacheMiss()
	}

	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	value, enabled, found, err := c.evaluateWithRetries(ctx, featureKey, req)

	c.metrics.ObserveEvaluateLatency(time.Since(start))
	if err != nil {
		c.metrics.IncEvaluateError(getErrorCode(err))
	}

	if err == nil && c.cache != nil {
		c.cache.Set(cacheKey, value, enabled, found)
	}

	return EvalResult{
		featureKey: featureKey,
		rawValue:   value,
		enabled:    enabled,
		found:      found,
		err:        err,
	}
}

func (c *Client) Evaluate(featureKey string, req RequestContext) EvalResult {
	return c.EvaluateWithContext(context.Background(), featureKey, req)
}

func (c *Client) IsEnabled(featureKey string, req RequestContext) (bool, error) {
	res := c.Evaluate(featureKey, req)
	if err := res.Err(); err != nil {
		return false, err
	}

	if !res.Found() {
		return false, ErrFeatureNotFound
	}

	return res.Enabled(), nil
}

func (c *Client) IsEnabledOrDefault(featureKey string, req RequestContext, def bool) bool {
	enabled, err := c.IsEnabled(featureKey, req)
	if err != nil {
		c.logger.Warn("evaluation failed, using default",
			"feature_key", featureKey, "error", err, "default", def)

		return def
	}

	return enabled
}

func (c *Client) evaluateWithRetries(
	ctx context.Context,
	featureKey string,
	req RequestContext,
) (string, bool, bool, error) {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return "", false, false, ctx.Err()
			case <-time.After(delay):
			}
		}

		evalReq := make(api.EvaluateRequest, len(req))
		for k, v := range req {
			if raw, err := json.Marshal(v); err == nil {
				evalReq[k] = jx.Raw(raw)
			}
		}

		params := api.SdkV1FeaturesFeatureKeyEvaluatePostParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.SdkV1FeaturesFeatureKeyEvaluatePost(ctx, evalReq, params)
		if err == nil {
			switch r := resp.(type) {
			case *api.EvaluateResponse:
				return r.Value, r.Enabled, true, nil
			case *api.ErrorNotFound:
				return "", false, false, nil
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

		if !shouldRetry(err) {
			c.logger.Debug("not retrying due to error type", "error", err)

			break
		}

		c.logger.Debug("retrying due to error", "attempt", attempt, "error", err)
	}

	var apiErr *APIError
	if errors.As(lastErr, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return "", false, false, nil
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

func shouldRetry(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return err != nil
}

func getErrorCode(err error) string {
	return "unknown"
}
