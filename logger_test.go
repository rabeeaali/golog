package golog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestLogger(t *testing.T) (*Logger, string) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	manager, err := NewManager(config)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	t.Cleanup(func() {
		manager.Close()
	})

	logger, err := manager.Channel("file")
	if err != nil {
		t.Fatalf("Failed to get channel: %v", err)
	}

	return logger, logPath
}

func TestLogger_Debug(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Debug("debug message", map[string]any{"key": "value"})

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "DEBUG") {
		t.Error("Expected DEBUG level in log")
	}
	if !strings.Contains(string(content), "debug message") {
		t.Error("Expected message in log")
	}
}

func TestLogger_Info(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Info("info message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "INFO") {
		t.Error("Expected INFO level in log")
	}
}

func TestLogger_Notice(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Notice("notice message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "NOTICE") {
		t.Error("Expected NOTICE level in log")
	}
}

func TestLogger_Warning(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Warning("warning message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "WARNING") {
		t.Error("Expected WARNING level in log")
	}
}

func TestLogger_Error(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Error("error message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "ERROR") {
		t.Error("Expected ERROR level in log")
	}
}

func TestLogger_Critical(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Critical("critical message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "CRITICAL") {
		t.Error("Expected CRITICAL level in log")
	}
}

func TestLogger_Alert(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Alert("alert message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "ALERT") {
		t.Error("Expected ALERT level in log")
	}
}

func TestLogger_Emergency(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Emergency("emergency message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "EMERGENCY") {
		t.Error("Expected EMERGENCY level in log")
	}
}

func TestLogger_Log(t *testing.T) {
	logger, logPath := createTestLogger(t)

	logger.Log(WarningLevel, "custom level message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "WARNING") {
		t.Error("Expected WARNING level in log")
	}
}

func TestLogger_ErrorWithException(t *testing.T) {
	logger, logPath := createTestLogger(t)

	err := &testError{message: "test error"}
	logger.ErrorWithException("error occurred", err, map[string]any{"context": "test"})

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "ERROR") {
		t.Error("Expected ERROR level in log")
	}
	if !strings.Contains(string(content), "Exception") {
		t.Error("Expected exception in log")
	}
}

func TestLogger_CriticalWithException(t *testing.T) {
	logger, logPath := createTestLogger(t)

	err := &testError{message: "critical error"}
	logger.CriticalWithException("critical error occurred", err)

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "CRITICAL") {
		t.Error("Expected CRITICAL level in log")
	}
}

func TestLogger_AlertWithException(t *testing.T) {
	logger, logPath := createTestLogger(t)

	err := &testError{message: "alert error"}
	logger.AlertWithException("alert occurred", err)

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "ALERT") {
		t.Error("Expected ALERT level in log")
	}
}

func TestLogger_EmergencyWithException(t *testing.T) {
	logger, logPath := createTestLogger(t)

	err := &testError{message: "emergency error"}
	logger.EmergencyWithException("emergency occurred", err)

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "EMERGENCY") {
		t.Error("Expected EMERGENCY level in log")
	}
}

func TestLogger_WithContext(t *testing.T) {
	logger, logPath := createTestLogger(t)

	contextLogger := logger.WithContext(map[string]any{
		"user_id":   123,
		"user_name": "test",
	})

	// Should return new logger
	if contextLogger == logger {
		t.Error("WithContext should return new logger")
	}

	contextLogger.Info("test message")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "user_id") {
		t.Error("Expected user_id in log")
	}
	if !strings.Contains(string(content), "123") {
		t.Error("Expected context value in log")
	}
}

func TestLogger_With(t *testing.T) {
	logger, logPath := createTestLogger(t)

	newLogger := logger.With("request_id", "abc123")

	// Should return new logger
	if newLogger == logger {
		t.Error("With should return new logger")
	}

	newLogger.Info("test")

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "request_id") {
		t.Error("Expected request_id in log")
	}
}

func TestLogger_WithoutContext(t *testing.T) {
	logger, _ := createTestLogger(t)

	// Add context
	contextLogger := logger.WithContext(map[string]any{
		"user_id":    123,
		"session_id": "abc",
		"keep":       "this",
	})

	// Remove specific keys
	cleanLogger := contextLogger.WithoutContext("user_id", "session_id")

	// Should return new logger
	if cleanLogger == contextLogger {
		t.Error("WithoutContext should return new logger")
	}

	// Check context
	if _, ok := cleanLogger.ctx["user_id"]; ok {
		t.Error("user_id should be removed")
	}
	if _, ok := cleanLogger.ctx["session_id"]; ok {
		t.Error("session_id should be removed")
	}
	if _, ok := cleanLogger.ctx["keep"]; !ok {
		t.Error("keep should remain")
	}
}

func TestLogger_ContextPersistence(t *testing.T) {
	logger, logPath := createTestLogger(t)

	// Create logger with context
	userLogger := logger.WithContext(map[string]any{
		"user_id": 123,
	})

	// Log multiple times
	userLogger.Info("action 1")
	userLogger.Info("action 2")
	userLogger.Info("action 3")

	content, _ := os.ReadFile(logPath)
	// Count occurrences of user_id
	count := strings.Count(string(content), "user_id")
	if count != 3 {
		t.Errorf("Expected user_id in all 3 logs, got %d occurrences", count)
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": {
				Driver: "file",
				Level:  "error", // Only ERROR and above
				FileConfig: &FileConfig{
					Path: logPath,
				},
			},
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	logger, _ := manager.Channel("file")

	// These should be filtered out
	logger.Debug("debug")
	logger.Info("info")
	logger.Warning("warning")

	// These should be logged
	logger.Error("error")
	logger.Critical("critical")

	content, _ := os.ReadFile(logPath)
	logContent := string(content)

	if strings.Contains(logContent, "DEBUG") {
		t.Error("DEBUG should be filtered")
	}
	if strings.Contains(logContent, "INFO") {
		t.Error("INFO should be filtered")
	}
	if strings.Contains(logContent, "WARNING") {
		t.Error("WARNING should be filtered")
	}
	if !strings.Contains(logContent, "ERROR") {
		t.Error("ERROR should be logged")
	}
	if !strings.Contains(logContent, "CRITICAL") {
		t.Error("CRITICAL should be logged")
	}
}

func TestLogger_MultipleContext(t *testing.T) {
	logger, logPath := createTestLogger(t)

	// Pass multiple context maps
	logger.Info("test", map[string]any{"key1": "val1"}, map[string]any{"key2": "val2"})

	content, _ := os.ReadFile(logPath)
	if !strings.Contains(string(content), "key1") {
		t.Error("Expected key1 in log")
	}
	if !strings.Contains(string(content), "key2") {
		t.Error("Expected key2 in log")
	}
}

func TestMergeContext(t *testing.T) {
	result := mergeContext(
		map[string]any{"a": 1},
		map[string]any{"b": 2},
		map[string]any{"c": 3},
	)

	if len(result) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(result))
	}

	if result["a"] != 1 || result["b"] != 2 || result["c"] != 3 {
		t.Error("Expected all values to be merged")
	}
}

func TestMergeContext_Override(t *testing.T) {
	result := mergeContext(
		map[string]any{"a": 1},
		map[string]any{"a": 2}, // Override
	)

	if result["a"] != 2 {
		t.Error("Later context should override earlier")
	}
}

func TestMergeContext_Empty(t *testing.T) {
	result := mergeContext()

	if result == nil {
		t.Error("Should return empty map, not nil")
	}

	if len(result) != 0 {
		t.Error("Should return empty map")
	}
}

