package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ganitzsh/12fact/delivery/httpv1"
	"github.com/ganitzsh/12fact/delivery/rpcv1"
	"github.com/ganitzsh/12fact/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errNoArgs      = errors.New("No arguments given, please read the help")
	errMissingArgs = errors.New("Missing arguments")
)

type CmdArgs interface {
	Read(cmd *cobra.Command, args []string) error
}

var (
	flagCfgFile string
	flagPort    int32
)

var rootCmd = &cobra.Command{
	Use:   "12fact",
	Short: "12fact is a simple application that can be used as a microservice",
	Long: `It will run as a server by default and the same binary can be used to
	execute some command.

	In order to use the cli you need to have setup a token either as a environment
	variable or directly in the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := service.NewConfig()
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		if cfg.DevMode {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if cfg.HTTPEnabled {
			go httpv1.NewHTTPServerV1(cfg).ListenAndServe()
		}
		if err := rpcv1.NewRPCServer(cfg).ListenAndServe(); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(
		&flagCfgFile, "config", "c", "", "config file (default is ./config.yml)",
	)
	rootCmd.PersistentFlags().Int32P(
		"port", "p", 0, "port on which the server will listen",
	)
	rootCmd.PersistentFlags().StringP(
		"addr", "a", "", "The server's address",
	)
	rootCmd.PersistentFlags().Bool(
		"http", false, "",
	)
}

func initConfig() {
	viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("API")

	viper.SetDefault("Port", 8080)
	viper.SetDefault("DevMode", true)
	viper.SetDefault("MaxImageSize", 25165824) // Default is 24MB
	viper.SetDefault("BufferSize", 2048)
	viper.SetDefault("HTTPEnabled", false)
	viper.SetDefault("HTTPPort", 8081)
	viper.BindPFlag("Port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("Host", rootCmd.Flags().Lookup("addr"))
	viper.BindPFlag("HTTPEnabled", rootCmd.Flags().Lookup("http"))
	viper.SetConfigFile(flagCfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	}
}
