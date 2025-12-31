// Package golog provides a Laravel-inspired logging system for Go
// with support for multiple drivers (file, Slack) and channels.
//
// Basic usage:
//
//	// Initialize with default configuration
//	golog.Init(nil)
//	defer golog.Close()
//
//	// Log to default channel
//	golog.Info("User logged in", map[string]any{
//	    "user_id": 123,
//	    "ip": "192.168.1.1",
//	})
//
//	// Log to specific channel
//	log, _ := golog.Channel("slack-alerts")
//	log.Error("Critical error occurred", map[string]any{
//	    "error": "Database connection failed",
//	})
package golog

import "sync"

var (
	defaultManager *Manager
	once           sync.Once
	mu             sync.RWMutex
)

// Init initializes the global log manager with the given configuration
func Init(config *Config) error {
	var err error
	once.Do(func() {
		defaultManager, err = NewManager(config)
	})
	return err
}

// SetManager sets a custom manager as the default
func SetManager(m *Manager) {
	mu.Lock()
	defer mu.Unlock()
	defaultManager = m
}

// GetManager returns the default manager
func GetManager() *Manager {
	mu.RLock()
	defer mu.RUnlock()
	return defaultManager
}

// Channel returns a logger for the specified channel
func Channel(name string) (*Logger, error) {
	m := GetManager()
	if m == nil {
		return nil, ErrNotInitialized
	}
	return m.Channel(name)
}

// Default returns the default channel logger
func Default() (*Logger, error) {
	m := GetManager()
	if m == nil {
		return nil, ErrNotInitialized
	}
	return m.Default()
}

// Close closes the global log manager
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if defaultManager != nil {
		err := defaultManager.Close()
		defaultManager = nil
		return err
	}
	return nil
}

// ShareContext adds context to be shared across all channels
func ShareContext(ctx map[string]any) {
	if m := GetManager(); m != nil {
		m.ShareContext(ctx)
	}
}

// --- Convenience logging functions using default channel ---

// Debug logs a debug message to the default channel
func Debug(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Debug(message, context...)
	}
}

// Info logs an info message to the default channel
func Info(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Info(message, context...)
	}
}

// Notice logs a notice message to the default channel
func Notice(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Notice(message, context...)
	}
}

// Warning logs a warning message to the default channel
func Warning(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Warning(message, context...)
	}
}

// Error logs an error message to the default channel
func Error(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Error(message, context...)
	}
}

// ErrorWithException logs an error with exception details to the default channel
func ErrorWithException(message string, err error, context ...map[string]any) {
	if log, logErr := Default(); logErr == nil {
		log.ErrorWithException(message, err, context...)
	}
}

// Critical logs a critical message to the default channel
func Critical(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Critical(message, context...)
	}
}

// CriticalWithException logs a critical message with exception details
func CriticalWithException(message string, err error, context ...map[string]any) {
	if log, logErr := Default(); logErr == nil {
		log.CriticalWithException(message, err, context...)
	}
}

// Alert logs an alert message to the default channel
func Alert(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Alert(message, context...)
	}
}

// Emergency logs an emergency message to the default channel
func Emergency(message string, context ...map[string]any) {
	if log, err := Default(); err == nil {
		log.Emergency(message, context...)
	}
}

