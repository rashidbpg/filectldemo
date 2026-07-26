package errors

import "fmt"

// FileError represents a file operation error with context.
type FileError struct {
	Op   string // operation name (create, copy, combine, delete)
	Path string // file path involved
	Err  error  // underlying error
}

func (e *FileError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// NotFoundError indicates a file was not found.
type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("file not found: %s", e.Path)
}

// AlreadyExistsError indicates a file already exists.
type AlreadyExistsError struct {
	Path string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("file already exists: %s", e.Path)
}

// PermissionError indicates insufficient permissions.
type PermissionError struct {
	Path string
	Err  error
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("permission denied: %s", e.Path)
}

func (e *PermissionError) Unwrap() error {
	return e.Err
}
