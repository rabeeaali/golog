# GoLog 📝

[![Go Reference](https://pkg.go.dev/badge/github.com/golog-pkg/golog.svg)](https://pkg.go.dev/github.com/golog-pkg/golog)
[![Go Report Card](https://goreportcard.com/badge/github.com/golog-pkg/golog)](https://goreportcard.com/report/github.com/golog-pkg/golog)
[![CI](https://github.com/golog-pkg/golog/actions/workflows/ci.yml/badge.svg)](https://github.com/golog-pkg/golog/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/golog-pkg/golog/branch/main/graph/badge.svg)](https://codecov.io/gh/golog-pkg/golog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A **Laravel-inspired** logging package for Go with support for multiple drivers and channels.

## ✨ Features

- 🎯 **Laravel-style API** - Familiar logging patterns for Laravel developers
- 📁 **File Driver** - Write logs to files with Laravel-style formatting
- 💬 **Slack Driver** - Send beautiful formatted logs to Slack webhooks
- 🔀 **Multiple Channels** - Configure different channels for different purposes
- 📚 **Stack Driver** - Log to multiple channels simultaneously
- 🏷️ **Context Support** - Add structured context data to your logs
- ⚡ **Async Support** - Send Slack messages asynchronously
- 🎨 **Beautiful Slack Messages** - Laravel-style formatted Slack attachments with colors and emojis

## 📦 Installation

```bash
go get github.com/golog-pkg/golog
```

## 🚀 Quick Start

```go
package main

import "github.com/golog-pkg/golog"

func main() {
    // Initialize with default configuration (logs to file)
    golog.Init(nil)
    defer golog.Close()

    // Log messages - just like Laravel's Log facade!
    golog.Info("User logged in", map[string]any{
        "user_id": 123,
        "ip":      "192.168.1.1",
    })

    golog.Error("Something went wrong", map[string]any{
        "error": "connection timeout",
    })
}
```

## 📖 Documentation

### Configuration

```go
config := &golog.Config{
    Default: "file",
    AppName: "MyApp",
    Channels: map[string]golog.ChannelConfig{
        // File logging
        "file": golog.NewFileChannelConfig("logs/app.log"),

        // Slack alerts
        "slack": golog.NewSlackChannelConfig(
            "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
            golog.WithSlackUsername("MyApp Alerts"),
            golog.WithSlackEmoji(":warning:"),
        ),
    },
}

golog.Init(config)
```

### Multiple Slack Channels

Perfect for sending different types of logs to different Slack channels:

```go
config := &golog.Config{
    Default: "file",
    Channels: map[string]golog.ChannelConfig{
        "file": golog.NewFileChannelConfig("logs/app.log"),

        // Orders channel
        "slack-orders": golog.NewSlackChannelConfig(
            os.Getenv("SLACK_ORDERS_WEBHOOK"),
            golog.WithSlackUsername("Order Bot"),
            golog.WithSlackEmoji(":shopping_cart:"),
            golog.WithSlackChannel("#orders"),
        ),

        // Payments channel
        "slack-payments": golog.NewSlackChannelConfig(
            os.Getenv("SLACK_PAYMENTS_WEBHOOK"),
            golog.WithSlackUsername("Payment Bot"),
            golog.WithSlackEmoji(":money_with_wings:"),
            golog.WithSlackChannel("#payments"),
        ),

        // Critical errors
        "slack-errors": golog.NewSlackChannelConfig(
            os.Getenv("SLACK_ERRORS_WEBHOOK"),
            golog.WithSlackUsername("Error Alert"),
            golog.WithSlackEmoji(":rotating_light:"),
        ),
    },
}
```

### Logging with Context (Laravel-style)

```go
// Like Laravel's Log::info('message', ['key' => 'value'])
golog.Info("create cart clone - START", map[string]any{
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
```

### Logging to Specific Channel

```go
// Like Laravel's Log::channel('slack')->error()
slackLog, _ := golog.Channel("slack-orders")
slackLog.Info("New order received", map[string]any{
    "order_id": 12345,
    "amount":   99.99,
})
```

### Logging Errors with Exception Details

```go
err := database.Query(...)
if err != nil {
    golog.ErrorWithException(
        "SQLSTATE[HY001]: Memory allocation error",
        err,
        map[string]any{
            "query": "SELECT * FROM products",
            "file":  "/app/db.go:42",
        },
    )
}
```

### Using WithContext (like Laravel)

```go
// Create a logger with persistent context
logger, _ := golog.Channel("file")
userLogger := logger.WithContext(map[string]any{
    "user_id":   123,
    "user_name": "Ahmed",
    "session":   "abc123",
})

// All subsequent logs include the user context
userLogger.Info("Viewed products")
userLogger.Info("Added to cart", map[string]any{"product_id": 456})
userLogger.Info("Completed checkout")
```

### Sharing Context Across All Channels

```go
// Like Laravel's Log::shareContext()
golog.ShareContext(map[string]any{
    "app_version": "1.0.0",
    "environment": "production",
    "server_id":   "web-01",
})

// Now ALL logs will include this context
golog.Info("Application started")
```

### Stack Driver (Multiple Outputs)

Log to multiple channels at once:

```go
config := &golog.Config{
    Default: "stack",
    Channels: map[string]golog.ChannelConfig{
        "file":  golog.NewFileChannelConfig("logs/app.log"),
        "slack": golog.NewSlackChannelConfig("https://..."),

        // Log to both file AND slack
        "stack": {
            Driver: "stack",
            StackConfig: &golog.StackConfig{
                Channels:         []string{"file", "slack"},
                IgnoreExceptions: true,
            },
        },
    },
}
```

## 📊 Log Levels

| Level     | Method        | Description                                 |
| --------- | ------------- | ------------------------------------------- |
| DEBUG     | `Debug()`     | Detailed debug information                  |
| INFO      | `Info()`      | Interesting events (user logins, SQL logs)  |
| NOTICE    | `Notice()`    | Normal but significant events               |
| WARNING   | `Warning()`   | Exceptional occurrences that are not errors |
| ERROR     | `Error()`     | Runtime errors                              |
| CRITICAL  | `Critical()`  | Critical conditions                         |
| ALERT     | `Alert()`     | Action must be taken immediately            |
| EMERGENCY | `Emergency()` | System is unusable                          |

## 📤 Output Examples

### File Output

```
[2024-01-15 10:30:45] local.INFO: create cart clone - START
  Cart_id: 32744811
  User_id: 795919
  User_phone: 551863966
  Total_100: 2090
  Total: 20.9
```

### Slack Output

The Slack driver formats messages beautifully with:

- Color-coded attachments based on log level
- Emoji indicators
- Structured fields for context data
- Timestamp and channel info in footer

![Slack Log Example](https://via.placeholder.com/400x200?text=Slack+Log+Preview)

## ⚙️ Configuration Options

### File Driver

```go
golog.NewFileChannelConfig("logs/app.log",
    golog.WithFileMaxSize(100),                    // Max size in MB
    golog.WithFileDateFormat("2006-01-02 15:04:05"),
)
```

### Slack Driver

```go
golog.NewSlackChannelConfig(webhookURL,
    golog.WithSlackUsername("Bot Name"),
    golog.WithSlackEmoji(":robot_face:"),
    golog.WithSlackChannel("#channel-name"),
    golog.WithSlackAsync(true),  // Send asynchronously
)
```

## 🔧 Custom Drivers

Register your own custom driver:

```go
golog.RegisterDriver("custom", func(config golog.ChannelConfig) (golog.Driver, error) {
    return &MyCustomDriver{}, nil
})
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by [Laravel's](https://laravel.com/) excellent logging system.
