// Example: Context and shared context usage
//
// This example demonstrates how to use context for structured logging,
// similar to Laravel's Log::withContext() and Log::shareContext().
//
// Run: go run main.go
package main

import (
	"fmt"

	"github.com/rabeeaali/golog"
)

func main() {
	// Initialize
	config := &golog.Config{
		Default: "file",
		AppName: "Context Example",
		Channels: map[string]golog.ChannelConfig{
			"file": golog.NewFileChannelConfig("logs/context.log"),
		},
	}

	if err := golog.Init(config); err != nil {
		panic(err)
	}
	defer golog.Close()

	// ========================================
	// Example 1: Share context globally
	// ========================================
	fmt.Println("🌍 Setting up shared context...")

	// This context will be included in ALL log entries
	golog.ShareContext(map[string]any{
		"app_version": "2.1.0",
		"environment": "production",
		"server_id":   "web-01",
		"region":      "us-east-1",
	})

	golog.Info("Application initialized")

	// ========================================
	// Example 2: Request-scoped context
	// ========================================
	fmt.Println("📝 Simulating request handling...")

	// Simulate handling a request
	handleRequest("req-abc123", 12345)
	handleRequest("req-xyz789", 67890)

	// ========================================
	// Example 3: Removing context
	// ========================================
	fmt.Println("🧹 Demonstrating context removal...")

	logger, _ := golog.Default()

	// Create logger with lots of context
	detailedLogger := logger.WithContext(map[string]any{
		"user_id":       123,
		"session_id":    "sess-abc",
		"request_id":    "req-123",
		"sensitive_key": "should-be-removed",
	})

	detailedLogger.Info("Before context removal")

	// Remove sensitive context
	cleanLogger := detailedLogger.WithoutContext("sensitive_key", "session_id")
	cleanLogger.Info("After context removal")

	fmt.Println("✅ Check logs/context.log for the output!")
}

// handleRequest simulates handling an HTTP request with request-scoped logging
func handleRequest(requestID string, userID int) {
	logger, _ := golog.Default()

	// Create request-scoped logger with context
	reqLogger := logger.WithContext(map[string]any{
		"request_id": requestID,
		"user_id":    userID,
	})

	// All logs in this request will have the request context
	reqLogger.Info("Request started")

	// Simulate some work
	reqLogger.Debug("Processing user data")
	reqLogger.Debug("Fetching from database")

	// Add more context as we learn more
	orderLogger := reqLogger.With("order_id", 999)
	orderLogger.Info("Processing order")

	reqLogger.Info("Request completed", map[string]any{
		"duration_ms": 150,
		"status":      200,
	})
}
