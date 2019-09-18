package cfg

import (
	"log"
	"strings"

	cfg "github.com/spf13/viper"
)

// Configuration related constants
const (
	AppName         string = "app"
	ConfigDir       string = "configs"
	ConfigFile      string = "app"
	ConfigType      string = "yaml"
	ConfigReadError string = "Fatal error config file"
)

func init() {
	cfg.AddConfigPath(ConfigDir)
	cfg.SetConfigName(ConfigFile)
	cfg.SetConfigType(ConfigType)
	cfg.SetEnvPrefix(AppName)

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	cfg.Set("Verbose", true)
	cfg.AutomaticEnv()
	cfg.AddConfigPath("/")

	err := cfg.ReadInConfig()
	if err != nil {
		// log.Panic is not used here, since logging depends ...
		log.Panicf("%s: %s", ConfigReadError, err)
	}
}

// HTTPDConfig ...
type HTTPDConfig struct {
	Host string
	Port int
}

// FileDBConfig ...
type FileDBConfig struct {
	Directory string
}

// GRPCConfig ...
type GRPCConfig struct {
}

// Config ...
type Config struct {
	HTTPD HTTPDConfig
	DB    FileDBConfig
	GRPC  GRPCConfig
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		HTTPD: HTTPDConfig{
			Host: cfg.GetString("httpd.host"),
			Port: cfg.GetInt("httpd.port"),
		},
		DB: FileDBConfig{
			Directory: cfg.GetString("database.file-based.directory"),
		},
	}
}
