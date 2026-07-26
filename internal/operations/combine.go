package operations

import (
	"errors"
	"fmt"
	"io"
	"os"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

// CombineFiles concatenates src1 and src2 into dst.
// Returns an error if any source doesn't exist or destination cannot be created.
func CombineFiles(src1, src2, dst string) error {
	dstFile, err := os.Create(dst)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return &fileerrors.PermissionError{Path: dst, Err: err}
		}
		return &fileerrors.FileError{Op: "combine", Path: dst, Err: err}
	}
	defer dstFile.Close()

	files := []string{src1, src2}
	for _, src := range files {
		srcFile, err := os.Open(src)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return &fileerrors.NotFoundError{Path: src}
			}
			return &fileerrors.FileError{Op: "combine", Path: src, Err: err}
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			return &fileerrors.FileError{Op: "combine", Path: fmt.Sprintf("%s -> %s", src, dst), Err: err}
		}
		srcFile.Close()
	}

	return nil
}
