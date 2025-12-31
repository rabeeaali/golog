package golog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFileDriver(t *testing.T) {
	// Create temp directory for tests
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := ChannelConfig{
		Driver: "file",
		Level:  "debug",
		FileConfig: &FileConfig{
			Path:       logPath,
			DateFormat: "2006-01-02 15:04:05",
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}
	defer driver.Close()

	if driver.Name() != "file" {
		t.Errorf("Expected driver name 'file', got %q", driver.Name())
	}
}

func TestNewFileDriver_NoConfig(t *testing.T) {
	config := ChannelConfig{
		Driver: "file",
	}

	_, err := NewFileDriver(config)
	if err == nil {
		t.Error("Expected error for missing FileConfig")
	}
}

func TestNewFileDriver_DefaultPath(t *testing.T) {
	tempDir := t.TempDir()

	config := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path: filepath.Join(tempDir, "logs", "app.log"),
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}
	defer driver.Close()

	// Directory should be created
	dirPath := filepath.Join(tempDir, "logs")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Error("Expected log directory to be created")
	}
}

func TestFileDriver_Log(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path:       logPath,
			DateFormat: "2006-01-02 15:04:05",
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}

	entry := NewEntry(InfoLevel, "test message")
	entry.SetChannel("test")
	entry.WithContext(map[string]any{
		"user_id": 123,
		"action":  "test",
	})

	if err := driver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	driver.Close()

	// Read the log file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify log content
	if !strings.Contains(logContent, "test.INFO") {
		t.Error("Log should contain 'test.INFO'")
	}

	if !strings.Contains(logContent, "test message") {
		t.Error("Log should contain 'test message'")
	}

	if !strings.Contains(logContent, "user_id") {
		t.Error("Log should contain context 'user_id'")
	}

	if !strings.Contains(logContent, "123") {
		t.Error("Log should contain context value '123'")
	}
}

func TestFileDriver_LogWithException(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path:       logPath,
			DateFormat: "2006-01-02 15:04:05",
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}

	entry := NewEntry(ErrorLevel, "database error")
	entry.SetChannel("test")
	entry.WithException("DatabaseError", "connection timeout", 500, "/app/db.go", 42, []string{
		"/app/main.go:10",
		"/app/handler.go:25",
	})

	if err := driver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	driver.Close()

	// Read the log file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify exception content
	if !strings.Contains(logContent, "Exception") {
		t.Error("Log should contain 'Exception'")
	}

	if !strings.Contains(logContent, "DatabaseError") {
		t.Error("Log should contain exception class")
	}

	if !strings.Contains(logContent, "connection timeout") {
		t.Error("Log should contain exception message")
	}

	if !strings.Contains(logContent, "Trace") {
		t.Error("Log should contain trace")
	}
}

func TestFileDriver_MultipleLogs(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path: logPath,
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}

	// Log multiple entries
	for i := 0; i < 10; i++ {
		entry := NewEntry(InfoLevel, "message")
		entry.With("index", i)
		if err := driver.Log(entry); err != nil {
			t.Fatalf("Log failed: %v", err)
		}
	}

	driver.Close()

	// Read and count lines
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Should have multiple log entries
	lines := strings.Split(string(content), "\n")
	nonEmpty := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmpty++
		}
	}

	if nonEmpty < 10 {
		t.Errorf("Expected at least 10 non-empty lines, got %d", nonEmpty)
	}
}

func TestFileDriver_Flush(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path: logPath,
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}
	defer driver.Close()

	entry := NewEntry(InfoLevel, "test message")
	driver.Log(entry)

	// Flush should not error
	fd := driver.(*FileDriver)
	if err := fd.Flush(); err != nil {
		t.Errorf("Flush failed: %v", err)
	}
}

func TestFileDriver_ConcurrentWrites(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	config := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path: logPath,
		},
	}

	driver, err := NewFileDriver(config)
	if err != nil {
		t.Fatalf("NewFileDriver failed: %v", err)
	}

	// Write concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			entry := NewEntry(InfoLevel, "concurrent message")
			entry.With("goroutine", idx)
			driver.Log(entry)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	driver.Close()

	// File should exist and have content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file should have content")
	}
}
