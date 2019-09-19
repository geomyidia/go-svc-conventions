package logging

import (
	cfg "github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	logger "github.com/geomyidia/zylog/logger"
)

func init() {
	logger.SetupLogging(&logger.ZyLogOptions{
		Colored:      cfg.GetBool("logging.colored"),
		Level:        cfg.GetString("logging.level"),
		Output:       cfg.GetString("logging.output"),
		ReportCaller: cfg.GetBool("logging.report-caller"),
	})
}

// Load pretends that the global is more functional in nature ...
func Load() *log.Logger {
	return log.StandardLogger()
}