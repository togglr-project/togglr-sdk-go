package togglr

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-faster/jx"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

type ErrorReport struct {
	ErrorType    string
	ErrorMessage string
	Context      map[string]any
}

func NewErrorReport(errorType, errorMessage string) *ErrorReport {
	return &ErrorReport{
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
		Context:      make(map[string]any),
	}
}

func (er *ErrorReport) WithContext(key string, value any) *ErrorReport {
	er.Context[key] = value
	return er
}

func (er *ErrorReport) WithContexts(contexts map[string]any) *ErrorReport {
	for k, v := range contexts {
		er.Context[k] = v
	}
	return er
}

func (c *Client) ReportError(
	ctx context.Context,
	featureKey string,
	report *ErrorReport,
) error {
	start := time.Now()
	c.metrics.IncErrorReportRequest()

	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	err := c.reportErrorWithRetries(ctx, featureKey, report)

	c.metrics.ObserveErrorReportLatency(time.Since(start))
	if err != nil {
		c.metrics.IncErrorReportError("report_error_failed")
	}

	return err
}

func (c *Client) reportErrorWithRetries(
	ctx context.Context,
	featureKey string,
	report *ErrorReport,
) error {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying error report after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		apiReq := &api.FeatureErrorReport{
			ErrorType:    report.ErrorType,
			ErrorMessage: report.ErrorMessage,
		}

		if len(report.Context) > 0 {
			contextData := make(api.FeatureErrorReportContext)
			for k, v := range report.Context {
				if raw, err := json.Marshal(v); err == nil {
					contextData[k] = jx.Raw(raw)
				}
			}
			apiReq.Context = api.NewOptFeatureErrorReportContext(contextData)
		}

		params := api.ReportFeatureErrorParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.ReportFeatureError(ctx, apiReq, params)
		if err == nil {
			switch resp.(type) {
			case *api.ReportFeatureErrorAccepted:
				return nil
			case *api.ErrorBadRequest:
				return ErrBadRequest
			case *api.ErrorUnauthorized:
				return ErrUnauthorized
			case *api.ErrorNotFound:
				return ErrFeatureNotFound
			case *api.ErrorInternalServerError:
				return ErrInternalServerError
			default:
				return fmt.Errorf("unexpected response type: %T", resp)
			}
		}

		lastErr = err

		if !shouldRetry(err) {
			c.logger.Debug("not retrying error report due to error type", "error", err)
			break
		}

		c.logger.Debug("retrying error report due to error", "attempt", attempt, "error", err)
	}

	return lastErr
}
