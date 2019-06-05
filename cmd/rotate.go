package cmd

import (
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rotateCmd represents the serve command
var rotateCmd = &cobra.Command{
	Use:   "rotate [angle] [file]",
	Short: "Rotates the given image with the given angle. Clockwise by default",
	Run: func(cmd *cobra.Command, args []string) {
		angle, _ := strconv.Atoi(args[0])
		direction := "clockwise"
		ccw, _ := cmd.Flags().GetBool("ccw")
		if ccw {
			direction = "counter-clockwise"
		}
		logrus.Info("Rotating ", args[1], " ", angle, " degrees ", direction)
	},
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errNoArgs
		}
		if len(args) != 2 {
			return errMissingArgs
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("angle must be a number")
		}
		return nil
	}),
}

func init() {
	rootCmd.AddCommand(rotateCmd)

	rotateCmd.PersistentFlags().Bool(
		"ccw", false, "Rotates the image counter-clockwise",
	)
}
