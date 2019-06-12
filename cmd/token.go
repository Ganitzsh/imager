package cmd

import (
	"github.com/spf13/cobra"
)

// tokenCmd represents the serve command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Acceess token manipulator tool",
	Run:   nil,
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
