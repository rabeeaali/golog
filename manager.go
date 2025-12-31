package golog

import (
	"fmt"
	"sync"
)

// Manager manages multiple log channels like Laravel's LogManager
type Manager struct {
	mu             sync.RWMutex
	config         *Config
	channels       map[string]*LogChannel
	defaultChannel string
	sharedContext  map[string]any
}

// LogChannel represents a logging channel with its driver and configuration
type LogChannel struct {
	name   string
	driver Driver
	level  Level
	ctx    map[string]any
}

// NewManager creates a new log manager with the given configuration
func NewManager(config *Config) (*Manager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	m := &Manager{
		config:         config,
		channels:       make(map[string]*LogChannel),
		defaultChannel: config.Default,
		sharedContext:  make(map[string]any),
	}

	return m, nil
}

// Channel returns a specific channel by name
func (m *Manager) Channel(name string) (*Logger, error) {
	m.mu.RLock()
	ch, exists := m.channels[name]
	m.mu.RUnlock()

	if exists {
		return NewLogger(ch, m), nil
	}

	// Create the channel
	ch, err := m.createChannel(name)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.channels[name] = ch
	m.mu.Unlock()

	return NewLogger(ch, m), nil
}

// createChannel creates a channel from configuration
func (m *Manager) createChannel(name string) (*LogChannel, error) {
	config, exists := m.config.Channels[name]
	if !exists {
		return nil, fmt.Errorf("channel [%s] is not defined", name)
	}

	// Handle stack driver
	if config.Driver == "stack" {
		return m.createStackChannel(name, config)
	}

	factory, exists := GetDriverFactory(config.Driver)
	if !exists {
		return nil, fmt.Errorf("driver [%s] is not supported", config.Driver)
	}

	driver, err := factory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver [%s]: %w", config.Driver, err)
	}

	level := ParseLevel(config.Level)

	return &LogChannel{
		name:   name,
		driver: driver,
		level:  level,
		ctx:    make(map[string]any),
	}, nil
}

// createStackChannel creates a stack channel that writes to multiple channels
func (m *Manager) createStackChannel(name string, config ChannelConfig) (*LogChannel, error) {
	if config.StackConfig == nil || len(config.StackConfig.Channels) == 0 {
		return nil, fmt.Errorf("stack channel [%s] requires channel list", name)
	}

	var drivers []Driver
	for _, chName := range config.StackConfig.Channels {
		chConfig, exists := m.config.Channels[chName]
		if !exists {
			return nil, fmt.Errorf("channel [%s] in stack is not defined", chName)
		}

		factory, exists := GetDriverFactory(chConfig.Driver)
		if !exists {
			return nil, fmt.Errorf("driver [%s] is not supported", chConfig.Driver)
		}

		driver, err := factory(chConfig)
		if err != nil {
			if !config.StackConfig.IgnoreExceptions {
				return nil, fmt.Errorf("failed to create driver [%s]: %w", chConfig.Driver, err)
			}
			continue
		}
		drivers = append(drivers, driver)
	}

	stackDriver := &StackDriver{
		drivers:          drivers,
		ignoreExceptions: config.StackConfig.IgnoreExceptions,
	}

	level := ParseLevel(config.Level)
	if config.Level == "" {
		level = DebugLevel
	}

	return &LogChannel{
		name:   name,
		driver: stackDriver,
		level:  level,
		ctx:    make(map[string]any),
	}, nil
}

// Default returns the default channel logger
func (m *Manager) Default() (*Logger, error) {
	return m.Channel(m.defaultChannel)
}

// SetDefault sets the default channel name
func (m *Manager) SetDefault(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultChannel = name
}

// ShareContext adds context that will be included in all log entries
func (m *Manager) ShareContext(ctx map[string]any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range ctx {
		m.sharedContext[k] = v
	}

	// Update existing channels
	for _, ch := range m.channels {
		for k, v := range ctx {
			ch.ctx[k] = v
		}
	}
}

// SharedContext returns the shared context
func (m *Manager) SharedContext() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ctx := make(map[string]any)
	for k, v := range m.sharedContext {
		ctx[k] = v
	}
	return ctx
}

// FlushSharedContext clears the shared context
func (m *Manager) FlushSharedContext() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sharedContext = make(map[string]any)
}

// Close closes all channels
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for _, ch := range m.channels {
		if err := ch.driver.Close(); err != nil {
			lastErr = err
		}
	}

	m.channels = make(map[string]*LogChannel)
	return lastErr
}

// StackDriver is a driver that writes to multiple drivers
type StackDriver struct {
	drivers          []Driver
	ignoreExceptions bool
}

// Log writes to all drivers in the stack
func (d *StackDriver) Log(entry *Entry) error {
	var lastErr error
	for _, driver := range d.drivers {
		if err := driver.Log(entry); err != nil {
			if !d.ignoreExceptions {
				lastErr = err
			}
		}
	}
	return lastErr
}

// Close closes all drivers
func (d *StackDriver) Close() error {
	var lastErr error
	for _, driver := range d.drivers {
		if err := driver.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Name returns the driver name
func (d *StackDriver) Name() string {
	return "stack"
}

