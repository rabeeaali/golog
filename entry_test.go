package golog

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestNewEntry(t *testing.T) {
	entry := NewEntry(InfoLevel, "test message")

	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got %q", entry.Message)
	}

	if entry.Level != InfoLevel {
		t.Errorf("Expected level INFO, got %v", entry.Level)
	}

	if entry.Context == nil {
		t.Error("Expected context to be initialized")
	}

	if entry.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	// Timestamp should be recent
	if time.Since(entry.Timestamp) > time.Second {
		t.Error("Timestamp should be recent")
	}
}

func TestEntry_WithContext(t *testing.T) {
	entry := NewEntry(InfoLevel, "test")

	result := entry.WithContext(map[string]any{
		"user_id": 123,
		"action":  "login",
	})

	// Should return same entry (fluent API)
	if result != entry {
		t.Error("WithContext should return same entry")
	}

	if entry.Context["user_id"] != 123 {
		t.Errorf("Expected user_id=123, got %v", entry.Context["user_id"])
	}

	if entry.Context["action"] != "login" {
		t.Errorf("Expected action=login, got %v", entry.Context["action"])
	}
}

func TestEntry_With(t *testing.T) {
	entry := NewEntry(InfoLevel, "test")

	result := entry.With("key", "value")

	if result != entry {
		t.Error("With should return same entry")
	}

	if entry.Context["key"] != "value" {
		t.Errorf("Expected key=value, got %v", entry.Context["key"])
	}
}

func TestEntry_WithError(t *testing.T) {
	entry := NewEntry(ErrorLevel, "something failed")
	err := errors.New("database connection failed")

	result := entry.WithError(err)

	if result != entry {
		t.Error("WithError should return same entry")
	}

	if entry.Exception == nil {
		t.Fatal("Expected exception to be set")
	}

	if entry.Exception.Message != "database connection failed" {
		t.Errorf("Expected error message, got %q", entry.Exception.Message)
	}

	if entry.Exception.Class == "" {
		t.Error("Expected error class to be set")
	}
}

func TestEntry_WithError_Nil(t *testing.T) {
	entry := NewEntry(ErrorLevel, "test")

	result := entry.WithError(nil)

	if result != entry {
		t.Error("WithError should return same entry")
	}

	if entry.Exception != nil {
		t.Error("Exception should be nil for nil error")
	}
}

func TestEntry_WithException(t *testing.T) {
	entry := NewEntry(ErrorLevel, "test")

	trace := []string{
		"/app/main.go:10",
		"/app/handler.go:25",
	}

	entry.WithException("DatabaseError", "connection timeout", 500, "/app/db.go", 42, trace)

	if entry.Exception == nil {
		t.Fatal("Expected exception to be set")
	}

	if entry.Exception.Class != "DatabaseError" {
		t.Errorf("Expected class DatabaseError, got %q", entry.Exception.Class)
	}

	if entry.Exception.Message != "connection timeout" {
		t.Errorf("Expected message 'connection timeout', got %q", entry.Exception.Message)
	}

	if entry.Exception.Code != 500 {
		t.Errorf("Expected code 500, got %d", entry.Exception.Code)
	}

	if entry.Exception.File != "/app/db.go" {
		t.Errorf("Expected file /app/db.go, got %q", entry.Exception.File)
	}

	if entry.Exception.Line != 42 {
		t.Errorf("Expected line 42, got %d", entry.Exception.Line)
	}

	if len(entry.Exception.Trace) != 2 {
		t.Errorf("Expected 2 trace entries, got %d", len(entry.Exception.Trace))
	}
}

func TestEntry_SetChannel(t *testing.T) {
	entry := NewEntry(InfoLevel, "test")

	result := entry.SetChannel("slack")

	if result != entry {
		t.Error("SetChannel should return same entry")
	}

	if entry.Channel != "slack" {
		t.Errorf("Expected channel 'slack', got %q", entry.Channel)
	}
}

func TestEntry_ToJSON(t *testing.T) {
	entry := NewEntry(InfoLevel, "test message")
	entry.With("user_id", 123)
	entry.SetChannel("file")

	jsonBytes, err := entry.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]any
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	if parsed["message"] != "test message" {
		t.Errorf("Expected message in JSON, got %v", parsed["message"])
	}
}

func TestEntry_ContextJSON(t *testing.T) {
	entry := NewEntry(InfoLevel, "test")
	entry.WithContext(map[string]any{
		"user_id": 123,
		"action":  "test",
	})

	jsonStr := entry.ContextJSON()

	if jsonStr == "" {
		t.Error("ContextJSON should not be empty")
	}

	// Should be valid JSON
	var parsed map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
}

func TestEntry_ContextJSON_Empty(t *testing.T) {
	entry := NewEntry(InfoLevel, "test")

	jsonStr := entry.ContextJSON()

	if jsonStr != "" {
		t.Errorf("ContextJSON should be empty for no context, got %q", jsonStr)
	}
}

func TestEntry_ExceptionJSON(t *testing.T) {
	entry := NewEntry(ErrorLevel, "test")
	entry.WithError(errors.New("test error"))

	jsonStr := entry.ExceptionJSON()

	if jsonStr == "" {
		t.Error("ExceptionJSON should not be empty")
	}

	// Should be valid JSON
	var parsed map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
}

func TestEntry_ExceptionJSON_Nil(t *testing.T) {
	entry := NewEntry(InfoLevel, "test")

	jsonStr := entry.ExceptionJSON()

	if jsonStr != "" {
		t.Errorf("ExceptionJSON should be empty for nil exception, got %q", jsonStr)
	}
}

// Custom error type for testing
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func TestEntry_WithError_CustomType(t *testing.T) {
	entry := NewEntry(ErrorLevel, "test")
	err := &CustomError{Code: 404, Message: "not found"}

	entry.WithError(err)

	if entry.Exception == nil {
		t.Fatal("Expected exception to be set")
	}

	// Should capture the type name
	if entry.Exception.Class == "" {
		t.Error("Expected class to be set for custom error")
	}
}

