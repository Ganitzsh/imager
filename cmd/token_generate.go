package cmd

import (
	"github.com/spf13/cobra"
)

func tokenGenerateCmdRun(cmd *cobra.Command, args []string) {
}

// tokenGenerateCmd represents the serve command
var tokenGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a new access token. Can only be used with the root token",
	Run:   tokenGenerateCmdRun,
}

func init() {
	tokenCmd.AddCommand(tokenGenerateCmd)
}
