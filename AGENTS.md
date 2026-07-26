# AI Context for filectl

## Project Overview
filectl is a Go CLI tool for file manipulation operations. It's designed as a single portable binary with no runtime dependencies. The project follows Go best practices with clean separation of concerns.

## Architecture
- **Language**: Go 1.22+
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Structure**: Clean separation between CLI commands (`cmd/`), business logic (`internal/operations/`), and public packages (`pkg/errors/`)

## Project Layout
```
filectl/
├── cmd/                    # CLI commands (Cobra)
├── internal/operations/    # Core business logic
├── internal/version/       # Version information
├── pkg/errors/             # Custom error types
├── .github/workflows/      # CI/CD pipeline
└── testdata/               # Test fixtures
```

## Key Commands
- `create` - Creates empty files or files with content
- `copy` - Copies files preserving permissions
- `combine` - Concatenates two files into a third
- `delete` - Removes files from filesystem

## Error Handling
Custom error types in `pkg/errors/`:
- `FileError` - General file operation error with context
- `NotFoundError` - File not found
- `AlreadyExistsError` - File already exists
- `PermissionError` - Insufficient permissions

Use `errors.As()` to check error types:
```go
var notFoundErr *fileerrors.NotFoundError
if errors.As(err, &notFoundErr) {
    // handle not found
}
```

## Testing
- Unit tests in `internal/operations/operations_test.go`
- CLI integration tests in `cmd/cmd_test.go`
- Uses table-driven tests pattern
- Uses `t.TempDir()` for isolated test environments
- Run with: `go test -v -race ./...`

## CI/CD
- GitHub Actions pipeline in `.github/workflows/ci.yml`
- Stages: lint → test → integration test → build → release
- Cross-compiles for linux/darwin/windows (amd64/arm64)
- Creates GitHub releases with:
  - Tarballs for Linux/macOS
  - Zip for Windows
  - SHA256 checksums
  - Debian packages (.deb)

## Conventions
- Errors are wrapped with `fmt.Errorf("context: %w", err)`
- All file operations check for existence before proceeding
- File permissions default to 0644
- Commands follow Cobra patterns (Args validation, RunE handlers)
- Version info injected via ldflags at build time

## Build Commands
```bash
make build              # Build for current platform
make test               # Run tests with race detection
make coverage           # Run tests with coverage report
make lint               # Run golangci-lint
make test-integration   # Run integration tests
make build-all          # Cross-compile for all platforms
make install            # Install to GOPATH/bin
make help               # Show all available targets
```

## Dependencies
- github.com/spf13/cobra - CLI framework
- No other external runtime dependencies

## Adding New Commands
1. Create `cmd/newcmd.go` with cobra.Command
2. Create `internal/operations/newop.go` with business logic
3. Add unit tests in `internal/operations/operations_test.go`
4. Add CLI tests in `cmd/cmd_test.go`
5. Register command in `cmd/root.go` via `init()`
