package cmd

import (
	"github.com/rashidbpg/filectldemo/internal/operations"
	"github.com/spf13/cobra"
)

var content string

var createCmd = &cobra.Command{
	Use:   "create <file>",
	Short: "Create a new file",
	Long:  `Create a new file, optionally with initial content.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if err := operations.CreateFile(path, content); err != nil {
			return err
		}
		cmd.Printf("Created: %s\n", path)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&content, "content", "c", "", "Initial file content")
	rootCmd.AddCommand(createCmd)
}
