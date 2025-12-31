package golog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	config := &Config{
		Default: "file",
		AppName: "TestApp",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(t.TempDir(), "test.log")),
		},
	}

	manager, err := NewManager(config)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if manager == nil {
		t.Fatal("Expected manager to be created")
	}
}

func TestNewManager_NilConfig(t *testing.T) {
	manager, err := NewManager(nil)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if manager == nil {
		t.Fatal("Expected manager with default config")
	}
}

func TestManager_Channel(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	logger, err := manager.Channel("file")
	if err != nil {
		t.Fatalf("Channel failed: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be returned")
	}
}

func TestManager_Channel_NotFound(t *testing.T) {
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(t.TempDir(), "test.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	_, err := manager.Channel("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent channel")
	}
}

func TestManager_Channel_Cached(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	// Get channel twice
	logger1, _ := manager.Channel("file")
	logger2, _ := manager.Channel("file")

	// Should use same underlying channel
	if logger1.channel != logger2.channel {
		t.Error("Expected channel to be cached")
	}
}

func TestManager_Default(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file":  NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
			"other": NewFileChannelConfig(filepath.Join(tempDir, "other.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	logger, err := manager.Default()
	if err != nil {
		t.Fatalf("Default failed: %v", err)
	}

	if logger.channel.name != "file" {
		t.Errorf("Expected default channel 'file', got %q", logger.channel.name)
	}
}

func TestManager_SetDefault(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file":  NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
			"other": NewFileChannelConfig(filepath.Join(tempDir, "other.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	manager.SetDefault("other")

	logger, _ := manager.Default()
	if logger.channel.name != "other" {
		t.Errorf("Expected default channel 'other', got %q", logger.channel.name)
	}
}

func TestManager_ShareContext(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(tempDir, "test.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	// Share context before getting channel
	manager.ShareContext(map[string]any{
		"app_version": "1.0.0",
		"environment": "test",
	})

	ctx := manager.SharedContext()
	if ctx["app_version"] != "1.0.0" {
		t.Error("Expected shared context to be stored")
	}

	// New loggers should have shared context
	logger, _ := manager.Channel("file")
	if logger.ctx["app_version"] != "1.0.0" {
		t.Error("Expected logger to have shared context")
	}
}

func TestManager_FlushSharedContext(t *testing.T) {
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(filepath.Join(t.TempDir(), "test.log")),
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	manager.ShareContext(map[string]any{"key": "value"})
	manager.FlushSharedContext()

	ctx := manager.SharedContext()
	if len(ctx) != 0 {
		t.Error("Expected shared context to be flushed")
	}
}

func TestManager_Close(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")
	config := &Config{
		Default: "file",
		Channels: map[string]ChannelConfig{
			"file": NewFileChannelConfig(logPath),
		},
	}

	manager, _ := NewManager(config)

	// Create a channel to ensure it's opened
	logger, _ := manager.Channel("file")
	logger.Info("test")

	// Close manager
	if err := manager.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// File should exist
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Expected log file to exist")
	}
}

func TestManager_StackChannel(t *testing.T) {
	tempDir := t.TempDir()
	logPath1 := filepath.Join(tempDir, "test1.log")
	logPath2 := filepath.Join(tempDir, "test2.log")

	config := &Config{
		Default: "stack",
		Channels: map[string]ChannelConfig{
			"file1": NewFileChannelConfig(logPath1),
			"file2": NewFileChannelConfig(logPath2),
			"stack": {
				Driver: "stack",
				Level:  "debug",
				StackConfig: &StackConfig{
					Channels:         []string{"file1", "file2"},
					IgnoreExceptions: true,
				},
			},
		},
	}

	manager, err := NewManager(config)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	defer manager.Close()

	logger, err := manager.Channel("stack")
	if err != nil {
		t.Fatalf("Channel failed: %v", err)
	}

	logger.Info("stack message", map[string]any{"test": "value"})

	// Both files should have the log
	content1, _ := os.ReadFile(logPath1)
	content2, _ := os.ReadFile(logPath2)

	if len(content1) == 0 {
		t.Error("Expected file1 to have content")
	}

	if len(content2) == 0 {
		t.Error("Expected file2 to have content")
	}
}

func TestManager_StackChannel_NoChannels(t *testing.T) {
	config := &Config{
		Default: "stack",
		Channels: map[string]ChannelConfig{
			"stack": {
				Driver:      "stack",
				StackConfig: &StackConfig{},
			},
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	_, err := manager.Channel("stack")
	if err == nil {
		t.Error("Expected error for stack with no channels")
	}
}

func TestManager_UnsupportedDriver(t *testing.T) {
	config := &Config{
		Default: "custom",
		Channels: map[string]ChannelConfig{
			"custom": {
				Driver: "unsupported",
				Level:  "debug",
			},
		},
	}

	manager, _ := NewManager(config)
	defer manager.Close()

	_, err := manager.Channel("custom")
	if err == nil {
		t.Error("Expected error for unsupported driver")
	}
}

func TestStackDriver(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	fileConfig := ChannelConfig{
		Driver: "file",
		FileConfig: &FileConfig{
			Path: logPath,
		},
	}

	fileDriver, _ := NewFileDriver(fileConfig)

	stackDriver := &StackDriver{
		drivers:          []Driver{fileDriver},
		ignoreExceptions: false,
	}

	entry := NewEntry(InfoLevel, "stack test")
	if err := stackDriver.Log(entry); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	if stackDriver.Name() != "stack" {
		t.Errorf("Expected name 'stack', got %q", stackDriver.Name())
	}

	if err := stackDriver.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

