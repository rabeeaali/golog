// Example: Basic usage of golog
//
// This example demonstrates the basic usage of golog with file logging.
//
// Run: go run main.go
package main

import (
	"fmt"

	"github.com/rabeeaali/golog"
)

func main() {
	// Initialize with default configuration (logs to logs/app.log)
	if err := golog.Init(nil); err != nil {
		panic(err)
	}
	defer golog.Close()

	// Simple logging - like Laravel's Log facade
	golog.Debug("Application starting...")
	golog.Info("User logged in", map[string]any{
		"user_id":  123,
		"username": "ahmed",
		"ip":       "192.168.1.1",
	})

	// Warning log
	golog.Warning("API rate limit approaching", map[string]any{
		"current_rate": 950,
		"max_rate":     1000,
	})

	// Error log
	golog.Error("Failed to process payment", map[string]any{
		"order_id": 456,
		"error":    "Card declined",
	})

	// Error with exception details
	err := fmt.Errorf("database connection timeout after 30s")
	golog.ErrorWithException("Database error", err, map[string]any{
		"host":     "db.example.com",
		"database": "production",
	})

	fmt.Println("✅ Logs written to logs/app.log")
}

