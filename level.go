package golog

import "strings"

// Level represents the severity of a log entry
type Level int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production
	DebugLevel Level = iota
	// InfoLevel is the default logging priority
	InfoLevel
	// NoticeLevel logs are normal but significant events
	NoticeLevel
	// WarningLevel logs are more important than Info, but don't need immediate attention
	WarningLevel
	// ErrorLevel logs are high-priority, but the application can still function
	ErrorLevel
	// CriticalLevel logs are critical conditions
	CriticalLevel
	// AlertLevel logs require immediate action
	AlertLevel
	// EmergencyLevel is when the system is unusable
	EmergencyLevel
)

// String returns the string representation of the level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case NoticeLevel:
		return "NOTICE"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	case AlertLevel:
		return "ALERT"
	case EmergencyLevel:
		return "EMERGENCY"
	default:
		return "UNKNOWN"
	}
}

// Emoji returns the emoji for the level (used for Slack)
func (l Level) Emoji() string {
	switch l {
	case DebugLevel:
		return "🔍"
	case InfoLevel:
		return "ℹ️"
	case NoticeLevel:
		return "📝"
	case WarningLevel:
		return "⚠️"
	case ErrorLevel:
		return "❌"
	case CriticalLevel:
		return "🔥"
	case AlertLevel:
		return "🚨"
	case EmergencyLevel:
		return "💀"
	default:
		return "📋"
	}
}

// Color returns the ANSI color code for the level
func (l Level) Color() string {
	switch l {
	case DebugLevel:
		return "\033[36m" // Cyan
	case InfoLevel:
		return "\033[32m" // Green
	case NoticeLevel:
		return "\033[34m" // Blue
	case WarningLevel:
		return "\033[33m" // Yellow
	case ErrorLevel:
		return "\033[31m" // Red
	case CriticalLevel:
		return "\033[35m" // Magenta
	case AlertLevel:
		return "\033[31;1m" // Bold Red
	case EmergencyLevel:
		return "\033[37;41m" // White on Red
	default:
		return "\033[0m" // Reset
	}
}

// SlackColor returns the Slack attachment color for the level
func (l Level) SlackColor() string {
	switch l {
	case DebugLevel:
		return "#36a64f" // Green
	case InfoLevel:
		return "#2196F3" // Blue
	case NoticeLevel:
		return "#9C27B0" // Purple
	case WarningLevel:
		return "#FF9800" // Orange
	case ErrorLevel:
		return "#f44336" // Red
	case CriticalLevel:
		return "#D32F2F" // Dark Red
	case AlertLevel:
		return "#B71C1C" // Darker Red
	case EmergencyLevel:
		return "#000000" // Black
	default:
		return "#9E9E9E" // Grey
	}
}

// ParseLevel parses a string into a Level
func ParseLevel(s string) Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "NOTICE":
		return NoticeLevel
	case "WARNING", "WARN":
		return WarningLevel
	case "ERROR", "ERR":
		return ErrorLevel
	case "CRITICAL", "CRIT":
		return CriticalLevel
	case "ALERT":
		return AlertLevel
	case "EMERGENCY", "EMERG":
		return EmergencyLevel
	default:
		return InfoLevel
	}
}
