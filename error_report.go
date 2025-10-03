package togglr

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-faster/jx"

	api "github.com/togglr-project/togglr-sdk-go/internal/generated/client"
)

// ErrorReport represents a feature error report
type ErrorReport struct {
	ErrorType    string
	ErrorMessage string
	Context      map[string]any
}

// NewErrorReport creates a new error report
func NewErrorReport(errorType, errorMessage string) *ErrorReport {
	return &ErrorReport{
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
		Context:      make(map[string]any),
	}
}

// WithContext adds context data to the error report
func (er *ErrorReport) WithContext(key string, value any) *ErrorReport {
	er.Context[key] = value
	return er
}

// WithContexts adds multiple context data to the error report
func (er *ErrorReport) WithContexts(contexts map[string]any) *ErrorReport {
	for k, v := range contexts {
		er.Context[k] = v
	}
	return er
}

// ReportError sends an error report for a feature
func (c *Client) ReportError(
	ctx context.Context,
	featureKey string,
	report *ErrorReport,
) error {
	start := time.Now()
	c.metrics.IncErrorReportRequest()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	// Make API call with retries
	err := c.reportErrorWithRetries(ctx, featureKey, report)

	// Record metrics
	c.metrics.ObserveErrorReportLatency(time.Since(start))
	if err != nil {
		c.metrics.IncErrorReportError("report_error_failed")
	}

	return err
}

// reportErrorWithRetries performs the actual API call with retry logic
func (c *Client) reportErrorWithRetries(
	ctx context.Context,
	featureKey string,
	report *ErrorReport,
) error {
	var lastErr error

	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			// Calculate backoff delay
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Debug("retrying error report after delay", "attempt", attempt, "delay", delay)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		// Convert to API request
		apiReq := &api.FeatureErrorReport{
			ErrorType:    report.ErrorType,
			ErrorMessage: report.ErrorMessage,
		}

		// Convert context to jx.Raw if provided
		if len(report.Context) > 0 {
			contextData := make(api.FeatureErrorReportContext)
			for k, v := range report.Context {
				if raw, err := json.Marshal(v); err == nil {
					contextData[k] = jx.Raw(raw)
				}
			}
			apiReq.Context = api.NewOptFeatureErrorReportContext(contextData)
		}

		// Make API call
		params := api.ReportFeatureErrorParams{
			FeatureKey: featureKey,
		}

		resp, err := c.apiClient.ReportFeatureError(ctx, apiReq, params)
		if err == nil {
			// Handle response
			switch resp.(type) {
			case *api.ReportFeatureErrorAccepted:
				// 202 response - error reported successfully, queued for processing
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

		// Check if we should retry
		if !shouldRetry(err) {
			c.logger.Debug("not retrying error report due to error type", "error", err)
			break
		}

		c.logger.Debug("retrying error report due to error", "attempt", attempt, "error", err)
	}

	return lastErr
}
