# filectl

A file manipulation CLI tool that provides simple commands for common file operations. Designed as a single portable binary with no runtime dependencies.

## Features

- **Create** empty files or files with content
- **Copy** files preserving permissions
- **Combine** two files into a third
- **Delete** files from the filesystem

## Installation

### From GitHub Releases (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/rashidbpg/filectldemo/releases).

#### Linux (Debian/Ubuntu)

```bash
# Download the .deb package
wget https://github.com/rashidbpg/filectldemo/releases/latest/download/filectl_<version>_amd64.deb

# Install with dpkg
sudo dpkg -i filectl_<version>_amd64.deb

# Or using apt (if repository is configured)
sudo apt install ./filectl_<version>_amd64.deb
```

#### Linux (Tarball)

```bash
# Download and extract
wget https://github.com/rashidbpg/filectldemo/releases/latest/download/filectl-linux-amd64.tar.gz
tar xzf filectl-linux-amd64.tar.gz

# Move to PATH
sudo mv filectl-linux-amd64 /usr/local/bin/filectl
sudo chmod +x /usr/local/bin/filectl
```

#### macOS

```bash
# Using Homebrew (coming soon)
# brew install rashidbpg/tap/filectldemo

# Or download manually
curl -L https://github.com/rashidbpg/filectldemo/releases/latest/download/filectl-darwin-arm64.tar.gz | tar xz
sudo mv filectl-darwin-arm64 /usr/local/bin/filectl
```

### From Source

```bash
# Clone the repository
git clone https://github.com/rashidbpg/filectldemo.git
cd filectl

# Build and install
make build
sudo make install
```

## Usage

```bash
# Show help
filectl --help

# Create an empty file
filectl create myfile.txt

# Create a file with content
filectl create myfile.txt --content "Hello World"

# Copy a file
filectl copy source.txt destination.txt

# Combine two files into a third
filectl combine file1.txt file2.txt combined.txt

# Delete a file
filectl delete myfile.txt
```

## Commands

| Command   | Description                                    | Example                                      |
|-----------|------------------------------------------------|----------------------------------------------|
| `create`  | Create a new file (empty or with content)      | `filectl create test.txt --content "hello"`  |
| `copy`    | Copy a file from source to destination         | `filectl copy src.txt dst.txt`               |
| `combine` | Combine two files into a third file            | `filectl combine a.txt b.txt out.txt`        |
| `delete`  | Delete a file from the filesystem              | `filectl delete test.txt`                    |

## Development

### Prerequisites

- Go 1.22 or later
- golangci-lint (for linting)

### Setup

```bash
# Clone the repository
git clone https://github.com/rashidbpg/filectldemo.git
cd filectl

# Install dependencies
go mod download
```

### Common Commands

```bash
# Build the binary
make build

# Run unit tests
make test

# Run tests with coverage
make coverage

# Run linter
make lint

# Run integration tests
make test-integration

# Cross-compile for all platforms
make build-all

# Install to GOPATH/bin
make install

# Show all available targets
make help
```

### Docker

Run and test in a Debian Linux container without installing Go locally:

```bash
# Run CLI integration tests on Debian bookworm
make docker-test

# Run Go unit tests in the builder container
make docker-test-go

# Run interactive demo
make docker-demo

# Build the Docker image
make docker-build

# Build and install .deb package in Docker
make docker-install
```

Or use docker compose directly:

```bash
# CLI integration tests
docker compose run --rm test

# Go unit tests
docker compose run --rm test-go

# Interactive demo
docker compose run --rm demo
```

## Project Structure

```
filectl/
├── cmd/                        # CLI command definitions
│   ├── root.go                 # Root command setup
│   ├── create.go               # Create command
│   ├── copy.go                 # Copy command
│   ├── combine.go              # Combine command
│   ├── delete.go               # Delete command
│   └── cmd_test.go             # CLI integration tests
├── internal/                   # Internal packages
│   ├── operations/             # Core business logic
│   │   ├── create.go
│   │   ├── copy.go
│   │   ├── combine.go
│   │   ├── delete.go
│   │   └── operations_test.go  # Unit tests
│   └── version/                # Version information
│       └── version.go
├── pkg/                        # Public packages
│   └── errors/                 # Custom error types
│       └── errors.go
├── .github/
│   └── workflows/
│       └── ci.yml              # CI/CD pipeline
├── .golangci.yml               # Linter configuration
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
├── main.go                     # Entry point
├── README.md                   # This file
└── AGENTS.md                   # AI context
```

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration and delivery:

1. **Lint** - Runs `golangci-lint` for code quality checks
2. **Test** - Runs unit tests with race detection and coverage
3. **Integration Test** - Builds and runs the binary with real commands
4. **Build** - Cross-compiles for multiple platforms (linux/darwin/windows, amd64/arm64)
5. **Release** - Creates GitHub releases with:
   - Compressed tarballs for Linux and macOS
   - Zip archive for Windows
   - SHA256 checksums
   - Debian package (.deb) for Ubuntu/Debian

### Releasing

To create a new release:

```bash
# Tag a new version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

The CI/CD pipeline will automatically build and publish the release.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
