package main

import (
	"github.com/rashidbpg/filectldemo/cmd"
)

// These variables are set at build time via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
