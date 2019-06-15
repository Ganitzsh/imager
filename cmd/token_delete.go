package cmd

import (
	"github.com/spf13/cobra"
)

func tokenDeleteCmdRun(cmd *cobra.Command, args []string) {

}

// tokenDeleteCmd represents the serve command
var tokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an access token. Can only be used with the root token",
	Run:   tokenDeleteCmdRun,
}

func init() {
	tokenCmd.AddCommand(tokenDeleteCmd)
}
