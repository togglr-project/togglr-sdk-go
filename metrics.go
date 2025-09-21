package togglr

import "time"

// Metrics interface for collecting metrics
type Metrics interface {
	IncEvaluateRequest()
	IncEvaluateError(code string)
	ObserveEvaluateLatency(d time.Duration)
	IncCacheHit()
	IncCacheMiss()
}

// NoOpMetrics is a no-op implementation of Metrics
type NoOpMetrics struct{}

func (NoOpMetrics) IncEvaluateRequest()                    {}
func (NoOpMetrics) IncEvaluateError(code string)           {}
func (NoOpMetrics) ObserveEvaluateLatency(d time.Duration) {}
func (NoOpMetrics) IncCacheHit()                           {}
func (NoOpMetrics) IncCacheMiss()                          {}
