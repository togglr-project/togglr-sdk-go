package togglr

import (
	"context"
	"fmt"
	"time"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

func (c *Client) TrackEvent(
	ctx context.Context,
	featureKey string,
	event *TrackEvent,
) error {
	start := time.Now()
	c.metrics.IncTrackEventRequest()

	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	err := c.trackEventWithRetries(ctx, featureKey, event)

	c.metrics.ObserveTrackEventLatency(time.Since(start))
	if err != nil {
		c.metrics.IncTrackEventError("track_event_failed")
	}

	return err
}

func (c *Client) trackEventWithRetries(
	ctx context.Context,
	featureKey string,
	event *TrackEvent,
) error {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying track event after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		apiReq := event.toAPIRequest()

		params := api.TrackFeatureEventParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.TrackFeatureEvent(ctx, apiReq, params)
		if err == nil {
			switch resp.(type) {
			case *api.TrackFeatureEventAccepted:
				return nil
			case *api.ErrorBadRequest:
				return ErrBadRequest
			case *api.ErrorUnauthorized:
				return ErrUnauthorized
			case *api.ErrorNotFound:
				return ErrFeatureNotFound
			case *api.ErrorTooManyRequests:
				return ErrTooManyRequests
			case *api.ErrorInternalServerError:
				return ErrInternalServerError
			default:
				return fmt.Errorf("unexpected response type: %T", resp)
			}
		}

		lastErr = err

		if !shouldRetry(err) {
			c.logger.Debug("not retrying track event due to error type", "error", err)
			break
		}

		c.logger.Debug("retrying track event due to error", "attempt", attempt, "error", err)
	}

	return lastErr
}
