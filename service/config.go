package service

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("API")

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("dev_mode", true)
	viper.SetDefault("max_image_size", "24 mib")
	viper.SetDefault("buffer_size", 2048)

	viper.SetDefault("http.enabled", false)
	viper.SetDefault("http.port", 8081)

	viper.SetDefault("store.type", StoreTypeRedis)
	viper.SetDefault("store.redis.host", "localhost:6379")
	viper.SetDefault("store.redis.user", "")
	viper.SetDefault("store.redis.password", "")
	viper.SetDefault("store.redis.db", 0)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Errorf("Failed to read config file: %v", err)
	}
}

type StoreType string

const (
	StoreTypeRedis StoreType = "redis"
)

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

type StoreConfig struct {
	Type  StoreType
	Redis *RedisConfig
}

type HTTPConfig struct {
	Enabled bool
	Port    int32
}

type Config struct {
	DevMode bool

	TLSEnabled bool
	TLSCert    string
	TLSKey     string

	MaxImageSize uint64
	Port         int32
	Host         string
	BufferSize   uint32

	Store *StoreConfig

	HTTPEnabled bool
	HTTPPort    int32
}

func NewConfig() (*Config, error) {
	imageSizeStr := viper.GetString("max_image_size")
	imageSize, err := humanize.ParseBytes(imageSizeStr)
	if err != nil {
		logrus.Errorf("Failed to read max_image_size: %v", err)
		logrus.Warn("Falling back to default (24 Mb)")
		imageSize = 25165824 // 24 MiB
	}
	c := Config{
		DevMode: viper.GetBool("dev_mode"),

		Host:         viper.GetString("host"),
		Port:         viper.GetInt32("port"),
		MaxImageSize: imageSize,
		BufferSize:   viper.GetUint32("buffer_size"),

		TLSEnabled: viper.GetBool("tls.enabled"),
		TLSCert:    viper.GetString("tls.cert"),
		TLSKey:     viper.GetString("tls.key"),

		HTTPEnabled: viper.GetBool("http.enabled"),
		HTTPPort:    viper.GetInt32("http.port"),

		Store: &StoreConfig{
			Type: StoreType(viper.GetString("store.type")),
			Redis: &RedisConfig{
				Host:     viper.GetString("store.redis.host"),
				Password: viper.GetString("store.redis.password"),
				DB:       viper.GetInt("store.redis.db"),
			},
		},
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
	if c.Store == nil {
		return ErrUnknownStoreType
	}
	switch c.Store.Type {
	case StoreTypeRedis:
		break
	default:
		return ErrUnknownStoreType
	}
	return nil
}
