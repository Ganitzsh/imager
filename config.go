package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.12fact")
	viper.AddConfigPath("/etc/12fact/")

	viper.SetEnvPrefix("12fact")

	viper.SetDefault("port", 8080)
}

type Config struct {
	DevMode bool

	TLSEnabled bool
	TLSCert    string
	TLSKey     string

	MaxImageSize int64
}

func NewConfig() (*Config, error) {
	c := Config{
		DevMode: viper.GetBool("dev_mode"),

		TLSEnabled: viper.GetBool("tls_enabled"),
		TLSCert:    viper.GetString("tls_cert"),
		TLSKey:     viper.GetString("tls_key"),

		MaxImageSize: viper.GetInt64("max_image_size"),
	}
	return &c, c.validate()
}

func (c Config) validate() error {
	if !c.DevMode {
		if c.TLSEnabled && (c.TLSKey == "" || c.TLSCert == "") {
			return ErrInvalidTLSconfiguration
		}
	} else {
		logrus.Warn("This instance is running in dev mode")
	}
	return nil
}
