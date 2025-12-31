package golog

// Driver is the interface that all log drivers must implement
type Driver interface {
	// Log writes a log entry
	Log(entry *Entry) error

	// Close closes the driver and releases any resources
	Close() error

	// Name returns the driver name
	Name() string
}

// DriverFactory creates a driver from configuration
type DriverFactory func(config ChannelConfig) (Driver, error)

// Built-in driver factories
var driverFactories = map[string]DriverFactory{
	"file":  NewFileDriver,
	"slack": NewSlackDriver,
}

// RegisterDriver registers a custom driver factory
func RegisterDriver(name string, factory DriverFactory) {
	driverFactories[name] = factory
}

// GetDriverFactory returns the factory for a driver type
func GetDriverFactory(name string) (DriverFactory, bool) {
	factory, ok := driverFactories[name]
	return factory, ok
}
