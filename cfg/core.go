package cfg

import (
	"strings"

	cfg "github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	logger "github.com/geomyidia/zylog/logger"
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

// GRPCDConfig ...
type GRPCDConfig struct {
	Host string
	Port int
}

// Config ...
type Config struct {
	HTTPD HTTPDConfig
	DB    FileDBConfig
	GRPCD GRPCDConfig
	Logging *logger.ZyLogOptions
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
		HTTPD: HTTPDConfig{
			Host: cfg.GetString("httpd.host"),
			Port: cfg.GetInt("httpd.port"),
		},
		DB: FileDBConfig{
			Directory: cfg.GetString("database.file-based.directory"),
		},
		GRPCD: GRPCDConfig{
			Host: cfg.GetString("grpc.host"),
			Port: cfg.GetInt("grpc.port"),
		},
		Logging: &logger.ZyLogOptions{
			Colored:      cfg.GetBool("logging.colored"),
			Level:        cfg.GetString("logging.level"),
			Output:       cfg.GetString("logging.output"),
			ReportCaller: cfg.GetBool("logging.report-caller"),
		},
	}
}
