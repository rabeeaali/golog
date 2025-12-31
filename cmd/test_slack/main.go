package main

import (
	"fmt"
	"time"

	"github.com/rabeeaali/golog"
)

func main() {
	fmt.Println("🚀 Testing GoLog Slack Integration...")

	// Configuration with your Slack webhook
	config := &golog.Config{
		Default: "slack",
		AppName: "GoLog Test",
		Channels: map[string]golog.ChannelConfig{
			// Your Slack channel
			"slack": golog.NewSlackChannelConfig(
				"https://hooks.slack.com/services/T04PVTV1GE5/B0A72TBG47J/XXXXXXXXXX",
				golog.WithSlackUsername("Laravel Log"),
				golog.WithSlackEmoji(":rocket:"),
			),
		},
	}

	// Initialize
	manager, err := golog.NewManager(config)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	logger, err := manager.Channel("slack")
	if err != nil {
		panic(err)
	}

	// =============================================
	// Test 1: INFO log (like Laravel example)
	// =============================================
	fmt.Println("📤 Sending INFO log...")

	logger.Info("create cart clone - START", map[string]any{
		"cart_id":    32744811,
		"user_id":    795919,
		"user_phone": "551863966",
		"total_100":  2090,
		"total":      20.9,
		"products": []map[string]any{
			{
				"id":       104,
				"title":    "Product Title",
				"quantity": 1,
			},
		},
	})

	time.Sleep(1 * time.Second)

	// =============================================
	// Test 2: ERROR log with exception
	// =============================================
	fmt.Println("📤 Sending ERROR log with exception...")

	logger.ErrorWithException(
		"SQLSTATE[HY001]: Memory allocation error: 1038 Out of sort memory, consider increasing server sort buffer size",
		fmt.Errorf("database query failed"),
		map[string]any{
			"query":      "select * from `products` where `is_active` = 1 and `products`.`deleted_at` is null order by `sort` asc limit 15 offset 0",
			"connection": "mysql",
			"file":       "/home/forge/app.example.com/vendor/laravel/framework/src/Illuminate/Database/Connection.php:712",
		},
	)

	time.Sleep(1 * time.Second)

	// =============================================
	// Test 3: WARNING log
	// =============================================
	fmt.Println("📤 Sending WARNING log...")

	logger.Warning("API Rate limit approaching", map[string]any{
		"current_rate": 950,
		"max_rate":     1000,
		"endpoint":     "/api/v1/products",
		"user_id":      795919,
	})

	time.Sleep(1 * time.Second)

	// =============================================
	// Test 4: CRITICAL log
	// =============================================
	fmt.Println("📤 Sending CRITICAL log...")

	logger.Critical("Payment gateway connection failed!", map[string]any{
		"gateway":     "Stripe",
		"error":       "Connection timeout after 30s",
		"retry_count": 3,
		"order_id":    12345,
		"amount":      99.99,
	})

	fmt.Println("")
	fmt.Println("✅ All logs sent! Check your Slack channel.")
}
