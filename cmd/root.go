package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

// SetVersionInfo sets the version information for the CLI.
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

var rootCmd = &cobra.Command{
	Use:   "filectl",
	Short: "A file manipulation CLI tool",
	Long: `filectl provides commands for creating, copying, combining, and deleting files.

It is designed as a single portable binary with no runtime dependencies.`,
	SilenceUsage: true,
}

func init() {
	rootCmd.Version = version
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
