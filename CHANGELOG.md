# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2024-XX-XX

### Added

- Initial release
- File driver for logging to files
- Slack driver for logging to Slack webhooks
- Stack driver for logging to multiple channels
- Laravel-style logging API
- Context support for structured logging
- Exception/error logging with stack traces
- Multiple channel support
- Async Slack message sending
- Log level filtering
- Shared context across channels
- Comprehensive test suite with 89%+ coverage

### Drivers

- **File Driver**
  - Write logs to files
  - Laravel-style formatting
  - Configurable date format
- **Slack Driver**

  - Send logs to Slack webhooks
  - Beautiful formatted attachments
  - Color-coded by log level
  - Emoji indicators
  - Async support
  - Multiple channels support

- **Stack Driver**
  - Log to multiple drivers simultaneously
  - Exception ignoring option

### Configuration

- Functional options pattern for driver configuration
- YAML/JSON compatible config structs
- Default configuration with sensible defaults

[Unreleased]: https://github.com/rabeeaali/golog/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/rabeeaali/golog/releases/tag/v1.0.0
