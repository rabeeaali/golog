.PHONY: all test lint build clean examples help

# Default target
all: lint test build

# Run all tests
test:
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests passed!"

# Run tests with coverage report
coverage: test
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Running go vet instead..."; \
		go vet ./...; \
	fi
	@echo "✅ Lint passed!"

# Build the package
build:
	@echo "🔨 Building..."
	go build -v ./...
	@echo "✅ Build successful!"

# Build examples
examples:
	@echo "📦 Building examples..."
	@for dir in _examples/*/; do \
		echo "Building $$dir..."; \
		cd "$$dir" && go build -v . && cd ../..; \
	done
	@echo "✅ Examples built!"

# Run example
run-basic:
	@echo "🚀 Running basic example..."
	cd _examples/basic && go run main.go

run-slack:
	@echo "🚀 Running slack example..."
	cd _examples/slack && go run main.go

run-context:
	@echo "🚀 Running context example..."
	cd _examples/context && go run main.go

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f coverage.out coverage.html
	rm -rf logs/
	go clean
	@echo "✅ Clean!"

# Install development dependencies
dev-deps:
	@echo "📥 Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Dependencies installed!"

# Format code
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...
	@echo "✅ Formatted!"

# Show help
help:
	@echo "GoLog Makefile Commands:"
	@echo ""
	@echo "  make          - Run lint, test, and build"
	@echo "  make test     - Run all tests"
	@echo "  make coverage - Run tests with coverage report"
	@echo "  make lint     - Run linter"
	@echo "  make build    - Build the package"
	@echo "  make examples - Build all examples"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make fmt      - Format code"
	@echo "  make dev-deps - Install development dependencies"
	@echo ""
	@echo "Run examples:"
	@echo "  make run-basic   - Run basic example"
	@echo "  make run-slack   - Run slack example"
	@echo "  make run-context - Run context example"

