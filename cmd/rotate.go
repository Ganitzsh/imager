package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/davecgh/go-spew/spew"
	pb "github.com/ganitzsh/12fact/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/golang/protobuf/ptypes"
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
		client, err := service.NewClient()
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		ext := filepath.Ext(args[1])
		if ext == "" {
			logrus.Error(errors.New("Unknown extension"))
			os.Exit(1)
		}
		f, err := os.Open(args[1])
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		stat, err := f.Stat()
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		stream, err := client.TransformImage(context.Background())
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		data, err := ptypes.MarshalAny(&pb.RotateImageRequest{
			Angle: int32(angle),
		})
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		if err := stream.Send(&pb.TransformImageRequest{
			Type: pb.TransformationType_ROTATE,
			Data: data,
			Image: &pb.Image{
				Size:   stat.Size(),
				Format: ext,
			},
		}); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		buff := make([]byte, 2048)
		for {
			_, e := f.Read(buff)
			if e == io.EOF {
				break
			}
			if err := stream.Send(&pb.TransformImageRequest{
				Image: &pb.Image{
					File: buff,
				},
			}); err != nil {
				logrus.Error(err)
				os.Exit(1)
			}
			time.Sleep(500 * time.Microsecond)
		}
		if err := stream.CloseSend(); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		spew.Dump(stream.Recv())
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
