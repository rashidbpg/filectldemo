# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-08-03

### Added
- Command help text now includes usage examples for every subcommand.
- `--version` flag now works on released builds (fixed version injection).
- `CHANGELOG.md` documenting project history.

## [1.1.0] - 2026-07-26

### Fixed
- Resolved remaining `golangci-lint` v2 lint issues (errcheck on `os.Chmod`).

## [1.0.1] - 2026-07-26

### Fixed
- Migrated `golangci-lint` configuration to v2 format.
- Fixed CI pipeline Go version and `golangci-lint` compatibility.
- Resolved all `golangci-lint` v2 lint errors.

## [1.0.0] - 2026-07-26

### Added
- Initial release of `filectl`, a file manipulation CLI tool.
- `create` - create empty files or files with content (`--content`).
- `copy` - copy a file preserving permissions.
- `combine` - concatenate two files into a third.
- `delete` - remove a file from the filesystem.
- Custom error types (`FileError`, `NotFoundError`, `AlreadyExistsError`, `PermissionError`).
- 100% unit test coverage with race detection.
- CI/CD pipeline (lint, test, integration test, cross-platform build, release).
- Cross-platform release assets (Linux/macOS tarballs, Windows zip, Debian package, SHA256 checksums).
- Docker support for testing and demos.

[Unreleased]: https://github.com/rashidbpg/filectldemo/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/rashidbpg/filectldemo/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/rashidbpg/filectldemo/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/rashidbpg/filectldemo/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/rashidbpg/filectldemo/releases/tag/v1.0.0
