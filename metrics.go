package togglr

import "time"

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
	IncTrackEventRequest()
	IncTrackEventError(code string)
	ObserveTrackEventLatency(d time.Duration)
}

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
func (NoOpMetrics) IncTrackEventRequest()                       {}
func (NoOpMetrics) IncTrackEventError(code string)              {}
func (NoOpMetrics) ObserveTrackEventLatency(d time.Duration)    {}
