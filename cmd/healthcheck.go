package cmd

import (
	"context"
	"os"

	"github.com/ganitzsh/12fact/delivery/rpcv1"
	pb "github.com/ganitzsh/12fact/delivery/rpcv1/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func healthcheckCmdRun(cmd *cobra.Command, args []string) {
	client, err := rpcv1.NewClient()
	if err != nil {
		logrus.Errorf("Failed to get client: %v", err)
		os.Exit(1)
	}
	ret, err := client.Healthcheck(
		context.Background(), &pb.HealthcheckRequest{},
	)
	if err != nil {
		logrus.Errorf("Failed to get healtcheck status: %v", err)
		os.Exit(1)
	}
	if !ret.Store {
		logrus.Error("Data store is having issues")
		os.Exit(1)
	}
}

// healthcheckCmd represents the serve command
var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Sends a request to analyze the status of the service",
	Run:   healthcheckCmdRun,
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
}
