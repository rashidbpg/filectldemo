package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Info contains version information.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	GoVer   string `json:"go_version"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
}

// Get returns the current version info.
func Get() Info {
	return Info{
		Version: version,
		Commit:  commit,
		Date:    date,
		GoVer:   runtime.Version(),
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
}

// String returns a human-readable version string.
func (v Info) String() string {
	return fmt.Sprintf("filectl %s (commit: %s, built: %s, %s/%s, %s)",
		v.Version, v.Commit, v.Date, v.OS, v.Arch, v.GoVer)
}
