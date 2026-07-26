package operations

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	fileerrors "github.com/rashidbpg/filectldemo/pkg/errors"
)

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
			path:    "test.txt",
			content: "",
			wantErr: false,
		},
		{
			name:    "file with content",
			path:    "test.txt",
			content: "hello world",
			wantErr: false,
		},
		{
			name:    "file with multiline content",
			path:    "test.txt",
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
	path := filepath.Join(dir, "test.txt")

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
		name     string
		content  string
		wantErr  bool
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
		name    string
		content1 string
		content2 string
		want    string
		wantErr bool
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
