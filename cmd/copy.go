package cmd

import (
	"github.com/rashidbpg/filectldemo/internal/operations"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy <source> <destination>",
	Short: "Copy a file",
	Long: `Copy a file from source to destination.

The destination file inherits the source file's permissions. An error
is returned if the source does not exist.`,
	Example: `  filectl copy notes.txt notes-backup.txt
  filectl copy /tmp/data.txt ./data.txt`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src, dst := args[0], args[1]
		if err := operations.CopyFile(src, dst); err != nil {
			return err
		}
		cmd.Printf("Copied: %s -> %s\n", src, dst)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
