package operations

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

const testFileName = "test.txt"

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		content string
		wantErr bool
		errType interface{}
	}{
		{
			name:    "empty file",
			path:    testFileName,
			content: "",
			wantErr: false,
		},
		{
			name:    "file with content",
			path:    testFileName,
			content: "hello world",
			wantErr: false,
		},
		{
			name:    "file with multiline content",
			path:    testFileName,
			content: "line 1\nline 2\nline 3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, tt.path)

			if err := CreateFile(path, tt.content); (err != nil) != tt.wantErr {
				t.Fatalf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile() error = %v", err)
				}
				if string(data) != tt.content {
					t.Errorf("content = %q, want %q", string(data), tt.content)
				}
			}
		})
	}
}

func TestCreateFileAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, testFileName)

	if err := CreateFile(path, ""); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	err := CreateFile(path, "")
	if err == nil {
		t.Fatal("expected error for existing file")
	}

	var existsErr *fileerrors.AlreadyExistsError
	if !errors.As(err, &existsErr) {
		t.Errorf("expected AlreadyExistsError, got %T", err)
	}
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "simple copy",
			content: "copy me",
			wantErr: false,
		},
		{
			name:    "empty file",
			content: "",
			wantErr: false,
		},
		{
			name:    "large content",
			content: string(make([]byte, 1024*1024)), // 1MB
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src.txt")
			dst := filepath.Join(dir, "dst.txt")

			if err := CreateFile(src, tt.content); err != nil {
				t.Fatalf("CreateFile() error = %v", err)
			}

			if err := CopyFile(src, dst); (err != nil) != tt.wantErr {
				t.Fatalf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				data, err := os.ReadFile(dst)
				if err != nil {
					t.Fatalf("ReadFile() error = %v", err)
				}
				if string(data) != tt.content {
					t.Errorf("content = %q, want %q", string(data), tt.content)
				}
			}
		})
	}
}

func TestCopyFileStatError(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := CreateFile(src, "content"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	origStat := osFileStat
	defer func() { osFileStat = origStat }()
	osFileStat = func(f *os.File) (os.FileInfo, error) {
		return nil, errors.New("stat failed")
	}

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected stat error")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCopyFileOpenGenericError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	subdir := filepath.Join(dir, "noperm")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	src := filepath.Join(subdir, "src.txt")
	if err := os.WriteFile(src, []byte("content"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Chmod(subdir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	dst := filepath.Join(dir, "dst.txt")
	err := CopyFile(src, dst)

	_ = os.Chmod(subdir, 0755)

	if err == nil {
		t.Fatal("expected error for unreadable source")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCopyFileDirectorySource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "srcdir")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected error when source is a directory")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCopyFileNotFound(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "nonexistent.txt")
	dst := filepath.Join(dir, "dst.txt")

	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected error for missing source")
	}

	var notFoundErr *fileerrors.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestCombineFiles(t *testing.T) {
	tests := []struct {
		name     string
		content1 string
		content2 string
		want     string
		wantErr  bool
	}{
		{
			name:     "simple combine",
			content1: "hello ",
			content2: "world",
			want:     "hello world",
			wantErr:  false,
		},
		{
			name:     "empty files",
			content1: "",
			content2: "",
			want:     "",
			wantErr:  false,
		},
		{
			name:     "with newlines",
			content1: "line1\n",
			content2: "line2",
			want:     "line1\nline2",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			src1 := filepath.Join(dir, "a.txt")
			src2 := filepath.Join(dir, "b.txt")
			dst := filepath.Join(dir, "combined.txt")

			if err := CreateFile(src1, tt.content1); err != nil {
				t.Fatalf("CreateFile() error = %v", err)
			}
			if err := CreateFile(src2, tt.content2); err != nil {
				t.Fatalf("CreateFile() error = %v", err)
			}

			if err := CombineFiles(src1, src2, dst); (err != nil) != tt.wantErr {
				t.Fatalf("CombineFiles() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				data, err := os.ReadFile(dst)
				if err != nil {
					t.Fatalf("ReadFile() error = %v", err)
				}
				if string(data) != tt.want {
					t.Errorf("content = %q, want %q", string(data), tt.want)
				}
			}
		})
	}
}

func TestCombineFilesGenericError(t *testing.T) {
	dir := t.TempDir()
	src1 := filepath.Join(dir, "a.txt")
	src2 := filepath.Join(dir, "b.txt")
	dst := filepath.Join(dir, "combined.txt")

	if err := CreateFile(src1, "hello"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	if err := os.Mkdir(src2, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	err := CombineFiles(src1, src2, dst)
	if err == nil {
		t.Fatal("expected error when source is a directory")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCombineFilesNotFound(t *testing.T) {
	dir := t.TempDir()
	src1 := filepath.Join(dir, "a.txt")
	src2 := filepath.Join(dir, "nonexistent.txt")
	dst := filepath.Join(dir, "combined.txt")

	if err := CreateFile(src1, "hello"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	err := CombineFiles(src1, src2, dst)
	if err == nil {
		t.Fatal("expected error for missing source")
	}

	var notFoundErr *fileerrors.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestDeleteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "delete_me.txt")

	if err := CreateFile(path, "bye"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	if err := DeleteFile(path); err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDeleteFileNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.txt")

	err := DeleteFile(path)
	if err == nil {
		t.Fatal("expected error for missing file")
	}

	var notFoundErr *fileerrors.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestForceDeleteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "delete_me.txt")

	// Should not error on non-existent file
	if err := ForceDeleteFile(path); err != nil {
		t.Fatalf("ForceDeleteFile() error = %v", err)
	}

	// Create and delete
	if err := CreateFile(path, "bye"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	if err := ForceDeleteFile(path); err != nil {
		t.Fatalf("ForceDeleteFile() error = %v", err)
	}
}

func TestCreateFilePermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	subdir := filepath.Join(dir, "readonly")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.Chmod(subdir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	path := filepath.Join(subdir, testFileName)
	err := CreateFile(path, "content")

	_ = os.Chmod(subdir, 0755)

	if err == nil {
		t.Fatal("expected permission error")
	}

	var permErr *fileerrors.PermissionError
	if !errors.As(err, &permErr) {
		t.Errorf("expected PermissionError, got %T: %v", err, err)
	}
}

func TestCreateFileGenericError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent", testFileName)

	err := CreateFile(path, "content")
	if err == nil {
		t.Fatal("expected error for nonexistent parent directory")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCopyFilePermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dstDir := filepath.Join(dir, "readonly")
	if err := os.Mkdir(dstDir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.Chmod(dstDir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	if err := CreateFile(src, "content"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	dst := filepath.Join(dstDir, "dst.txt")
	err := CopyFile(src, dst)

	_ = os.Chmod(dstDir, 0755)

	if err == nil {
		t.Fatal("expected permission error")
	}

	var permErr *fileerrors.PermissionError
	if !errors.As(err, &permErr) {
		t.Errorf("expected PermissionError, got %T: %v", err, err)
	}
}

func TestCopyFileGenericDstError(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	if err := CreateFile(src, "content"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	dst := filepath.Join(dir, "nonexistent", "dst.txt")
	err := CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected error for nonexistent destination parent")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCombineFilesPermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	src1 := filepath.Join(dir, "a.txt")
	src2 := filepath.Join(dir, "b.txt")
	dstDir := filepath.Join(dir, "readonly")
	if err := os.Mkdir(dstDir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.Chmod(dstDir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	if err := CreateFile(src1, "hello"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}
	if err := CreateFile(src2, "world"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	dst := filepath.Join(dstDir, "combined.txt")
	err := CombineFiles(src1, src2, dst)

	_ = os.Chmod(dstDir, 0755)

	if err == nil {
		t.Fatal("expected permission error")
	}

	var permErr *fileerrors.PermissionError
	if !errors.As(err, &permErr) {
		t.Errorf("expected PermissionError, got %T: %v", err, err)
	}
}

func TestCombineFilesGenericDstError(t *testing.T) {
	dir := t.TempDir()
	src1 := filepath.Join(dir, "a.txt")
	src2 := filepath.Join(dir, "b.txt")

	if err := CreateFile(src1, "hello"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}
	if err := CreateFile(src2, "world"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	dst := filepath.Join(dir, "nonexistent", "combined.txt")
	err := CombineFiles(src1, src2, dst)
	if err == nil {
		t.Fatal("expected error for nonexistent destination parent")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestCombineFilesGenericSrcError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	src1 := filepath.Join(dir, "a.txt")
	dst := filepath.Join(dir, "combined.txt")

	if err := CreateFile(src1, "hello"); err != nil {
		t.Fatalf("CreateFile() error = %v", err)
	}

	src2Dir := filepath.Join(dir, "noperm")
	if err := os.Mkdir(src2Dir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	src2Blocked := filepath.Join(src2Dir, "b.txt")
	if err := os.WriteFile(src2Blocked, []byte("world"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Chmod(src2Dir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	err := CombineFiles(src1, src2Blocked, dst)

	_ = os.Chmod(src2Dir, 0755)

	if err == nil {
		t.Fatal("expected error for unreadable source")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestDeleteFilePermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	subdir := filepath.Join(dir, "readonly")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	path := filepath.Join(subdir, "file.txt")
	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Chmod(subdir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	err := DeleteFile(path)

	_ = os.Chmod(subdir, 0755)

	if err == nil {
		t.Fatal("expected permission error")
	}

	var permErr *fileerrors.PermissionError
	if !errors.As(err, &permErr) {
		t.Errorf("expected PermissionError, got %T: %v", err, err)
	}
}

func TestDeleteFileGenericError(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "notempty")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "child"), []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := DeleteFile(subdir)
	if err == nil {
		t.Fatal("expected error for non-empty directory")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}

func TestForceDeleteFilePermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	dir := t.TempDir()
	subdir := filepath.Join(dir, "readonly")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	path := filepath.Join(subdir, "file.txt")
	if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Chmod(subdir, 0000); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}

	err := ForceDeleteFile(path)

	os.Chmod(subdir, 0755)

	if err == nil {
		t.Fatal("expected permission error")
	}

	var permErr *fileerrors.PermissionError
	if !errors.As(err, &permErr) {
		t.Errorf("expected PermissionError, got %T: %v", err, err)
	}
}

func TestForceDeleteFileGenericError(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "notempty")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "child"), []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := ForceDeleteFile(subdir)
	if err == nil {
		t.Fatal("expected error for non-empty directory")
	}

	var fileErr *fileerrors.FileError
	if !errors.As(err, &fileErr) {
		t.Errorf("expected FileError, got %T: %v", err, err)
	}
}
