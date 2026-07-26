package operations

import (
	"errors"
	"os"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

// DeleteFile removes a file from the filesystem.
// Returns an error if the file doesn't exist or cannot be removed.
func DeleteFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &fileerrors.NotFoundError{Path: path}
	}

	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return &fileerrors.PermissionError{Path: path, Err: err}
		}
		return &fileerrors.FileError{Op: "delete", Path: path, Err: err}
	}

	return nil
}

// ForceDeleteFile removes a file, ignoring "not found" errors.
// Useful for cleanup operations.
func ForceDeleteFile(path string) error {
	err := os.Remove(path)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if errors.Is(err, os.ErrPermission) {
		return &fileerrors.PermissionError{Path: path, Err: err}
	}
	return &fileerrors.FileError{Op: "delete", Path: path, Err: err}
}
