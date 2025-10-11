package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/togglr-project/togglr-sdk-go"
)

func main() {
	// Example of using TLS certificates for secure connection
	client, err := togglr.NewClientWithDefaults("42b6f8f1-630c-400c-97bd-a3454a07f700",
		togglr.WithBaseURL("https://localhost"),
		// Use client certificate and key for mutual TLS authentication
		togglr.WithClientCertAndKey("/path/to/client.crt", "/path/to/client.key"),
		// Use custom CA certificate for server verification
		//togglr.WithCACert("/path/to/ca.crt"),
		togglr.WithTimeout(5*time.Second),
		togglr.WithRetries(3),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Build request context
	reqCtx := togglr.NewContext().
		WithUserID("user123").
		WithCountry("US").
		WithUserEmail("user@example.com").
		WithDeviceType("mobile").
		WithOS("iOS").
		WithOSVersion("15.0")

	// Evaluate a feature
	featureKey := "new_ui"
	res := client.Evaluate(featureKey, reqCtx)
	if err := res.Err(); err != nil {
		log.Printf("Error evaluating feature %s: %v", featureKey, err)
		return
	}

	if !res.Found() {
		fmt.Printf("Feature %s not found\n", featureKey)
		return
	}

	fmt.Printf("Feature %s: enabled=%t, value=%s\n", featureKey, res.Enabled(), res.Value())

	// Use convenience method for boolean flags
	isEnabled, err := client.IsEnabled(featureKey, reqCtx)
	if err != nil {
		log.Printf("Error checking if feature is enabled: %v", err)
		return
	}

	fmt.Printf("Feature %s is enabled: %t\n", featureKey, isEnabled)

	// Health check
	if err := client.HealthCheck(context.Background()); err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Println("Health check passed")
	}

	// Example: Report an error for a feature
	errorReport := togglr.NewErrorReport("timeout", "Service did not respond in 5s").
		WithContext("service", "payment-gateway").
		WithContext("timeout_ms", 5000).
		WithContext("retry_count", 3)

	err = client.ReportError(context.Background(), featureKey, errorReport)
	if err != nil {
		log.Printf("Error reporting feature error: %v", err)
	} else {
		fmt.Printf("Error reported successfully - queued for processing\n")
	}

	// Example: Get feature health status
	featureHealth, err := client.GetFeatureHealth(context.Background(), featureKey)
	if err != nil {
		log.Printf("Error getting feature health: %v", err)
	} else {
		fmt.Printf("Feature health: enabled=%t, auto_disabled=%t\n",
			featureHealth.Enabled, featureHealth.AutoDisabled)
		if featureHealth.ErrorRate != nil {
			fmt.Printf("Error rate: %.2f%%\n", *featureHealth.ErrorRate*100)
		}
		if featureHealth.LastErrorAt != nil {
			fmt.Printf("Last error at: %s\n", featureHealth.LastErrorAt.Format(time.RFC3339))
		}
	}

	// Example: Check if feature is healthy
	isHealthy, err := client.IsFeatureHealthy(context.Background(), featureKey)
	if err != nil {
		log.Printf("Error checking feature health: %v", err)
	} else {
		fmt.Printf("Feature is healthy: %t\n", isHealthy)
	}
}
