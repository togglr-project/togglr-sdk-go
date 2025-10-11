package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/togglr-project/togglr-sdk-go"
)

func main() {
	// Create client with advanced configuration
	client, err := togglr.NewClientWithDefaults("42b6f8f1-630c-400c-97bd-a3454a07f700",
		togglr.WithBaseURL("https://localhost"),
		togglr.WithTimeout(2*time.Second),
		togglr.WithCache(1000, 30*time.Second),
		togglr.WithRetries(3),
		togglr.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	featureKey := "payment_processing"

	// Example 1: Basic feature evaluation
	fmt.Println("=== Basic Feature Evaluation ===")
	reqCtx := togglr.NewContext().
		WithUserID("user123").
		WithCountry("US").
		WithUserEmail("user@example.com")

	res := client.Evaluate(featureKey, reqCtx)
	if err := res.Err(); err != nil {
		log.Printf("Error evaluating feature: %v", err)
		return
	}

	if res.Found() {
		fmt.Printf("Feature %s: enabled=%t, value=%s\n", featureKey, res.Enabled(), res.Value())
	} else {
		fmt.Printf("Feature %s not found\n", featureKey)
	}

	// Example 2: Error reporting with different error types
	fmt.Println("\n=== Error Reporting Examples ===")

	// Report a timeout error
	timeoutError := togglr.NewErrorReport("timeout", "Payment gateway timeout").
		WithContext("gateway", "stripe").
		WithContext("timeout_ms", 5000).
		WithContext("retry_count", 2).
		WithContext("request_id", "req_12345")

	err = client.ReportError(context.Background(), featureKey, timeoutError)
	if err != nil {
		log.Printf("Error reporting timeout: %v", err)
	} else {
		fmt.Printf("Timeout error reported successfully - queued for processing\n")
	}

	// Report a validation error
	validationError := togglr.NewErrorReport("validation", "Invalid payment data").
		WithContext("field", "card_number").
		WithContext("error_code", "INVALID_FORMAT").
		WithContext("user_id", "user123")

	err = client.ReportError(context.Background(), featureKey, validationError)
	if err != nil {
		log.Printf("Error reporting validation error: %v", err)
	} else {
		fmt.Printf("Validation error reported successfully - queued for processing\n")
	}

	// Report a service unavailable error
	serviceError := togglr.NewErrorReport("service_unavailable", "Payment service is down").
		WithContext("service", "payment-processor").
		WithContext("status_code", 503).
		WithContext("region", "us-east-1").
		WithContext("timestamp", time.Now().Unix())

	err = client.ReportError(context.Background(), featureKey, serviceError)
	if err != nil {
		log.Printf("Error reporting service error: %v", err)
	} else {
		fmt.Printf("Service error reported successfully - queued for processing\n")
	}

	// Example 3: Feature health monitoring
	fmt.Println("\n=== Feature Health Monitoring ===")

	// Get current health status
	featureHealth, err := client.GetFeatureHealth(context.Background(), featureKey)
	if err != nil {
		log.Printf("Error getting feature health: %v", err)
	} else {
		fmt.Printf("Feature Health Status:\n")
		fmt.Printf("  Feature Key: %s\n", featureHealth.FeatureKey)
		fmt.Printf("  Environment: %s\n", featureHealth.EnvironmentKey)
		fmt.Printf("  Enabled: %t\n", featureHealth.Enabled)
		fmt.Printf("  Auto Disabled: %t\n", featureHealth.AutoDisabled)

		if featureHealth.ErrorRate != nil {
			fmt.Printf("  Error Rate: %.2f%%\n", *featureHealth.ErrorRate*100)
		}

		if featureHealth.Threshold != nil {
			fmt.Printf("  Threshold: %.2f%%\n", *featureHealth.Threshold*100)
		}

		if featureHealth.LastErrorAt != nil {
			fmt.Printf("  Last Error: %s\n", featureHealth.LastErrorAt.Format(time.RFC3339))
		}
	}

	// Check if feature is healthy
	isHealthy, err := client.IsFeatureHealthy(context.Background(), featureKey)
	if err != nil {
		log.Printf("Error checking feature health: %v", err)
	} else {
		fmt.Printf("Feature is healthy: %t\n", isHealthy)
	}

	// Example 4: Batch error reporting simulation
	fmt.Println("\n=== Batch Error Reporting Simulation ===")

	// Simulate multiple errors to trigger auto-disable
	for i := 0; i < 5; i++ {
		batchError := togglr.NewErrorReport("batch_error", fmt.Sprintf("Batch error #%d", i+1)).
			WithContext("batch_id", "batch_123").
			WithContext("error_index", i).
			WithContext("timestamp", time.Now().Unix())

		err := client.ReportError(context.Background(), featureKey, batchError)
		if err != nil {
			log.Printf("Error reporting batch error %d: %v", i+1, err)
		} else {
			fmt.Printf("Batch error %d reported successfully - queued for processing\n", i+1)
		}

		// Small delay between errors
		time.Sleep(100 * time.Millisecond)
	}

	// Final health check
	fmt.Println("\n=== Final Health Check ===")
	finalHealth, err := client.GetFeatureHealth(context.Background(), featureKey)
	if err != nil {
		log.Printf("Error getting final health: %v", err)
	} else {
		fmt.Printf("Final Status - enabled=%t, auto_disabled=%t\n",
			finalHealth.Enabled, finalHealth.AutoDisabled)

		if finalHealth.ErrorRate != nil {
			fmt.Printf("Final Error Rate: %.2f%%\n", *finalHealth.ErrorRate*100)
		}
	}

	// Example 5: Error reporting with context timeout
	fmt.Println("\n=== Error Reporting with Context Timeout ===")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contextTimeoutError := togglr.NewErrorReport("timeout_test", "Testing with context timeout").
		WithContext("test_id", "timeout_001").
		WithContext("timeout_sec", 5)

	err = client.ReportError(ctx, featureKey, contextTimeoutError)
	if err != nil {
		log.Printf("Error with timeout context: %v", err)
	} else {
		fmt.Printf("Timeout test completed - error reported successfully - queued for processing\n")
	}

	// Example 6: Health check with context timeout
	fmt.Println("\n=== Health Check with Context Timeout ===")

	health, err := client.GetFeatureHealth(ctx, featureKey)
	if err != nil {
		log.Printf("Error getting health with timeout: %v", err)
	} else {
		fmt.Printf("Health with timeout - enabled=%t, auto_disabled=%t\n",
			health.Enabled, health.AutoDisabled)
	}

	// Example 7: Health check with context timeout
	fmt.Println("\n=== Is Healthy with Context Timeout ===")

	isHealthy, err = client.IsFeatureHealthy(ctx, featureKey)
	if err != nil {
		log.Printf("Error checking health with timeout: %v", err)
	} else {
		fmt.Printf("Is healthy with timeout: %t\n", isHealthy)
	}
}
