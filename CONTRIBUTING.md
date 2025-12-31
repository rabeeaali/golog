# Contributing to GoLog

First off, thank you for considering contributing to GoLog! 🎉

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

**When creating a bug report, include:**

- Go version (`go version`)
- GoLog version
- Operating system
- A clear description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Code samples if applicable

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:

- A clear description of the feature
- Use cases and examples
- Why this feature would be useful

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Run linter (`golangci-lint run`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/golog.git
cd golog

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run linter (install golangci-lint first)
golangci-lint run
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Add comments for exported functions and types
- Write tests for new functionality
- Keep functions small and focused

## Testing

- Write unit tests for all new code
- Aim for high test coverage
- Use table-driven tests where appropriate
- Mock external dependencies (HTTP calls, etc.)

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "input1", "expected1"},
        {"case 2", "input2", "expected2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Commit Messages

- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Keep the first line under 72 characters

Examples:

- `Add Slack driver async support`
- `Fix file driver concurrent write issue`
- `Update README with new examples`

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Focus on what's best for the community

## Questions?

Feel free to open an issue for any questions about contributing!
