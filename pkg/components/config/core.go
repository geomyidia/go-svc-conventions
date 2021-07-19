package config

import (
	"fmt"
	"strings"

	logger "github.com/geomyidia/zylog/logger"
	log "github.com/sirupsen/logrus"
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

// DBConfig ...
type DBConfig struct {
	Directory string
}

// HTTPDConfig ...
type HTTPDConfig struct {
	Host           string
	Port           int
	RequestLogging bool
}

// ConnectionString ...
func (c *HTTPDConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GRPCDConfig ...
type GRPCDConfig struct {
	Host string
	Port int
}

// ConnectionString ...
func (c *GRPCDConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Config ...
type Config struct {
	HTTPD         *HTTPDConfig
	DB            *DBConfig
	GRPCD         *GRPCDConfig
	Logging       *logger.ZyLogOptions
	ClientLogging *logger.ZyLogOptions
}

// NewConfig is a constructor that creates the full coniguration data structure
// for use by our application(s) and client(s) as an in-memory copy of the
// config data (saving from having to make repeated and somewhat expensive
// calls to the viper library).
//
// Note that Viper does provide both the AllSettings() and Unmarshall()
// functions, but these require that you have a struct defined that will be
// used to dump the Viper config data into. We've already got that set up, so
// there's no real benefit to switching.
//
// Furthermore, in our case, we're utilizing structs from other libraries to
// be used when setting those up (see how we initialize the logging component
// in ./components/logging.go, Setup).
func NewConfig() *Config {
	return &Config{
		HTTPD: &HTTPDConfig{
			Host:           cfg.GetString("httpd.host"),
			Port:           cfg.GetInt("httpd.port"),
			RequestLogging: cfg.GetBool("httpd.request-logging"),
		},
		DB: &DBConfig{
			Directory: cfg.GetString("db.directory"),
		},
		GRPCD: &GRPCDConfig{
			Host: cfg.GetString("grpc.host"),
			Port: cfg.GetInt("grpc.port"),
		},
		Logging: &logger.ZyLogOptions{
			Colored:      cfg.GetBool("logging.colored"),
			Level:        cfg.GetString("logging.level"),
			Output:       cfg.GetString("logging.output"),
			ReportCaller: cfg.GetBool("logging.report-caller"),
		},
		ClientLogging: &logger.ZyLogOptions{
			Colored:      cfg.GetBool("client-logging.colored"),
			Level:        cfg.GetString("client-logging.level"),
			Output:       cfg.GetString("client-logging.output"),
			ReportCaller: cfg.GetBool("client-logging.report-caller"),
		},
	}
}
