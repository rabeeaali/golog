package golog

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Entry represents a log entry with all its data
type Entry struct {
	// Message is the log message
	Message string `json:"message"`

	// Level is the severity level
	Level Level `json:"level"`

	// Timestamp is when the log was created
	Timestamp time.Time `json:"timestamp"`

	// Context contains additional structured data
	Context map[string]any `json:"context,omitempty"`

	// Exception contains error/exception details if applicable
	Exception *ExceptionInfo `json:"exception,omitempty"`

	// Channel is the name of the log channel
	Channel string `json:"channel,omitempty"`
}

// ExceptionInfo contains structured exception/error information
type ExceptionInfo struct {
	Class   string   `json:"class"`
	Message string   `json:"message"`
	Code    int      `json:"code,omitempty"`
	File    string   `json:"file,omitempty"`
	Line    int      `json:"line,omitempty"`
	Trace   []string `json:"trace,omitempty"`
}

// NewEntry creates a new log entry
func NewEntry(level Level, message string) *Entry {
	return &Entry{
		Message:   message,
		Level:     level,
		Timestamp: time.Now(),
		Context:   make(map[string]any),
	}
}

// WithContext adds context data to the entry
func (e *Entry) WithContext(ctx map[string]any) *Entry {
	for k, v := range ctx {
		e.Context[k] = v
	}
	return e
}

// With adds a single context key-value pair
func (e *Entry) With(key string, value any) *Entry {
	e.Context[key] = value
	return e
}

// WithError adds error information to the entry
func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}

	e.Exception = &ExceptionInfo{
		Class:   getErrorType(err),
		Message: err.Error(),
	}

	// Capture stack trace
	e.Exception.Trace = captureStackTrace(3) // Skip captureStackTrace, WithError, and caller

	// Try to get file and line from stack
	if _, file, line, ok := runtime.Caller(1); ok {
		e.Exception.File = file
		e.Exception.Line = line
	}

	return e
}

// WithException adds detailed exception information
func (e *Entry) WithException(class, message string, code int, file string, line int, trace []string) *Entry {
	e.Exception = &ExceptionInfo{
		Class:   class,
		Message: message,
		Code:    code,
		File:    file,
		Line:    line,
		Trace:   trace,
	}
	return e
}

// SetChannel sets the channel name
func (e *Entry) SetChannel(channel string) *Entry {
	e.Channel = channel
	return e
}

// getErrorType returns the type name of an error
func getErrorType(err error) string {
	if err == nil {
		return ""
	}
	t := fmt.Sprintf("%T", err)
	// Clean up common prefixes
	t = strings.TrimPrefix(t, "*")
	return t
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) []string {
	var trace []string
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		// Skip runtime frames
		if !strings.Contains(frame.File, "runtime/") {
			trace = append(trace, fmt.Sprintf("%s:%d (%s)", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
		// Limit trace length
		if len(trace) >= 20 {
			break
		}
	}
	return trace
}

// ToJSON converts the entry to JSON
func (e *Entry) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// ContextJSON returns the context as pretty-printed JSON
func (e *Entry) ContextJSON() string {
	if len(e.Context) == 0 {
		return ""
	}
	b, err := json.MarshalIndent(e.Context, "", "    ")
	if err != nil {
		return fmt.Sprintf("%v", e.Context)
	}
	return string(b)
}

// ExceptionJSON returns the exception as pretty-printed JSON
func (e *Entry) ExceptionJSON() string {
	if e.Exception == nil {
		return ""
	}
	b, err := json.MarshalIndent(e.Exception, "", "    ")
	if err != nil {
		return fmt.Sprintf("%v", e.Exception)
	}
	return string(b)
}
