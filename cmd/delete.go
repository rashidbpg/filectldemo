package cmd

import (
	"github.com/rashidbpg/filectldemo/internal/operations"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <file>",
	Short: "Delete a file",
	Long: `Delete a file from the filesystem.

The file is removed permanently. An error is returned if the file
does not exist.`,
	Example: `  filectl delete notes.txt
  filectl delete /tmp/old-data.txt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if err := operations.DeleteFile(path); err != nil {
			return err
		}
		cmd.Printf("Deleted: %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
