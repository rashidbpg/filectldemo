package errors

import (
	"errors"
	"testing"
)

func TestFileErrorWithPath(t *testing.T) {
	inner := errors.New("inner error")
	fe := &FileError{Op: "create", Path: "/tmp/test.txt", Err: inner}
	expected := "create /tmp/test.txt: inner error"
	if got := fe.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestFileErrorWithoutPath(t *testing.T) {
	inner := errors.New("inner error")
	fe := &FileError{Op: "create", Err: inner}
	expected := "create: inner error"
	if got := fe.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestFileErrorUnwrap(t *testing.T) {
	inner := errors.New("inner error")
	fe := &FileError{Op: "create", Path: "/tmp/test.txt", Err: inner}
	if !errors.Is(fe, inner) {
		t.Error("expected Unwrap to return inner error")
	}
}

func TestNotFoundError(t *testing.T) {
	nf := &NotFoundError{Path: "/tmp/missing.txt"}
	expected := "file not found: /tmp/missing.txt"
	if got := nf.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestAlreadyExistsError(t *testing.T) {
	ae := &AlreadyExistsError{Path: "/tmp/exists.txt"}
	expected := "file already exists: /tmp/exists.txt"
	if got := ae.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestPermissionError(t *testing.T) {
	inner := errors.New("permission denied")
	pe := &PermissionError{Path: "/tmp/readonly.txt", Err: inner}
	expected := "permission denied: /tmp/readonly.txt"
	if got := pe.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestPermissionErrorUnwrap(t *testing.T) {
	inner := errors.New("permission denied")
	pe := &PermissionError{Path: "/tmp/readonly.txt", Err: inner}
	if !errors.Is(pe, inner) {
		t.Error("expected Unwrap to return inner error")
	}
}
