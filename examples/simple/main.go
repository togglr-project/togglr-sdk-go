package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/togglr-project/togglr-sdk-go"
)

func main() {
	// Create client with default configuration
	client, err := togglr.NewClientWithDefaults("42b6f8f1-630c-400c-97bd-a3454a07f700",
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

	// Use default value fallback
	isEnabled = client.IsEnabledOrDefault(featureKey, reqCtx, false)
	fmt.Printf("Feature %s with default fallback: %t\n", featureKey, isEnabled)

	// Example with context cancellation
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	res = client.EvaluateWithContext(ctxWithTimeout, "another_feature", reqCtx)
	if err := res.Err(); err != nil {
		log.Printf("Error with context: %v", err)
		return
	}

	fmt.Printf("Another feature: enabled=%t, value=%s, found=%t\n", res.Enabled(), res.Value(), res.Found())

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

	// Example: Track events for analytics
	// Track impression event (recommended for each evaluation)
	impressionEvent := togglr.NewTrackEvent("A", togglr.EventTypeSuccess).
		WithContext("user.id", "user123").
		WithContext("country", "US").
		WithContext("device_type", "mobile").
		WithDedupKey("impression-user123-new_ui")

	err = client.TrackEvent(context.Background(), featureKey, impressionEvent)
	if err != nil {
		log.Printf("Error tracking impression event: %v", err)
	} else {
		fmt.Printf("Impression event tracked successfully\n")
	}

	// Track conversion event with reward
	conversionEvent := togglr.NewTrackEvent("A", togglr.EventTypeSuccess).
		WithReward(1.0).
		WithContext("user.id", "user123").
		WithContext("conversion_type", "purchase").
		WithContext("order_value", 99.99).
		WithDedupKey("conversion-user123-new_ui")

	err = client.TrackEvent(context.Background(), featureKey, conversionEvent)
	if err != nil {
		log.Printf("Error tracking conversion event: %v", err)
	} else {
		fmt.Printf("Conversion event tracked successfully\n")
	}

	// Track error event
	errorEvent := togglr.NewTrackEvent("B", togglr.EventTypeError).
		WithContext("user.id", "user123").
		WithContext("error_type", "timeout").
		WithContext("error_message", "Service did not respond in 5s").
		WithDedupKey("error-user123-new_ui")

	err = client.TrackEvent(context.Background(), featureKey, errorEvent)
	if err != nil {
		log.Printf("Error tracking error event: %v", err)
	} else {
		fmt.Printf("Error event tracked successfully\n")
	}
}
