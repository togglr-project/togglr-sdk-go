# Togglr Go SDK

Go SDK for working with Togglr - feature flag management system.

## Installation

```bash
go get github.com/togglr-project/togglr-sdk-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/togglr-project/togglr-sdk-go"
)

func main() {
    // Create client with default configuration
    client, err := togglr.NewClientWithDefaults("your-api-key-here",
        togglr.WithBaseURL("http://localhost:8090"),
        togglr.WithTimeout(1*time.Second),
        togglr.WithCache(1000, 10*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create request context
    ctx := togglr.NewContext().
        WithUserID("user123").
        WithCountry("US").
        WithDeviceType("mobile")

    // Evaluate feature flag
    res := client.Evaluate("new_ui", ctx)
    if err:= res.Err(); err != nil {
        log.Fatal(err)
    }

    if res.Found() {
        fmt.Printf("Feature enabled: %t, value: %s\n", res.Enabled(), res.Value())
    }
}
```

## Configuration

### Creating a client

```go
// With default settings
client, err := togglr.NewClientWithDefaults("api-key")

// With custom configuration
cfg := togglr.DefaultConfig("api-key")
cfg.BaseURL = "https://api.togglr.com"
cfg.Timeout = 2 * time.Second
cfg.Retries = 3

client, err := togglr.NewClient(cfg)
```

### Functional options

```go
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithBaseURL("https://api.togglr.com"),
    togglr.WithTimeout(2*time.Second),
    togglr.WithRetries(3),
    togglr.WithCache(1000, 10*time.Second),
    togglr.WithCircuitBreaker(true),
)
```

## Usage

### Creating request context

```go
ctx := togglr.NewContext().
    WithUserID("user123").
    WithUserEmail("user@example.com").
    WithCountry("US").
    WithDeviceType("mobile").
    WithOS("iOS").
    WithOSVersion("15.0").
    WithBrowser("Safari").
    WithLanguage("en-US")
```

### Evaluating feature flags

```go
// Full evaluation
res := client.Evaluate("feature_key", ctx)

// Simple enabled check
isEnabled, err := client.IsEnabled("feature_key", ctx)

// With default value
isEnabled = client.IsEnabledOrDefault("feature_key", ctx, false)
```

### Working with context

```go
// With cancellation context
ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
defer cancel()

res := client.EvaluateWithContext(ctx, "feature_key", reqCtx)
```

### Error Reporting and Auto-Disable

The SDK supports reporting feature execution errors for auto-disable functionality:

```go
// Create an error report
errorReport := togglr.NewErrorReport("timeout", "Service did not respond in 5s").
    WithContext("service", "payment-gateway").
    WithContext("timeout_ms", 5000).
    WithContext("retry_count", 3)

// Report the error
health, isPending, err := client.ReportError(context.Background(), "feature_key", errorReport)
if err != nil {
    log.Printf("Error reporting: %v", err)
} else {
    fmt.Printf("Feature enabled: %t, auto_disabled: %t, pending_change: %t\n", 
        health.Enabled, health.AutoDisabled, isPending)
    
    if isPending {
        fmt.Println("Change is pending approval")
    }
}
```

### Feature Health Monitoring

Check the health status of features:

```go
// Get feature health status
health, err := client.GetFeatureHealth(context.Background(), "feature_key")
if err != nil {
    log.Printf("Error getting health: %v", err)
} else {
    fmt.Printf("Feature Health:\n")
    fmt.Printf("  Enabled: %t\n", health.Enabled)
    fmt.Printf("  Auto Disabled: %t\n", health.AutoDisabled)
    if health.ErrorRate != nil {
        fmt.Printf("  Error Rate: %.2f%%\n", *health.ErrorRate*100)
    }
    if health.LastErrorAt != nil {
        fmt.Printf("  Last Error: %s\n", health.LastErrorAt.Format(time.RFC3339))
    }
}

// Simple health check
isHealthy, err := client.IsFeatureHealthy(context.Background(), "feature_key")
if err != nil {
    log.Printf("Error checking health: %v", err)
} else {
    fmt.Printf("Feature is healthy: %t\n", isHealthy)
}
```


## Caching

The SDK supports optional caching of evaluation results:

```go
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithCache(1000, 10*time.Second), // cache size and TTL
)
```

## Retries

The SDK automatically retries requests on temporary errors:

```go
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithRetries(3), // number of attempts
    togglr.WithBackoff(togglr.Backoff{
        BaseDelay: 100 * time.Millisecond,
        MaxDelay:  2 * time.Second,
        Factor:    2.0,
    }),
)
```

## Logging and Metrics

```go
// Custom logger
logger := &MyLogger{}
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithLogger(logger),
)

// Custom metrics
metrics := &MyMetrics{}
client, err := togglr.NewClientWithDefaults("api-key",
    togglr.WithMetrics(metrics),
)
```

## Error Handling

```go
res := client.Evaluate("feature_key", ctx)
if err := res.Err(); err != nil {
    switch {
    case errors.Is(err, togglr.ErrUnauthorized):
        // Authorization error
    case errors.Is(err, togglr.ErrBadRequest):
        // Bad request
    default:
        // Other error
    }
}
```

### Error Report Types

```go
// Create different types of error reports
timeoutError := togglr.NewErrorReport("timeout", "Service timeout").
    WithContext("service", "payment-gateway").
    WithContext("timeout_ms", 5000)

validationError := togglr.NewErrorReport("validation", "Invalid data").
    WithContext("field", "email").
    WithContext("error_code", "INVALID_FORMAT")

serviceError := togglr.NewErrorReport("service_unavailable", "Service down").
    WithContext("service", "database").
    WithContext("status_code", 503)
```

### Feature Health Types

```go
// FeatureHealth provides detailed health information
type FeatureHealth struct {
    FeatureKey     string     // Feature identifier
    EnvironmentKey string     // Environment identifier
    Enabled        bool       // Whether feature is enabled
    AutoDisabled   bool       // Whether feature was auto-disabled
    ErrorRate      *float32   // Error rate percentage (optional)
    Threshold      *float32   // Error threshold (optional)
    LastErrorAt    *time.Time // Last error timestamp (optional)
}
```

## Client Generation

To update the generated client from OpenAPI specification:

```bash
make generate
```

## Building and Testing

```bash
# Build
make build

# Testing
make test

# Linting
make lint

# Clean
make clean
```

## Examples

Complete usage examples are located in `examples/`.
