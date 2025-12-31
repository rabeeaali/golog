package golog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileDriver writes log entries to a file
type FileDriver struct {
	mu         sync.Mutex
	file       *os.File
	path       string
	dateFormat string
}

// NewFileDriver creates a new file driver from configuration
func NewFileDriver(config ChannelConfig) (Driver, error) {
	if config.FileConfig == nil {
		return nil, fmt.Errorf("file configuration is required")
	}

	path := config.FileConfig.Path
	if path == "" {
		path = "logs/app.log"
	}

	dateFormat := config.FileConfig.DateFormat
	if dateFormat == "" {
		dateFormat = "2006-01-02 15:04:05"
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open file for appending
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &FileDriver{
		file:       file,
		path:       path,
		dateFormat: dateFormat,
	}, nil
}

// Log writes a log entry to the file
func (d *FileDriver) Log(entry *Entry) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	formatted := d.format(entry)
	_, err := d.file.WriteString(formatted)
	return err
}

// format formats the entry for file output (Laravel-style)
func (d *FileDriver) format(entry *Entry) string {
	// Format: [2024-01-15 10:30:45] production.INFO: Message {"context":"data"}
	timestamp := entry.Timestamp.Format(d.dateFormat)
	channel := entry.Channel
	if channel == "" {
		channel = "local"
	}

	// Build the log line
	line := fmt.Sprintf("[%s] %s.%s: %s", timestamp, channel, entry.Level.String(), entry.Message)

	// Add context if present
	if len(entry.Context) > 0 {
		line += "\n"
		for key, value := range entry.Context {
			line += fmt.Sprintf("  %s: %v\n", key, formatValue(value))
		}
	}

	// Add exception if present
	if entry.Exception != nil {
		line += "\n  Exception:\n"
		line += fmt.Sprintf("    Class: %s\n", entry.Exception.Class)
		line += fmt.Sprintf("    Message: %s\n", entry.Exception.Message)
		if entry.Exception.Code != 0 {
			line += fmt.Sprintf("    Code: %d\n", entry.Exception.Code)
		}
		if entry.Exception.File != "" {
			line += fmt.Sprintf("    File: %s:%d\n", entry.Exception.File, entry.Exception.Line)
		}
		if len(entry.Exception.Trace) > 0 {
			line += "    Trace:\n"
			for i, t := range entry.Exception.Trace {
				line += fmt.Sprintf("      #%d %s\n", i, t)
				if i >= 10 {
					line += fmt.Sprintf("      ... and %d more\n", len(entry.Exception.Trace)-10)
					break
				}
			}
		}
	}

	line += "\n"
	return line
}

// formatValue formats a value for log output
func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Close closes the file
func (d *FileDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.file != nil {
		return d.file.Close()
	}
	return nil
}

// Name returns the driver name
func (d *FileDriver) Name() string {
	return "file"
}

// Flush ensures all data is written to disk
func (d *FileDriver) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.file != nil {
		return d.file.Sync()
	}
	return nil
}
