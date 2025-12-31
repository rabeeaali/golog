package golog

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func resetGlobalState() {
	mu.Lock()
	defaultManager = nil
	once = sync.Once{}
	mu.Unlock()
}

func TestInit(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	err := Init(config)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	manager := GetManager()
	if manager == nil {
		t.Error("Expected manager to be set")
	}

	Close()
}

func TestInit_NilConfig(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	err := Init(nil)
	if err != nil {
		t.Fatalf("Init with nil config failed: %v", err)
	}

	manager := GetManager()
	if manager == nil {
		t.Error("Expected manager with default config")
	}

	Close()
}

func TestInit_OnlyOnce(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config1 := &Config{
		Default: "file",
		AppName: "App1",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test1.log")),
		},
	}

	config2 := &Config{
		Default: "file",
		AppName: "App2",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test2.log")),
		},
	}

	Init(config1)
	Init(config2) // Should be ignored

	// Manager should have first config
	manager := GetManager()
	if manager.config.AppName != "App1" {
		t.Error("Second Init should be ignored")
	}

	Close()
}

func TestSetManager(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		AppName: "CustomManager",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	manager, _ := NewManager(config)
	SetManager(manager)

	got := GetManager()
	if got != manager {
		t.Error("Expected custom manager to be set")
	}

	manager.Close()
}

func TestChannel(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file":  NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
			"other": NewFileChannelConfig(filepath.Join(tempDir, "other.log")),
		},
	}

	Init(config)
	defer Close()

	logger, err := Channel("other")
	if err != nil {
		t.Fatalf("Channel failed: %v", err)
	}

	if logger == nil {
		t.Error("Expected logger")
	}
}

func TestChannel_NotInitialized(t *testing.T) {
	resetGlobalState()

	_, err := Channel("file")
	if err != ErrNotInitialized {
		t.Errorf("Expected ErrNotInitialized, got %v", err)
	}
}

func TestDefault(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	Init(config)
	defer Close()

	logger, err := Default()
	if err != nil {
		t.Fatalf("Default failed: %v", err)
	}

	if logger == nil {
		t.Error("Expected logger")
	}
}

func TestDefault_NotInitialized(t *testing.T) {
	resetGlobalState()

	_, err := Default()
	if err != ErrNotInitialized {
		t.Errorf("Expected ErrNotInitialized, got %v", err)
	}
}

func TestClose(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	Init(config)

	err := Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Manager should be nil after close
	if GetManager() != nil {
		t.Error("Manager should be nil after close")
	}
}

func TestClose_NotInitialized(t *testing.T) {
	resetGlobalState()

	// Should not error
	err := Close()
	if err != nil {
		t.Errorf("Close on uninitialized should not error: %v", err)
	}
}

func TestShareContext(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	Init(config)
	defer Close()

	ShareContext(map[string]any{
		"app_version": "1.0.0",
	})

	manager := GetManager()
	ctx := manager.SharedContext()
	if ctx["app_version"] != "1.0.0" {
		t.Error("Expected shared context to be set")
	}
}

func TestShareContext_NotInitialized(t *testing.T) {
	resetGlobalState()

	// Should not panic
	ShareContext(map[string]any{"key": "value"})
}

// Test convenience logging functions
func TestConvenienceLogging(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	Init(config)
	defer Close()

	// Test all convenience functions
	Debug("debug message")
	Info("info message")
	Notice("notice message")
	Warning("warning message")
	Error("error message")
	Critical("critical message")
	Alert("alert message")
	Emergency("emergency message")

	content, _ := os.ReadFile(logPath)
	logContent := string(content)

	levels := []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"}
	for _, level := range levels {
		if !strings.Contains(logContent, level) {
			t.Errorf("Expected %s in log", level)
		}
	}
}

func TestConvenienceLogging_WithContext(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	Init(config)
	defer Close()

	Info("test message", map[string]any{
		"user_id": 123,
		"action":  "login",
	})

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "user_id") {
		t.Error("Expected context in log")
	}
}

func TestErrorWithException_Convenience(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	Init(config)
	defer Close()

	err := &testError{message: "test error"}
	ErrorWithException("error occurred", err, map[string]any{"context": "test"})

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "Exception") {
		t.Error("Expected exception in log")
	}
}

func TestCriticalWithException_Convenience(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	Init(config)
	defer Close()

	err := &testError{message: "critical error"}
	CriticalWithException("critical occurred", err)

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "CRITICAL") {
		t.Error("Expected CRITICAL in log")
	}
}

func TestConvenienceLogging_NotInitialized(t *testing.T) {
	resetGlobalState()

	// These should not panic when not initialized
	Debug("debug")
	Info("info")
	Notice("notice")
	Warning("warning")
	Error("error")
	Critical("critical")
	Alert("alert")
	Emergency("emergency")
	ErrorWithException("error", nil)
	CriticalWithException("critical", nil)
}

// Integration test: Laravel-like usage
func TestLaravelStyleUsage(t *testing.T) {
	resetGlobalState()
	defer resetGlobalState()

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "laravel.log")

	config := &Config{
		Default: "file",
		AppName: "Laravel Log",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	Init(config)
	defer Close()

	// Laravel-style: Log::info('message', ['context' => 'data'])
	Info("create cart clone - START", map[string]any{
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

	content, _ := os.ReadFile(logPath)
	logContent := string(content)

	// Verify Laravel-style format
	if !strings.Contains(logContent, "INFO") {
		t.Error("Expected INFO level")
	}
	if !strings.Contains(logContent, "create cart clone - START") {
		t.Error("Expected message")
	}
	if !strings.Contains(logContent, "cart_id") {
		t.Error("Expected cart_id context")
	}
	if !strings.Contains(logContent, "32744811") {
		t.Error("Expected cart_id value")
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

