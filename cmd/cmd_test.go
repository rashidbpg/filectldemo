package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "filectl") {
		t.Error("expected help output to contain 'filectl'")
	}
}

func TestCreateCommand(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	output, err := executeCommand("create", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Created:") {
		t.Errorf("expected output to contain 'Created:', got %q", output)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist")
	}
}

func TestCreateCommandWithContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	output, err := executeCommand("create", path, "--content", "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Created:") {
		t.Errorf("expected output to contain 'Created:', got %q", output)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(data) != "hello world" {
		t.Errorf("content = %q, want %q", string(data), "hello world")
	}
}

func TestCopyCommand(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.WriteFile(src, []byte("copy me"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output, err := executeCommand("copy", src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Copied:") {
		t.Errorf("expected output to contain 'Copied:', got %q", output)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(data) != "copy me" {
		t.Errorf("content = %q, want %q", string(data), "copy me")
	}
}

func TestCombineCommand(t *testing.T) {
	dir := t.TempDir()
	src1 := filepath.Join(dir, "a.txt")
	src2 := filepath.Join(dir, "b.txt")
	dst := filepath.Join(dir, "combined.txt")

	if err := os.WriteFile(src1, []byte("hello "), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(src2, []byte("world"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output, err := executeCommand("combine", src1, src2, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Combined:") {
		t.Errorf("expected output to contain 'Combined:', got %q", output)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(data) != "hello world" {
		t.Errorf("content = %q, want %q", string(data), "hello world")
	}
}

func TestDeleteCommand(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "delete_me.txt")

	if err := os.WriteFile(path, []byte("bye"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output, err := executeCommand("delete", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Deleted:") {
		t.Errorf("expected output to contain 'Deleted:', got %q", output)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDeleteCommandNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.txt")

	_, err := executeCommand("delete", path)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCommandMissingArgs(t *testing.T) {
	_, err := executeCommand("create")
	if err == nil {
		t.Error("expected error for missing args")
	}

	_, err = executeCommand("copy")
	if err == nil {
		t.Error("expected error for missing args")
	}

	_, err = executeCommand("combine")
	if err == nil {
		t.Error("expected error for missing args")
	}

	_, err = executeCommand("delete")
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestEndToEndWorkflow(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")
	combined := filepath.Join(dir, "combined.txt")
	backup := filepath.Join(dir, "backup.txt")

	// Create first file
	if _, err := executeCommand("create", file1, "--content", "line 1"); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Create second file
	if _, err := executeCommand("create", file2, "--content", "line 2"); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Combine files
	if _, err := executeCommand("combine", file1, file2, combined); err != nil {
		t.Fatalf("combine failed: %v", err)
	}

	// Copy combined file
	if _, err := executeCommand("copy", combined, backup); err != nil {
		t.Fatalf("copy failed: %v", err)
	}

	// Verify backup
	data, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(data) != "line 1line 2" {
		t.Errorf("unexpected content: %q", string(data))
	}

	// Delete original files
	if _, err := executeCommand("delete", file1); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := executeCommand("delete", file2); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify originals are gone
	if _, err := os.Stat(file1); !os.IsNotExist(err) {
		t.Error("expected file1 to be deleted")
	}
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		t.Error("expected file2 to be deleted")
	}
}
