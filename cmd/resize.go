package cmd

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/ganitzsh/12fact/delivery/rpcv1"
	pb "github.com/ganitzsh/12fact/delivery/rpcv1/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ResizeCmdArgs struct {
	CmdArgs
	pb.ResizeImageRequest
	File string
}

func NewResizeCmdArgs() *ResizeCmdArgs {
	return &ResizeCmdArgs{}
}

func (a *ResizeCmdArgs) Read(cmd *cobra.Command, args []string) error {
	var width, height int
	var err error

	if width, err = strconv.Atoi(args[1]); err != nil {
		return err
	}
	if height, err = strconv.Atoi(args[2]); err != nil {
		return err
	}
	a.Width = int32(width)
	a.Height = int32(height)
	a.File = args[0]
	return nil
}

func resizeCmdRun(cmd *cobra.Command, args []string) {
	a := NewResizeCmdArgs()
	if err := a.Read(cmd, args); err != nil {
		logrus.Errorf("Could not read arguments: %v", err)
		os.Exit(1)
	}
	logrus.Infof("Resizing %s to %dx%d", a.File, a.Width, a.Height)
	client, err := rpcv1.NewClient()
	if err != nil {
		logrus.Errorf("Failed to get client: %v", err)
		os.Exit(1)
	}
	r, err := client.Transform(
		a.File, pb.TransformationType_RESIZE, &pb.ResizeImageRequest{
			Width:  a.Width,
			Height: a.Height,
		})
	if err != nil {
		logrus.Errorf("Failed to resize image: %v", err)
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

// resizeCmd represents the serve command
var resizeCmd = &cobra.Command{
	Use:   "resize [file] [width] [height]",
	Short: "Resizes the given image to fit [width]x[height]",
	Long:  "If [height] or [width] is 0 then the aspect ration will be preserved",
	Run:   resizeCmdRun,
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errNoArgs
		}
		if len(args) != 3 {
			return errMissingArgs
		}
		if _, err := strconv.Atoi(args[1]); err != nil {
			return err
		}
		if _, err := strconv.Atoi(args[2]); err != nil {
			return err
		}
		return nil
	}),
}

func init() {
	rootCmd.AddCommand(resizeCmd)

	resizeCmd.PersistentFlags().StringP(
		"out", "o", "", "The path to the output file",
	)
}
