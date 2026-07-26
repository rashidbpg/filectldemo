package operations

import (
	"errors"
	"fmt"
	"io"
	"os"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

var osFileStat = (*os.File).Stat

const opCopy = "copy"

// CopyFile copies a file from src to dst, preserving permissions.
// Returns an error if the source doesn't exist or destination cannot be created.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &fileerrors.NotFoundError{Path: src}
		}
		return &fileerrors.FileError{Op: opCopy, Path: src, Err: err}
	}
	defer func() { _ = srcFile.Close() }()

	srcInfo, err := osFileStat(srcFile)
	if err != nil {
		return &fileerrors.FileError{Op: opCopy, Path: src, Err: err}
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return &fileerrors.PermissionError{Path: dst, Err: err}
		}
		return &fileerrors.FileError{Op: opCopy, Path: dst, Err: err}
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return &fileerrors.FileError{Op: opCopy, Path: fmt.Sprintf("%s -> %s", src, dst), Err: err}
	}

	return nil
}
