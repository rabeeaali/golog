package golog

import "errors"

var (
	// ErrNotInitialized is returned when the log manager is not initialized
	ErrNotInitialized = errors.New("golog: manager not initialized, call Init() first")

	// ErrChannelNotFound is returned when a channel is not configured
	ErrChannelNotFound = errors.New("golog: channel not found")

	// ErrDriverNotSupported is returned when a driver is not supported
	ErrDriverNotSupported = errors.New("golog: driver not supported")
)

