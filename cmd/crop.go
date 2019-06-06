package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	pb "github.com/ganitzsh/12fact/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CropCmdArgs struct {
	CmdArgs
	pb.CropImageRequest
	File string
}

func NewCropCmdArgs() *CropCmdArgs {
	return &CropCmdArgs{}
}

func (a *CropCmdArgs) Read(cmd *cobra.Command, args []string) error {
	topLeftX, _ := strconv.Atoi(args[1])
	topLeftY, _ := strconv.Atoi(args[2])
	width, _ := strconv.Atoi(args[3])
	height, _ := strconv.Atoi(args[4])
	a.File = args[0]
	a.TopLeftX = int32(topLeftX)
	a.TopLeftY = int32(topLeftY)
	a.Width = int32(width)
	a.Height = int32(height)
	return nil
}

func cropCmdRun(cmd *cobra.Command, args []string) {
	a := NewCropCmdArgs()
	if err := a.Read(cmd, args); err != nil {
		logrus.Errorf("Could not read arguments: %v", err)
		os.Exit(1)
	}
	logrus.Infof(
		"Cropping %s starting at (%d;%d) of size %dx%d",
		a.File, a.TopLeftX, a.TopLeftY, a.Width, a.Height,
	)
	client, err := rpcv1.NewClient()
	if err != nil {
		logrus.Errorf("Failed to get client: %v", err)
		os.Exit(1)
	}
	r, err := client.Transform(
		a.File, pb.TransformationType_CROP, &pb.CropImageRequest{
			TopLeftX: a.TopLeftX,
			TopLeftY: a.TopLeftY,
			Width:    a.Width,
			Height:   a.Height,
		})
	if err != nil {
		logrus.Errorf("Failed to crop image: %v", err)
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

// cropCmd represents the serve command
var cropCmd = &cobra.Command{
	Use:   "crop [file] [topLeftX] [topLeftY] [width] [height]",
	Short: "Crops the given [file] starting at topLeft of size [width]x[height]",
	Run:   cropCmdRun,
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errNoArgs
		}
		if len(args) != 5 {
			return errMissingArgs
		}
		validate := map[string]string{
			"topLeftX": args[1],
			"topLeftY": args[2],
			"width":    args[3],
			"height":   args[4],
		}
		for a, val := range validate {
			if _, err := strconv.Atoi(val); err != nil {
				return errors.New(fmt.Sprintf("Invalid value for [%s]", a))
			}
		}
		return nil
	}),
}

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.PersistentFlags().StringP(
		"out", "o", "", "The path to the output file",
	)
}
