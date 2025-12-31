package golog

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSlackDriver(t *testing.T) {
	config := ChannelConfig{
		Driver: "slack",
		Level:  "error",
		SlackConfig: &SlackConfig{
			WebhookURL: "https://hooks.slack.com/test",
			Username:   "TestBot",
			IconEmoji:  ":robot:",
			Timeout:    5 * time.Second,
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	if driver.Name() != "slack" {
		t.Errorf("Expected driver name 'slack', got %q", driver.Name())
	}
}

func TestNewSlackDriver_NoConfig(t *testing.T) {
	config := ChannelConfig{
		Driver: "slack",
	}

	_, err := NewSlackDriver(config)
	if err == nil {
		t.Error("Expected error for missing SlackConfig")
	}
}

func TestNewSlackDriver_NoWebhookURL(t *testing.T) {
	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			Username: "TestBot",
		},
	}

	_, err := NewSlackDriver(config)
	if err == nil {
		t.Error("Expected error for missing webhook URL")
	}
}

func TestNewSlackDriver_Defaults(t *testing.T) {
	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: "https://hooks.slack.com/test",
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	sd := driver.(*SlackDriver)

	if sd.username != "GoLog" {
		t.Errorf("Expected default username 'GoLog', got %q", sd.username)
	}

	if sd.iconEmoji != ":robot_face:" {
		t.Errorf("Expected default emoji ':robot_face:', got %q", sd.iconEmoji)
	}

	if sd.timeout != 10*time.Second {
		t.Errorf("Expected default timeout 10s, got %v", sd.timeout)
	}
}

func TestSlackDriver_Log(t *testing.T) {
	var receivedPayload []byte

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPayload, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL:   server.URL,
			Username:     "TestBot",
			IconEmoji:    ":test:",
			SlackChannel: "#test",
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	entry := NewEntry(InfoLevel, "test message")
	entry.SetChannel("slack")
	entry.WithContext(map[string]any{
		"user_id": 123,
		"action":  "login",
	})

	if err := driver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Parse the payload
	var msg SlackMessage
	if err := json.Unmarshal(receivedPayload, &msg); err != nil {
		t.Fatalf("Failed to parse payload: %v", err)
	}

	// Verify payload
	if msg.Username != "TestBot" {
		t.Errorf("Expected username 'TestBot', got %q", msg.Username)
	}

	if msg.IconEmoji != ":test:" {
		t.Errorf("Expected emoji ':test:', got %q", msg.IconEmoji)
	}

	if msg.Channel != "#test" {
		t.Errorf("Expected channel '#test', got %q", msg.Channel)
	}

	if len(msg.Attachments) == 0 {
		t.Fatal("Expected at least one attachment")
	}

	// Check attachment
	attachment := msg.Attachments[0]
	if attachment.Color != InfoLevel.SlackColor() {
		t.Errorf("Expected color %q, got %q", InfoLevel.SlackColor(), attachment.Color)
	}

	// Check fields contain message
	hasMessage := false
	for _, field := range attachment.Fields {
		if field.Title == "Message" && field.Value == "test message" {
			hasMessage = true
			break
		}
	}
	if !hasMessage {
		t.Error("Expected Message field in attachment")
	}
}

func TestSlackDriver_LogWithException(t *testing.T) {
	var receivedPayload []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPayload, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: server.URL,
			Username:   "ErrorBot",
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	entry := NewEntry(ErrorLevel, "database error")
	entry.WithException("DatabaseError", "connection timeout", 500, "/app/db.go", 42, []string{
		"/app/main.go:10",
	})

	if err := driver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	var msg SlackMessage
	if err := json.Unmarshal(receivedPayload, &msg); err != nil {
		t.Fatalf("Failed to parse payload: %v", err)
	}

	// Check for Exception field
	hasException := false
	for _, field := range msg.Attachments[0].Fields {
		if field.Title == "Exception" {
			hasException = true
			break
		}
	}
	if !hasException {
		t.Error("Expected Exception field in attachment")
	}
}

func TestSlackDriver_LogAsync(t *testing.T) {
	received := make(chan bool, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: server.URL,
			Async:      true,
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	entry := NewEntry(InfoLevel, "async message")

	// Log should return immediately
	if err := driver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Wait for async send
	select {
	case <-received:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("Async log was not sent within timeout")
	}
}

func TestSlackDriver_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: server.URL,
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	entry := NewEntry(ErrorLevel, "test")

	err = driver.Log(entry)
	if err == nil {
		t.Error("Expected error for non-OK response")
	}
}

func TestSlackDriver_Close(t *testing.T) {
	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: "https://hooks.slack.com/test",
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	// Close should not error
	if err := driver.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestSlackDriver_IconURL(t *testing.T) {
	var receivedPayload []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPayload, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: server.URL,
			IconURL:    "https://example.com/icon.png",
			IconEmoji:  ":ignored:", // Should be ignored when IconURL is set
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	entry := NewEntry(InfoLevel, "test")
	driver.Log(entry)

	var msg SlackMessage
	json.Unmarshal(receivedPayload, &msg)

	if msg.IconURL != "https://example.com/icon.png" {
		t.Errorf("Expected IconURL to be set, got %q", msg.IconURL)
	}

	if msg.IconEmoji != "" {
		t.Error("IconEmoji should be empty when IconURL is set")
	}
}

func TestFormatFieldTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_id", "User_Id"},
		{"cart_id", "Cart_Id"},
		{"total_100", "Total_100"},
		{"simple", "Simple"},
		{"already_Title", "Already_Title"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := formatFieldTitle(tt.input)
			if got != tt.expected {
				t.Errorf("formatFieldTitle(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatSlackValue(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		contains string
	}{
		{"string", "hello", "hello"},
		{"int", 123, "123"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"slice", []int{1, 2, 3}, "["},
		{"map", map[string]int{"a": 1}, "{"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSlackValue(tt.input)
			if got == "" {
				t.Error("formatSlackValue should not return empty string")
			}
		})
	}
}

func TestSlackDriver_ComplexContext(t *testing.T) {
	var receivedPayload []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPayload, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := ChannelConfig{
		Driver: "slack",
		SlackConfig: &SlackConfig{
			WebhookURL: server.URL,
		},
	}

	driver, err := NewSlackDriver(config)
	if err != nil {
		t.Fatalf("NewSlackDriver failed: %v", err)
	}

	// Complex context like Laravel cart example
	entry := NewEntry(InfoLevel, "create cart clone - START")
	entry.WithContext(map[string]any{
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

	if err := driver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Just verify it was sent successfully and payload is valid JSON
	var msg SlackMessage
	if err := json.Unmarshal(receivedPayload, &msg); err != nil {
		t.Fatalf("Failed to parse payload: %v", err)
	}

	if len(msg.Attachments) == 0 {
		t.Error("Expected attachments")
	}

	// Should have multiple fields
	if len(msg.Attachments[0].Fields) < 3 {
		t.Error("Expected multiple fields for context")
	}
}
