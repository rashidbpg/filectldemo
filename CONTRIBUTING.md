# Contributing to filectl

Thank you for your interest in contributing to filectl! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.22 or later
- golangci-lint
- Git

### Setting Up the Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/filectl.git
   cd filectl
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/rashidbpg/filectldemo.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build and test:
   ```bash
   make build
   make test
   ```

## Development Workflow

### Branch Naming

Use descriptive branch names:
- `feature/add-new-command` - for new features
- `fix/resolve-copy-error` - for bug fixes
- `docs/update-readme` - for documentation changes
- `refactor/improve-error-handling` - for refactoring

### Making Changes

1. Create a new branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding standards below

3. Write or update tests as needed

4. Run the full test suite:
   ```bash
   make test
   make lint
   make test-integration
   ```

5. Commit your changes with a clear message:
   ```bash
   git commit -m "Add feature: describe what you added"
   ```

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Open a Pull Request against `main`

### Commit Messages

Follow conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(create): add --force flag to overwrite existing files
fix(copy): handle permission errors gracefully
docs: update installation instructions
```

## Coding Standards

### Go Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` and `goimports` to format code
- Run `golangci-lint` before committing

### Error Handling

- Use custom error types from `pkg/errors`
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Return meaningful error messages

### Testing

- Write table-driven tests
- Use `t.TempDir()` for test isolation
- Aim for >80% code coverage
- Include both unit and integration tests

### Documentation

- Add comments to exported functions
- Update README.md for user-facing changes
- Update AGENTS.md for AI context changes

## Pull Request Guidelines

### Before Submitting

- [ ] Code compiles without errors
- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Integration tests pass (`make test-integration`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format

### PR Description

Include:
1. What the change does
2. Why the change is needed
3. How to test the changes
4. Any breaking changes

### Review Process

1. At least one approval required
2. All CI checks must pass
3. Address review feedback
4. Squash and merge

## Reporting Issues

### Bug Reports

Include:
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version)

### Feature Requests

Include:
- Use case
- Proposed solution
- Alternatives considered

## Release Process

1. Update version in code (if applicable)
2. Update CHANGELOG.md
3. Create a release PR
4. After merge, tag the release:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
5. CI/CD will automatically build and publish

## Questions?

Open an issue or reach out to the maintainers.
