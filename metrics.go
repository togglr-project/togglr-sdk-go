package togglr

import "time"

// Metrics interface for collecting metrics
type Metrics interface {
	IncEvaluateRequest()
	IncEvaluateError(code string)
	ObserveEvaluateLatency(d time.Duration)
	IncCacheHit()
	IncCacheMiss()
	IncErrorReportRequest()
	IncErrorReportError(code string)
	ObserveErrorReportLatency(d time.Duration)
	IncFeatureHealthRequest()
	IncFeatureHealthError(code string)
	ObserveFeatureHealthLatency(d time.Duration)
}

// NoOpMetrics is a no-op implementation of Metrics
type NoOpMetrics struct{}

func (NoOpMetrics) IncEvaluateRequest()                         {}
func (NoOpMetrics) IncEvaluateError(code string)                {}
func (NoOpMetrics) ObserveEvaluateLatency(d time.Duration)      {}
func (NoOpMetrics) IncCacheHit()                                {}
func (NoOpMetrics) IncCacheMiss()                               {}
func (NoOpMetrics) IncErrorReportRequest()                      {}
func (NoOpMetrics) IncErrorReportError(code string)             {}
func (NoOpMetrics) ObserveErrorReportLatency(d time.Duration)   {}
func (NoOpMetrics) IncFeatureHealthRequest()                    {}
func (NoOpMetrics) IncFeatureHealthError(code string)           {}
func (NoOpMetrics) ObserveFeatureHealthLatency(d time.Duration) {}
