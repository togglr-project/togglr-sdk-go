# Changelog

## [Unreleased] - 2025-10-02

### Added
- **Error Reporting**: New methods for reporting feature execution errors
  - `ReportError(ctx, featureKey, errorReport)` - Report a single error with automatic retries, returns `(health, isPending, error)`
  - `NewErrorReport(errorType, errorMessage)` - Create error report with builder pattern
  - `ErrorReport.WithContext(key, value)` - Add context data to error reports
  - `ErrorReport.WithContexts(contexts)` - Add multiple context data at once

- **Feature Health Monitoring**: New methods for monitoring feature health
  - `GetFeatureHealth(ctx, featureKey)` - Get detailed health status with automatic retries
  - `IsFeatureHealthy(ctx, featureKey)` - Simple boolean health check

- **New Types**:
  - `ErrorReport` - Structure for error reporting
  - `FeatureHealth` - Structure for health monitoring
  - Support for 202 responses with `isPending` boolean return value

- **Enhanced Examples**:
  - Updated simple example with error reporting and health monitoring
  - New advanced example demonstrating comprehensive usage
  - Separate example directories for better organization

### Changed
- **Retry Logic**: All methods now automatically apply retries based on client configuration
- **202 Response Handling**: 202 responses now return `(health, isPending, error)` with `isPending = true` instead of error
- **ReportError Signature**: Now returns `(health, isPending, error)` instead of `(health, error)`
- Updated README with comprehensive documentation for new features
- Enhanced error handling and response processing
- Improved example structure and organization

### Removed
- All `*WithRetries` methods - retries are now applied automatically
- `ErrPendingChange` error type - replaced with `isPending` boolean return value
- `PendingChange` field from `FeatureHealth` struct - replaced with `isPending` return value

### Technical Details
- All new methods follow existing SDK patterns and conventions
- Full integration with generated OpenAPI client
- Automatic retry logic based on client configuration
- Proper handling of 202 responses with pending change indication
- Backward compatible - no breaking changes to existing API
