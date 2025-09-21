package main

import (
	"context"
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
		togglr.WithRetries(3),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Build request context using builder methods
	reqCtx := togglr.NewContext().
		WithUserID("user123").
		WithCountry("US").
		WithUserEmail("user@example.com").
		WithDeviceType("mobile").
		WithOS("iOS").
		WithOSVersion("15.0")

	// Evaluate a feature
	featureKey := "new_ui"
	value, enabled, found, err := client.Evaluate(featureKey, reqCtx)
	if err != nil {
		log.Printf("Error evaluating feature %s: %v", featureKey, err)
		return
	}

	if !found {
		fmt.Printf("Feature %s not found\n", featureKey)
		return
	}

	fmt.Printf("Feature %s: enabled=%t, value=%s\n", featureKey, enabled, value)

	// Use convenience method for boolean flags
	isEnabled, err := client.IsEnabled(featureKey, reqCtx)
	if err != nil {
		log.Printf("Error checking if feature is enabled: %v", err)
		return
	}

	fmt.Printf("Feature %s is enabled: %t\n", featureKey, isEnabled)

	// Use default value fallback
	isEnabled = client.IsEnabledOrDefault(featureKey, reqCtx, false)
	fmt.Printf("Feature %s with default fallback: %t\n", featureKey, isEnabled)

	// Example with context cancellation
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	value2, enabled2, found2, err := client.EvaluateWithContext(ctxWithTimeout, "another_feature", reqCtx)
	if err != nil {
		log.Printf("Error with context: %v", err)
		return
	}

	fmt.Printf("Another feature: enabled=%t, value=%s, found=%t\n", enabled2, value2, found2)

	// Health check
	if err := client.HealthCheck(context.Background()); err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Println("Health check passed")
	}
}
