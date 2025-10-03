# Changelog

## [Unreleased] - 2025-01-02

### Changed
- **Error Reporting API Simplification**: Updated error reporting to use asynchronous processing
  - `ReportError(ctx, featureKey, errorReport)` - Now returns only `error` (simplified API)
  - 202 responses now always indicate successful queuing for processing (no more pending changes)
  - Removed `isPending` return value as it's no longer needed
- **Metrics Interface Enhancement**: Added dedicated metrics for error reporting and feature health
  - Added `IncErrorReportRequest()`, `IncErrorReportError(code)`, `ObserveErrorReportLatency(d)`
  - Added `IncFeatureHealthRequest()`, `IncFeatureHealthError(code)`, `ObserveFeatureHealthLatency(d)`
  - Fixed incorrect use of evaluation metrics in error reporting methods

### Added
- **Error Reporting**: New methods for reporting feature execution errors
  - `ReportError(ctx, featureKey, errorReport)` - Report a single error with automatic retries, returns `error`
  - `NewErrorReport(errorType, errorMessage)` - Create error report with builder pattern
  - `ErrorReport.WithContext(key, value)` - Add context data to error reports
  - `ErrorReport.WithContexts(contexts)` - Add multiple context data at once

- **Feature Health Monitoring**: New methods for monitoring feature health
  - `GetFeatureHealth(ctx, featureKey)` - Get detailed health status with automatic retries
  - `IsFeatureHealthy(ctx, featureKey)` - Simple boolean health check

- **New Types**:
  - `ErrorReport` - Structure for error reporting
  - `FeatureHealth` - Structure for health monitoring

- **Enhanced Examples**:
  - Updated simple example with error reporting and health monitoring
  - New advanced example demonstrating comprehensive usage
  - Separate example directories for better organization

### Changed
- **Retry Logic**: All methods now automatically apply retries based on client configuration
- Updated README with comprehensive documentation for new features
- Enhanced error handling and response processing
- Improved example structure and organization

### Removed
- All `*WithRetries` methods - retries are now applied automatically
- `isPending` return value from `ReportError` - simplified API with asynchronous processing

### Technical Details
- All new methods follow existing SDK patterns and conventions
- Full integration with generated OpenAPI client
- Automatic retry logic based on client configuration
- Proper handling of 202 responses with pending change indication
- Backward compatible - no breaking changes to existing API
