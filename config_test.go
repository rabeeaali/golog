package golog

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig should not return nil")
	}

	if config.Default != "file" {
		t.Errorf("Expected default channel 'file', got %q", config.Default)
	}

	if config.AppName != "GoLog" {
		t.Errorf("Expected app name 'GoLog', got %q", config.AppName)
	}

	if len(config.Channels) == 0 {
		t.Error("Expected at least one default channel")
	}

	fileChannel, ok := config.Channels["file"]
	if !ok {
		t.Fatal("Expected 'file' channel in default config")
	}

	if fileChannel.Driver != "file" {
		t.Errorf("Expected file driver, got %q", fileChannel.Driver)
	}

	if fileChannel.FileConfig == nil {
		t.Error("Expected FileConfig to be set")
	}

	if fileChannel.FileConfig.Path != "logs/app.log" {
		t.Errorf("Expected path 'logs/app.log', got %q", fileChannel.FileConfig.Path)
	}
}

func TestNewSlackChannelConfig(t *testing.T) {
	webhookURL := "https://hooks.slack.com/test"

	config := NewSlackChannelConfig(webhookURL)

	if config.Driver != "slack" {
		t.Errorf("Expected driver 'slack', got %q", config.Driver)
	}

	if config.Level != "error" {
		t.Errorf("Expected default level 'error', got %q", config.Level)
	}

	if config.SlackConfig == nil {
		t.Fatal("Expected SlackConfig to be set")
	}

	if config.SlackConfig.WebhookURL != webhookURL {
		t.Errorf("Expected webhook URL %q, got %q", webhookURL, config.SlackConfig.WebhookURL)
	}

	if config.SlackConfig.Username != "GoLog" {
		t.Errorf("Expected default username 'GoLog', got %q", config.SlackConfig.Username)
	}

	if config.SlackConfig.IconEmoji != ":robot_face:" {
		t.Errorf("Expected default emoji ':robot_face:', got %q", config.SlackConfig.IconEmoji)
	}

	if config.SlackConfig.Timeout != 10*time.Second {
		t.Errorf("Expected default timeout 10s, got %v", config.SlackConfig.Timeout)
	}
}

func TestNewSlackChannelConfig_WithOptions(t *testing.T) {
	webhookURL := "https://hooks.slack.com/test"

	config := NewSlackChannelConfig(
		webhookURL,
		WithSlackUsername("Custom Bot"),
		WithSlackEmoji(":fire:"),
		WithSlackChannel("#alerts"),
		WithSlackAsync(true),
	)

	if config.SlackConfig.Username != "Custom Bot" {
		t.Errorf("Expected username 'Custom Bot', got %q", config.SlackConfig.Username)
	}

	if config.SlackConfig.IconEmoji != ":fire:" {
		t.Errorf("Expected emoji ':fire:', got %q", config.SlackConfig.IconEmoji)
	}

	if config.SlackConfig.SlackChannel != "#alerts" {
		t.Errorf("Expected channel '#alerts', got %q", config.SlackConfig.SlackChannel)
	}

	if !config.SlackConfig.Async {
		t.Error("Expected async to be true")
	}
}

func TestNewFileChannelConfig(t *testing.T) {
	path := "logs/custom.log"

	config := NewFileChannelConfig(path)

	if config.Driver != "file" {
		t.Errorf("Expected driver 'file', got %q", config.Driver)
	}

	if config.Level != "debug" {
		t.Errorf("Expected default level 'debug', got %q", config.Level)
	}

	if config.FileConfig == nil {
		t.Fatal("Expected FileConfig to be set")
	}

	if config.FileConfig.Path != path {
		t.Errorf("Expected path %q, got %q", path, config.FileConfig.Path)
	}

	if config.FileConfig.MaxSize != 100 {
		t.Errorf("Expected default max size 100, got %d", config.FileConfig.MaxSize)
	}

	if config.FileConfig.MaxBackups != 3 {
		t.Errorf("Expected default max backups 3, got %d", config.FileConfig.MaxBackups)
	}

	if !config.FileConfig.Compress {
		t.Error("Expected compress to be true by default")
	}
}

func TestNewFileChannelConfig_WithOptions(t *testing.T) {
	path := "logs/custom.log"

	config := NewFileChannelConfig(
		path,
		WithFileMaxSize(50),
		WithFileDateFormat("2006-01-02"),
	)

	if config.FileConfig.MaxSize != 50 {
		t.Errorf("Expected max size 50, got %d", config.FileConfig.MaxSize)
	}

	if config.FileConfig.DateFormat != "2006-01-02" {
		t.Errorf("Expected date format '2006-01-02', got %q", config.FileConfig.DateFormat)
	}
}

func TestChannelConfig_StackConfig(t *testing.T) {
	config := ChannelConfig{
		Driver: "stack",
		Level:  "debug",
		StackConfig: &StackConfig{
			Channels:         []string{"file", "slack"},
			IgnoreExceptions: true,
		},
	}

	if config.Driver != "stack" {
		t.Errorf("Expected driver 'stack', got %q", config.Driver)
	}

	if config.StackConfig == nil {
		t.Fatal("Expected StackConfig to be set")
	}

	if len(config.StackConfig.Channels) != 2 {
		t.Errorf("Expected 2 channels, got %d", len(config.StackConfig.Channels))
	}

	if !config.StackConfig.IgnoreExceptions {
		t.Error("Expected IgnoreExceptions to be true")
	}
}

