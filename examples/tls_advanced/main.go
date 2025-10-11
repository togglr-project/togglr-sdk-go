package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/togglr-project/togglr-sdk-go"
)

func main() {
	fmt.Println("=== TLS Configuration Examples ===")

	// Example 1: Basic TLS with client certificate
	fmt.Println("\n1. Basic TLS with client certificate:")
	client1, err := togglr.NewClientWithDefaults("42b6f8f1-630c-400c-97bd-a3454a07f700",
		togglr.WithBaseURL("https://localhost"),
		togglr.WithClientCertAndKey("/path/to/client.crt", "/path/to/client.key"),
		togglr.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create client 1: %v", err)
	} else {
		defer client1.Close()
		fmt.Println("Client 1 created successfully with client certificate")
	}

	// Example 2: TLS with custom CA certificate
	fmt.Println("\n2. TLS with custom CA certificate:")
	client2, err := togglr.NewClientWithDefaults("your-api-key-here",
		togglr.WithBaseURL("https://localhost:8090"),
		togglr.WithCACert("/path/to/ca.crt"),
		togglr.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create client 2: %v", err)
	} else {
		defer client2.Close()
		fmt.Println("Client 2 created successfully with custom CA certificate")
	}

	// Example 3: Full TLS configuration (client cert + CA cert)
	fmt.Println("\n3. Full TLS configuration:")
	client3, err := togglr.NewClientWithDefaults("your-api-key-here",
		togglr.WithBaseURL("https://localhost:8090"),
		togglr.WithClientCertAndKey("/path/to/client.crt", "/path/to/client.key"),
		togglr.WithCACert("/path/to/ca.crt"),
		togglr.WithTimeout(5*time.Second),
		togglr.WithRetries(3),
	)
	if err != nil {
		log.Printf("Failed to create client 3: %v", err)
	} else {
		defer client3.Close()
		fmt.Println("Client 3 created successfully with full TLS configuration")
	}

	// Example 4: Insecure mode (skip TLS verification)
	fmt.Println("\n4. Insecure mode (skip TLS verification):")
	client4, err := togglr.NewClientWithDefaults("your-api-key-here",
		togglr.WithBaseURL("https://localhost:8090"),
		togglr.WithInsecure(),
		togglr.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create client 4: %v", err)
	} else {
		defer client4.Close()
		fmt.Println("Client 4 created successfully in insecure mode")
	}

	// Example 5: Using individual certificate options
	fmt.Println("\n5. Using individual certificate options:")
	client5, err := togglr.NewClientWithDefaults("your-api-key-here",
		togglr.WithBaseURL("https://localhost:8090"),
		togglr.WithClientCert("/path/to/client.crt"),
		togglr.WithClientKey("/path/to/client.key"),
		togglr.WithCACert("/path/to/ca.crt"),
		togglr.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create client 5: %v", err)
	} else {
		defer client5.Close()
		fmt.Println("Client 5 created successfully with individual certificate options")
	}

	// Example 6: Test feature evaluation with TLS client
	fmt.Println("\n6. Testing feature evaluation with TLS client:")
	if client3 != nil {
		reqCtx := togglr.NewContext().
			WithUserID("user123").
			WithCountry("US").
			WithUserEmail("user@example.com").
			WithDeviceType("mobile").
			WithOS("iOS").
			WithOSVersion("15.0")

		featureKey := "new_ui"
		res := client3.Evaluate(featureKey, reqCtx)
		if err := res.Err(); err != nil {
			log.Printf("Error evaluating feature %s: %v", featureKey, err)
		} else if res.Found() {
			fmt.Printf("Feature %s: enabled=%t, value=%s\n", featureKey, res.Enabled(), res.Value())
		} else {
			fmt.Printf("Feature %s not found\n", featureKey)
		}
	}

	// Example 7: Health check with TLS client
	fmt.Println("\n7. Testing health check with TLS client:")
	if client3 != nil {
		if err := client3.HealthCheck(context.Background()); err != nil {
			log.Printf("Health check failed: %v", err)
		} else {
			fmt.Println("Health check passed")
		}
	}

	// Example 8: Error reporting with TLS client
	fmt.Println("\n8. Testing error reporting with TLS client:")
	if client3 != nil {
		errorReport := togglr.NewErrorReport("tls_test", "Testing error reporting with TLS").
			WithContext("tls_enabled", true).
			WithContext("client_cert", "/path/to/client.crt").
			WithContext("ca_cert", "/path/to/ca.crt")

		err := client3.ReportError(context.Background(), "test_feature", errorReport)
		if err != nil {
			log.Printf("Error reporting failed: %v", err)
		} else {
			fmt.Println("Error reported successfully with TLS client")
		}
	}

	fmt.Println("\n=== TLS Examples Completed ===")
}
