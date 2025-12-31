package golog

import "testing"

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{NoticeLevel, "NOTICE"},
		{WarningLevel, "WARNING"},
		{ErrorLevel, "ERROR"},
		{CriticalLevel, "CRITICAL"},
		{AlertLevel, "ALERT"},
		{EmergencyLevel, "EMERGENCY"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"  DEBUG  ", DebugLevel},
		{"info", InfoLevel},
		{"INFO", InfoLevel},
		{"notice", NoticeLevel},
		{"warning", WarningLevel},
		{"WARN", WarningLevel},
		{"error", ErrorLevel},
		{"ERR", ErrorLevel},
		{"critical", CriticalLevel},
		{"CRIT", CriticalLevel},
		{"alert", AlertLevel},
		{"emergency", EmergencyLevel},
		{"EMERG", EmergencyLevel},
		{"unknown", InfoLevel}, // defaults to INFO
		{"", InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLevel(tt.input); got != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLevel_Emoji(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{DebugLevel, "🔍"},
		{InfoLevel, "ℹ️"},
		{ErrorLevel, "❌"},
		{EmergencyLevel, "💀"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := tt.level.Emoji(); got != tt.want {
				t.Errorf("Level.Emoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_SlackColor(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{DebugLevel, "#36a64f"},
		{InfoLevel, "#2196F3"},
		{ErrorLevel, "#f44336"},
		{EmergencyLevel, "#000000"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := tt.level.SlackColor(); got != tt.want {
				t.Errorf("Level.SlackColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_Color(t *testing.T) {
	// Just ensure colors are non-empty ANSI codes
	levels := []Level{DebugLevel, InfoLevel, WarningLevel, ErrorLevel, CriticalLevel}
	for _, l := range levels {
		if got := l.Color(); got == "" {
			t.Errorf("Level.Color() for %v should not be empty", l)
		}
	}
}
