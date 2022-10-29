package cmd

import (
	"github.com/spf13/cobra"
)

var orphanCmd = &cobra.Command{
	Use:   "orphan",
	Short: "Delete macOS ._* generated files that are left alone. This does not delete the .DS_Store file.",
	RunE:  runOrphan,
}

func init() {
	rootCmd.AddCommand(orphanCmd)
}

func runOrphan(cmd *cobra.Command, args []string) error {
	return nil
}
