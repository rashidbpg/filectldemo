package cmd

import (
	"github.com/rashidbpg/filectldemo/internal/operations"
	"github.com/spf13/cobra"
)

var combineCmd = &cobra.Command{
	Use:   "combine <source1> <source2> <destination>",
	Short: "Combine two files into a third",
	Long:  `Combine the contents of two files into a single destination file.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		src1, src2, dst := args[0], args[1], args[2]
		if err := operations.CombineFiles(src1, src2, dst); err != nil {
			return err
		}
		cmd.Printf("Combined: %s + %s -> %s\n", src1, src2, dst)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(combineCmd)
}
