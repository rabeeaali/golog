package golog

import "sync"

// Logger provides logging methods for a specific channel
type Logger struct {
	channel *LogChannel
	manager *Manager
	mu      sync.RWMutex
	ctx     map[string]any
}

// NewLogger creates a new logger for a channel
func NewLogger(channel *LogChannel, manager *Manager) *Logger {
	ctx := make(map[string]any)

	// Copy channel context
	for k, v := range channel.ctx {
		ctx[k] = v
	}

	// Copy shared context from manager
	for k, v := range manager.SharedContext() {
		ctx[k] = v
	}

	return &Logger{
		channel: channel,
		manager: manager,
		ctx:     ctx,
	}
}

// WithContext returns a new logger with additional context
func (l *Logger) WithContext(ctx map[string]any) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newCtx := make(map[string]any)
	for k, v := range l.ctx {
		newCtx[k] = v
	}
	for k, v := range ctx {
		newCtx[k] = v
	}

	return &Logger{
		channel: l.channel,
		manager: l.manager,
		ctx:     newCtx,
	}
}

// With adds a single key-value pair to context and returns a new logger
func (l *Logger) With(key string, value any) *Logger {
	return l.WithContext(map[string]any{key: value})
}

// WithoutContext removes specific keys from context
func (l *Logger) WithoutContext(keys ...string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newCtx := make(map[string]any)
	for k, v := range l.ctx {
		newCtx[k] = v
	}
	for _, key := range keys {
		delete(newCtx, key)
	}

	return &Logger{
		channel: l.channel,
		manager: l.manager,
		ctx:     newCtx,
	}
}

// log writes a log entry at the given level
func (l *Logger) log(level Level, message string, context map[string]any) {
	// Check if level meets minimum
	if level < l.channel.level {
		return
	}

	entry := NewEntry(level, message)
	entry.SetChannel(l.channel.name)

	// Add context
	l.mu.RLock()
	for k, v := range l.ctx {
		entry.Context[k] = v
	}
	l.mu.RUnlock()

	for k, v := range context {
		entry.Context[k] = v
	}

	// Write to driver
	_ = l.channel.driver.Log(entry)
}

// logWithError writes a log entry with error information
func (l *Logger) logWithError(level Level, message string, err error, context map[string]any) {
	// Check if level meets minimum
	if level < l.channel.level {
		return
	}

	entry := NewEntry(level, message)
	entry.SetChannel(l.channel.name)
	entry.WithError(err)

	// Add context
	l.mu.RLock()
	for k, v := range l.ctx {
		entry.Context[k] = v
	}
	l.mu.RUnlock()

	for k, v := range context {
		entry.Context[k] = v
	}

	// Write to driver
	_ = l.channel.driver.Log(entry)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(DebugLevel, message, ctx)
}

// Info logs an info message
func (l *Logger) Info(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(InfoLevel, message, ctx)
}

// Notice logs a notice message
func (l *Logger) Notice(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(NoticeLevel, message, ctx)
}

// Warning logs a warning message
func (l *Logger) Warning(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(WarningLevel, message, ctx)
}

// Error logs an error message
func (l *Logger) Error(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(ErrorLevel, message, ctx)
}

// ErrorWithException logs an error message with exception details
func (l *Logger) ErrorWithException(message string, err error, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.logWithError(ErrorLevel, message, err, ctx)
}

// Critical logs a critical message
func (l *Logger) Critical(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(CriticalLevel, message, ctx)
}

// CriticalWithException logs a critical message with exception details
func (l *Logger) CriticalWithException(message string, err error, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.logWithError(CriticalLevel, message, err, ctx)
}

// Alert logs an alert message
func (l *Logger) Alert(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(AlertLevel, message, ctx)
}

// AlertWithException logs an alert message with exception details
func (l *Logger) AlertWithException(message string, err error, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.logWithError(AlertLevel, message, err, ctx)
}

// Emergency logs an emergency message
func (l *Logger) Emergency(message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(EmergencyLevel, message, ctx)
}

// EmergencyWithException logs an emergency message with exception details
func (l *Logger) EmergencyWithException(message string, err error, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.logWithError(EmergencyLevel, message, err, ctx)
}

// Log logs a message at the specified level
func (l *Logger) Log(level Level, message string, context ...map[string]any) {
	ctx := mergeContext(context...)
	l.log(level, message, ctx)
}

// mergeContext merges multiple context maps
func mergeContext(contexts ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, ctx := range contexts {
		for k, v := range ctx {
			result[k] = v
		}
	}
	return result
}

