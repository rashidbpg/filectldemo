package version

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version != "dev" {
		t.Errorf("Version = %q, want %q", info.Version, "dev")
	}
	if info.Commit != "none" {
		t.Errorf("Commit = %q, want %q", info.Commit, "none")
	}
	if info.Date != "unknown" {
		t.Errorf("Date = %q, want %q", info.Date, "unknown")
	}
	if info.GoVer == "" {
		t.Error("GoVer should not be empty")
	}
	if info.OS == "" {
		t.Error("OS should not be empty")
	}
	if info.Arch == "" {
		t.Error("Arch should not be empty")
	}
}

func TestString(t *testing.T) {
	info := Info{
		Version: "1.0.0",
		Commit:  "abc123",
		Date:    "2024-01-01",
		GoVer:   "go1.22.0",
		OS:      "linux",
		Arch:    "amd64",
	}

	s := info.String()
	if !strings.Contains(s, "1.0.0") {
		t.Errorf("String() should contain version, got %q", s)
	}
	if !strings.Contains(s, "abc123") {
		t.Errorf("String() should contain commit, got %q", s)
	}
	if !strings.Contains(s, "2024-01-01") {
		t.Errorf("String() should contain date, got %q", s)
	}
	if !strings.Contains(s, "linux") {
		t.Errorf("String() should contain OS, got %q", s)
	}
	if !strings.Contains(s, "amd64") {
		t.Errorf("String() should contain arch, got %q", s)
	}
	if !strings.Contains(s, "go1.22.0") {
		t.Errorf("String() should contain Go version, got %q", s)
	}
}
