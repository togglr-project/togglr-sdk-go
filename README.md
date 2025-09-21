# Togglr Go SDK

Go SDK for working with Togglr - feature flag management system.

## Installation

```bash
go get github.com/rom8726/togglr-sdk-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/rom8726/togglr-sdk-go"
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
    value, enabled, found, err := client.Evaluate("new_ui", ctx)
    if err != nil {
        log.Fatal(err)
    }

    if found {
        fmt.Printf("Feature enabled: %t, value: %s\n", enabled, value)
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
value, enabled, found, err := client.Evaluate("feature_key", ctx)

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

value, enabled, found, err := client.EvaluateWithContext(ctx, "api-key", "feature_key", reqCtx)
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
value, enabled, found, err := client.Evaluate("feature_key", ctx)
if err != nil {
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
