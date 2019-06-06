package service

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	DevMode bool

	TLSEnabled bool
	TLSCert    string
	TLSKey     string

	MaxImageSize int64
	Port         int32
	Host         string
	BufferSize   uint32

	HTTPEnabled bool
	HTTPPort    int32
}

func NewConfig() (*Config, error) {
	c := Config{
		DevMode: viper.GetBool("DevMode"),

		TLSEnabled: viper.GetBool("TLSEnabled"),
		TLSCert:    viper.GetString("TLSCert"),
		TLSKey:     viper.GetString("TLSKey"),

		Port: viper.GetInt32("Port"),
		Host: viper.GetString("Host"),

		MaxImageSize: viper.GetInt64("MaxImageSize"),
		BufferSize:   viper.GetUint32("BufferSize"),

		HTTPEnabled: viper.GetBool("HTTPEnabled"),
		HTTPPort:    viper.GetInt32("HTTPPort"),
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
