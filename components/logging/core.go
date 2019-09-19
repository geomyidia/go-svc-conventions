package logging

import (
	"github.com/geomyidia/go-svc-conventions/cfg"
	log "github.com/sirupsen/logrus"
	logger "github.com/geomyidia/zylog/logger"
)

// Setup ...
func Setup(config *cfg.Config) {
	logger.SetupLogging(config.Logging)
}

// Load pretends that the global is more functional in nature ...
func Load(config *cfg.Config) *log.Logger {
	Setup(config)
	return log.StandardLogger()
}