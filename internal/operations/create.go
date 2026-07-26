package operations

import (
	"errors"
	"os"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

// CreateFile creates a new file at the given path with optional content.
// Returns an error if the file already exists or cannot be created.
func CreateFile(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		return &fileerrors.AlreadyExistsError{Path: path}
	}

	var data []byte
	if content != "" {
		data = []byte(content)
	}

	if err := os.WriteFile(path, data, 0644); err != nil { //nolint:gosec // 0644 is intentional for user-created files
		if errors.Is(err, os.ErrPermission) {
			return &fileerrors.PermissionError{Path: path, Err: err}
		}
		return &fileerrors.FileError{Op: "create", Path: path, Err: err}
	}

	return nil
}
