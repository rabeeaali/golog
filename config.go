package golog

import "time"

// Config is the main configuration for the log manager
type Config struct {
	// Default is the default channel name to use
	Default string `json:"default" yaml:"default"`

	// Channels is a map of channel configurations
	Channels map[string]ChannelConfig `json:"channels" yaml:"channels"`

	// AppName is the application name (used in Slack messages)
	AppName string `json:"app_name" yaml:"app_name"`
}

// ChannelConfig represents configuration for a single logging channel
type ChannelConfig struct {
	// Driver is the type of driver: "file", "slack", "stack"
	Driver string `json:"driver" yaml:"driver"`

	// Level is the minimum log level for this channel
	Level string `json:"level" yaml:"level"`

	// FileConfig contains file-specific configuration
	*FileConfig `json:",inline" yaml:",inline"`

	// SlackConfig contains Slack-specific configuration
	*SlackConfig `json:",inline" yaml:",inline"`

	// StackConfig contains stack-specific configuration (for combining channels)
	*StackConfig `json:",inline" yaml:",inline"`
}

// FileConfig contains configuration for the file driver
type FileConfig struct {
	// Path is the file path for logs
	Path string `json:"path" yaml:"path"`

	// MaxSize is the maximum size in MB before rotation (0 = no rotation)
	MaxSize int `json:"max_size" yaml:"max_size"`

	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `json:"max_backups" yaml:"max_backups"`

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `json:"max_age" yaml:"max_age"`

	// Compress determines if the rotated log files should be compressed
	Compress bool `json:"compress" yaml:"compress"`

	// Permission is the file permission mode
	Permission string `json:"permission" yaml:"permission"`

	// DateFormat is the date format for log entries
	DateFormat string `json:"date_format" yaml:"date_format"`
}

// SlackConfig contains configuration for the Slack driver
type SlackConfig struct {
	// WebhookURL is the Slack webhook URL
	WebhookURL string `json:"webhook_url" yaml:"webhook_url"`

	// Username is the bot username shown in Slack
	Username string `json:"username" yaml:"username"`

	// IconEmoji is the emoji icon for the bot (e.g., ":boom:")
	IconEmoji string `json:"icon_emoji" yaml:"icon_emoji"`

	// IconURL is the URL of the icon to use (alternative to IconEmoji)
	IconURL string `json:"icon_url" yaml:"icon_url"`

	// SlackChannel is the Slack channel to post to (overrides webhook default)
	SlackChannel string `json:"slack_channel" yaml:"slack_channel"`

	// Timeout is the HTTP timeout for sending to Slack
	Timeout time.Duration `json:"timeout" yaml:"timeout"`

	// Async determines if messages should be sent asynchronously
	Async bool `json:"async" yaml:"async"`
}

// StackConfig contains configuration for the stack driver (multiple channels)
type StackConfig struct {
	// Channels is a list of channel names to log to
	Channels []string `json:"channels" yaml:"channels"`

	// IgnoreExceptions determines if exceptions from individual channels should be ignored
	IgnoreExceptions bool `json:"ignore_exceptions" yaml:"ignore_exceptions"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Default: "file",
		AppName: "GoLog",
		Channels: map[string]ChannelConfig{
			"file": {
				Driver: "file",
				Level:  "debug",
				FileConfig: &FileConfig{
					Path:       "logs/app.log",
					MaxSize:    100,
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
					DateFormat: "2006-01-02 15:04:05",
				},
			},
		},
	}
}

// NewSlackChannelConfig creates a new Slack channel configuration
func NewSlackChannelConfig(webhookURL string, options ...SlackOption) ChannelConfig {
	cfg := ChannelConfig{
		Driver: "slack",
		Level:  "error",
		SlackConfig: &SlackConfig{
			WebhookURL: webhookURL,
			Username:   "GoLog",
			IconEmoji:  ":robot_face:",
			Timeout:    10 * time.Second,
		},
	}

	for _, opt := range options {
		opt(cfg.SlackConfig)
	}

	return cfg
}

// SlackOption is a function that configures a SlackConfig
type SlackOption func(*SlackConfig)

// WithSlackUsername sets the Slack username
func WithSlackUsername(username string) SlackOption {
	return func(c *SlackConfig) {
		c.Username = username
	}
}

// WithSlackEmoji sets the Slack icon emoji
func WithSlackEmoji(emoji string) SlackOption {
	return func(c *SlackConfig) {
		c.IconEmoji = emoji
	}
}

// WithSlackChannel sets the Slack channel
func WithSlackChannel(channel string) SlackOption {
	return func(c *SlackConfig) {
		c.SlackChannel = channel
	}
}

// WithSlackAsync enables async sending
func WithSlackAsync(async bool) SlackOption {
	return func(c *SlackConfig) {
		c.Async = async
	}
}

// NewFileChannelConfig creates a new file channel configuration
func NewFileChannelConfig(path string, options ...FileOption) ChannelConfig {
	cfg := ChannelConfig{
		Driver: "file",
		Level:  "debug",
		FileConfig: &FileConfig{
			Path:       path,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
			DateFormat: "2006-01-02 15:04:05",
		},
	}

	for _, opt := range options {
		opt(cfg.FileConfig)
	}

	return cfg
}

// FileOption is a function that configures a FileConfig
type FileOption func(*FileConfig)

// WithFileMaxSize sets the max file size in MB
func WithFileMaxSize(size int) FileOption {
	return func(c *FileConfig) {
		c.MaxSize = size
	}
}

// WithFileDateFormat sets the date format
func WithFileDateFormat(format string) FileOption {
	return func(c *FileConfig) {
		c.DateFormat = format
	}
}
