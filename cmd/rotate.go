package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	pb "github.com/ganitzsh/12fact/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RotateCmdArgs struct {
	CmdArgs
	pb.RotateImageRequest
	File      string
	Direction string
}

func NewRotateCmdArgs() *RotateCmdArgs {
	return &RotateCmdArgs{}
}

func (a *RotateCmdArgs) Read(cmd *cobra.Command, args []string) error {
	angle, _ := strconv.Atoi(args[1])
	a.Angle = int32(angle)
	a.Direction = "counter-clockwise"
	cw, _ := cmd.Flags().GetBool("cw")
	if cw {
		a.Direction = "clockwise"
	}
	a.File = args[0]
	return nil
}

func rotateCmdRun(cmd *cobra.Command, args []string) {
	a := NewRotateCmdArgs()
	if err := a.Read(cmd, args); err != nil {
		logrus.Errorf("Could not read arguments: %v", err)
		os.Exit(1)
	}
	logrus.Info("Rotating ", a.File, " ", a.Angle, " degrees ", a.Direction)
	client, err := rpcv1.NewClient()
	if err != nil {
		logrus.Errorf("Failed to get client: %v", err)
		os.Exit(1)
	}
	r, err := client.Transform(
		a.File, pb.TransformationType_ROTATE, &pb.RotateImageRequest{
			Angle:     a.Angle,
			ClockWise: a.ClockWise,
		})
	if err != nil {
		logrus.Errorf("Failed to transform image: %v", err)
		os.Exit(1)
	}
	outPath := "./out" + filepath.Ext(a.File)
	if outFlag, _ := cmd.Flags().GetString("out"); outFlag != "" {
		outPath = outFlag
	}
	if err := client.SaveToFile(r, outPath); err != nil {
		logrus.Errorf("Failed to save to %v: %ev", outPath, err)
	}
}

// rotateCmd represents the serve command
var rotateCmd = &cobra.Command{
	Use:   "rotate [file] [angle]",
	Short: "Rotates the given [file] with the given [angle]. Clockwise by default",
	Run:   rotateCmdRun,
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
		"cw", false, "Rotates the image clockwise",
	)
	rotateCmd.PersistentFlags().StringP(
		"out", "o", "", "The path to the output file",
	)
}
