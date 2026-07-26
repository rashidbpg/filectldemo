.PHONY: build test lint clean install run build-all test-integration coverage docker-build docker-test docker-demo docker-install

BINARY=filectl
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Default target
all: lint test build

# Build for current platform
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

# Run all tests
test:
	go test -v -race ./...

# Run tests with coverage
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	golangci-lint run

# Run integration tests
test-integration: build
	@echo "Running integration tests..."
	./$(BINARY) create /tmp/test-create.txt --content "integration test"
	./$(BINARY) copy /tmp/test-create.txt /tmp/test-copy.txt
	./$(BINARY) create /tmp/test-a.txt --content "hello "
	./$(BINARY) create /tmp/test-b.txt --content "world"
	./$(BINARY) combine /tmp/test-a.txt /tmp/test-b.txt /tmp/test-combined.txt
	./$(BINARY) delete /tmp/test-create.txt
	./$(BINARY) delete /tmp/test-copy.txt
	./$(BINARY) delete /tmp/test-a.txt
	./$(BINARY) delete /tmp/test-b.txt
	./$(BINARY) delete /tmp/test-combined.txt
	@echo "Integration tests passed!"

# Build for all platforms
build-all:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe .

# Install to GOPATH/bin
install: build
	cp $(BINARY) $(GOPATH)/bin/

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -rf dist/
	rm -f coverage.out coverage.html

# Run the built binary
run: build
	./$(BINARY)

# Show help
help:
	@echo "Targets:"
	@echo "  build          - Build for current platform"
	@echo "  test           - Run unit tests"
	@echo "  coverage       - Run tests with coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  test-integration - Run integration tests"
	@echo "  build-all      - Cross-compile for all platforms"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  clean          - Remove build artifacts"
	@echo "  run            - Build and run"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-test    - Run tests in Docker"
	@echo "  docker-test-go - Run Go tests in builder container"
	@echo "  docker-demo    - Interactive demo in Docker"
	@echo "  docker-install - Build and install .deb in Docker"
	@echo "  help           - Show this help"

# Docker targets
docker-build:
	docker compose build build

docker-test:
	docker compose run --rm test

docker-test-go:
	docker compose run --rm test-go

docker-demo:
	docker compose run --rm demo

docker-install:
	docker compose run --rm debian-install
