package golog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackDriver sends log entries to Slack via webhook
type SlackDriver struct {
	webhookURL string
	username   string
	iconEmoji  string
	iconURL    string
	channel    string
	timeout    time.Duration
	async      bool
	client     *http.Client
}

// SlackMessage represents a Slack message payload
type SlackMessage struct {
	Username    string            `json:"username,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	IconURL     string            `json:"icon_url,omitempty"`
	Channel     string            `json:"channel,omitempty"`
	Text        string            `json:"text,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment represents a Slack message attachment
type SlackAttachment struct {
	Color      string       `json:"color,omitempty"`
	Title      string       `json:"title,omitempty"`
	Text       string       `json:"text,omitempty"`
	Fields     []SlackField `json:"fields,omitempty"`
	Footer     string       `json:"footer,omitempty"`
	FooterIcon string       `json:"footer_icon,omitempty"`
	Timestamp  int64        `json:"ts,omitempty"`
	MarkdownIn []string     `json:"mrkdwn_in,omitempty"`
}

// SlackField represents a field in a Slack attachment
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// NewSlackDriver creates a new Slack driver from configuration
func NewSlackDriver(config ChannelConfig) (Driver, error) {
	if config.SlackConfig == nil {
		return nil, fmt.Errorf("slack configuration is required")
	}

	if config.SlackConfig.WebhookURL == "" {
		return nil, fmt.Errorf("slack webhook URL is required")
	}

	timeout := config.SlackConfig.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	username := config.SlackConfig.Username
	if username == "" {
		username = "GoLog"
	}

	iconEmoji := config.SlackConfig.IconEmoji
	if iconEmoji == "" {
		iconEmoji = ":robot_face:"
	}

	return &SlackDriver{
		webhookURL: config.SlackConfig.WebhookURL,
		username:   username,
		iconEmoji:  iconEmoji,
		iconURL:    config.SlackConfig.IconURL,
		channel:    config.SlackConfig.SlackChannel,
		timeout:    timeout,
		async:      config.SlackConfig.Async,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// Log sends a log entry to Slack
func (d *SlackDriver) Log(entry *Entry) error {
	msg := d.buildMessage(entry)

	if d.async {
		go func() {
			_ = d.send(msg)
		}()
		return nil
	}

	return d.send(msg)
}

// buildMessage builds a Slack message from a log entry (Laravel-style)
func (d *SlackDriver) buildMessage(entry *Entry) *SlackMessage {
	msg := &SlackMessage{
		Username:  d.username,
		IconEmoji: d.iconEmoji,
		Channel:   d.channel,
	}

	if d.iconURL != "" {
		msg.IconURL = d.iconURL
		msg.IconEmoji = ""
	}

	// Build the main attachment
	attachment := SlackAttachment{
		Color:      entry.Level.SlackColor(),
		Title:      fmt.Sprintf("%s %s", entry.Level.Emoji(), entry.Level.String()),
		Timestamp:  entry.Timestamp.Unix(),
		MarkdownIn: []string{"text", "fields"},
	}

	// Add message field
	attachment.Fields = append(attachment.Fields, SlackField{
		Title: "Message",
		Value: entry.Message,
		Short: false,
	})

	// Add level field
	attachment.Fields = append(attachment.Fields, SlackField{
		Title: "Level",
		Value: entry.Level.String(),
		Short: true,
	})

	// Add context fields (like Laravel)
	for key, value := range entry.Context {
		fieldValue := formatSlackValue(value)
		// Determine if the field should be short based on value length
		isShort := len(fieldValue) < 40

		attachment.Fields = append(attachment.Fields, SlackField{
			Title: formatFieldTitle(key),
			Value: fieldValue,
			Short: isShort,
		})
	}

	// Add exception information if present
	if entry.Exception != nil {
		exceptionJSON := entry.ExceptionJSON()
		attachment.Fields = append(attachment.Fields, SlackField{
			Title: "Exception",
			Value: fmt.Sprintf("```%s```", exceptionJSON),
			Short: false,
		})
	}

	// Add footer with channel and timestamp
	channel := entry.Channel
	if channel == "" {
		channel = "default"
	}
	attachment.Footer = fmt.Sprintf("%s | %s", d.username, channel)
	attachment.FooterIcon = "https://avatars.slack-edge.com/2019-01-17/123456789_abc123_48.png"

	msg.Attachments = []SlackAttachment{attachment}
	return msg
}

// formatSlackValue formats a value for Slack display
func formatSlackValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case map[string]any, []any:
		// Pretty print JSON for complex types
		b, err := json.MarshalIndent(val, "", "    ")
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return fmt.Sprintf("```\n%s\n```", string(b))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%v", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		// Try JSON marshal for structs/slices
		if b, err := json.MarshalIndent(val, "", "    "); err == nil {
			return fmt.Sprintf("```\n%s\n```", string(b))
		}
		return fmt.Sprintf("%v", val)
	}
}

// formatFieldTitle converts snake_case or camelCase to Title Case
func formatFieldTitle(s string) string {
	// Convert snake_case to Title_Case for Laravel-like display
	result := ""
	capitalize := true

	for _, c := range s {
		if c == '_' {
			result += "_"
			capitalize = true
			continue
		}

		if capitalize {
			if c >= 'a' && c <= 'z' {
				c = c - 32 // Convert to uppercase
			}
			capitalize = false
		}
		result += string(c)
	}

	return result
}

// send sends a message to Slack
func (d *SlackDriver) send(msg *SlackMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequest("POST", d.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create slack request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

// Close closes the driver
func (d *SlackDriver) Close() error {
	return nil
}

// Name returns the driver name
func (d *SlackDriver) Name() string {
	return "slack"
}

