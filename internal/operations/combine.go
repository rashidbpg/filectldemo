package operations

import (
	"errors"
	"fmt"
	"io"
	"os"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

const opCombine = "combine"

// CombineFiles concatenates src1 and src2 into dst.
// Returns an error if any source doesn't exist or destination cannot be created.
func CombineFiles(src1, src2, dst string) error {
	dstFile, err := os.Create(dst)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return &fileerrors.PermissionError{Path: dst, Err: err}
		}
		return &fileerrors.FileError{Op: opCombine, Path: dst, Err: err}
	}
	defer func() { _ = dstFile.Close() }()

	files := []string{src1, src2}
	for _, src := range files {
		srcFile, err := os.Open(src)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return &fileerrors.NotFoundError{Path: src}
			}
			return &fileerrors.FileError{Op: opCombine, Path: src, Err: err}
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			_ = srcFile.Close()
			return &fileerrors.FileError{Op: opCombine, Path: fmt.Sprintf("%s -> %s", src, dst), Err: err}
		}
		_ = srcFile.Close()
	}

	return nil
}
