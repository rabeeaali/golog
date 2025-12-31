package golog

import (
	"testing"
)

func TestGetDriverFactory(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		expected bool
	}{
		{"file driver exists", "file", true},
		{"slack driver exists", "slack", true},
		{"unknown driver", "unknown", false},
		{"custom driver", "custom", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := GetDriverFactory(tt.driver)
			if exists != tt.expected {
				t.Errorf("GetDriverFactory(%q) exists = %v, want %v", tt.driver, exists, tt.expected)
			}
		})
	}
}

func TestRegisterDriver(t *testing.T) {
	// Register a custom driver
	RegisterDriver("custom", func(config ChannelConfig) (Driver, error) {
		return &mockDriver{name: "custom"}, nil
	})

	factory, exists := GetDriverFactory("custom")
	if !exists {
		t.Error("Expected custom driver to be registered")
	}

	driver, err := factory(ChannelConfig{})
	if err != nil {
		t.Fatalf("Factory failed: %v", err)
	}

	if driver.Name() != "custom" {
		t.Errorf("Expected driver name 'custom', got %q", driver.Name())
	}

	// Clean up
	delete(driverFactories, "custom")
}

func TestRegisterDriver_Override(t *testing.T) {
	// Save original
	originalFile := driverFactories["file"]

	// Override file driver
	RegisterDriver("file", func(config ChannelConfig) (Driver, error) {
		return &mockDriver{name: "overridden"}, nil
	})

	factory, _ := GetDriverFactory("file")
	driver, _ := factory(ChannelConfig{})

	if driver.Name() != "overridden" {
		t.Error("Expected driver to be overridden")
	}

	// Restore original
	driverFactories["file"] = originalFile
}

// Mock driver for testing
type mockDriver struct {
	name    string
	entries []*Entry
}

func (d *mockDriver) Log(entry *Entry) error {
	d.entries = append(d.entries, entry)
	return nil
}

func (d *mockDriver) Close() error {
	return nil
}

func (d *mockDriver) Name() string {
	return d.name
}

