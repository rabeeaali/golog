// Example: Slack logging with multiple channels
//
// This example demonstrates how to configure multiple Slack channels
// for different types of notifications (orders, payments, errors).
//
// Set environment variables before running:
//
//	export SLACK_ORDERS_WEBHOOK="https://hooks.slack.com/services/..."
//	export SLACK_PAYMENTS_WEBHOOK="https://hooks.slack.com/services/..."
//	export SLACK_ERRORS_WEBHOOK="https://hooks.slack.com/services/..."
//
// Run: go run main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golog-pkg/golog"
)

func main() {
	// Configuration with multiple Slack channels
	config := &golog.Config{
		Default: "file",
		AppName: "E-Commerce App",
		Channels: map[string]golog.ChannelConfig{
			// Default file logging
			"file": golog.NewFileChannelConfig("logs/app.log"),

			// Slack channel for order notifications
			"slack-orders": golog.NewSlackChannelConfig(
				getEnvOrDefault("SLACK_ORDERS_WEBHOOK", "https://hooks.slack.com/services/YOUR/ORDERS/WEBHOOK"),
				golog.WithSlackUsername("Order Bot"),
				golog.WithSlackEmoji(":shopping_cart:"),
				golog.WithSlackChannel("#orders"),
			),

			// Slack channel for payment notifications
			"slack-payments": golog.NewSlackChannelConfig(
				getEnvOrDefault("SLACK_PAYMENTS_WEBHOOK", "https://hooks.slack.com/services/YOUR/PAYMENTS/WEBHOOK"),
				golog.WithSlackUsername("Payment Bot"),
				golog.WithSlackEmoji(":money_with_wings:"),
				golog.WithSlackChannel("#payments"),
			),

			// Slack channel for error alerts (async for performance)
			"slack-errors": golog.NewSlackChannelConfig(
				getEnvOrDefault("SLACK_ERRORS_WEBHOOK", "https://hooks.slack.com/services/YOUR/ERRORS/WEBHOOK"),
				golog.WithSlackUsername("Error Alert"),
				golog.WithSlackEmoji(":rotating_light:"),
				golog.WithSlackAsync(true),
			),

			// Stack driver: log to both file and errors slack
			"critical": {
				Driver: "stack",
				Level:  "error",
				StackConfig: &golog.StackConfig{
					Channels:         []string{"file", "slack-errors"},
					IgnoreExceptions: true,
				},
			},
		},
	}

	// Initialize
	manager, err := golog.NewManager(config)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	// Share common context across all channels
	manager.ShareContext(map[string]any{
		"app":         "E-Commerce",
		"environment": "production",
		"server":      "web-01",
	})

	// ========================================
	// Example 1: Order Notification (Laravel-style)
	// ========================================
	fmt.Println("📦 Sending order notification...")

	orderLogger, _ := manager.Channel("slack-orders")
	orderLogger.Info("create cart clone - START", map[string]any{
		"cart_id":    32744811,
		"user_id":    795919,
		"user_phone": "551863966",
		"total_100":  2090,
		"total":      20.9,
		"products": []map[string]any{
			{
				"id":       104,
				"title":    "بطاقة 5$ ايتونز - أمريكي",
				"quantity": 1,
			},
		},
	})

	// ========================================
	// Example 2: Payment Notification
	// ========================================
	fmt.Println("💳 Sending payment notification...")

	paymentLogger, _ := manager.Channel("slack-payments")
	paymentLogger.Info("Payment received", map[string]any{
		"order_id":       12345,
		"amount":         99.99,
		"currency":       "SAR",
		"payment_method": "credit_card",
		"card_last_four": "4242",
	})

	// ========================================
	// Example 3: Critical Error (Stack: file + slack)
	// ========================================
	fmt.Println("🚨 Sending critical error notification...")

	criticalLogger, _ := manager.Channel("critical")
	criticalLogger.ErrorWithException(
		"SQLSTATE[HY001]: Memory allocation error: 1038 Out of sort memory",
		fmt.Errorf("database query failed"),
		map[string]any{
			"query":      "SELECT * FROM products WHERE is_active = 1 ORDER BY sort ASC LIMIT 15",
			"connection": "mysql",
			"database":   "app_production",
			"file":       "/app/ProductController.php:142",
		},
	)

	// ========================================
	// Example 4: Using WithContext for user tracking
	// ========================================
	fmt.Println("👤 Logging user actions...")

	fileLogger, _ := manager.Channel("file")
	userLogger := fileLogger.WithContext(map[string]any{
		"user_id":    795919,
		"user_phone": "551863966",
		"session_id": "abc123xyz",
	})

	// All these logs will include the user context
	userLogger.Info("User viewed product catalog")
	userLogger.Info("User added item to cart", map[string]any{
		"product_id": 104,
		"quantity":   2,
	})
	userLogger.Info("User proceeded to checkout")

	// Wait a bit for async messages
	time.Sleep(2 * time.Second)

	fmt.Println("✅ All logs sent!")
	fmt.Println("📁 Check logs/app.log for file logs")
	fmt.Println("💬 Check your Slack channels for notifications")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

