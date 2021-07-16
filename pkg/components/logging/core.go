package logging

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/zylog/logger"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
)

// Setup ...
func Setup(cfg *config.Config) {
	logger.SetupLogging(cfg.Logging)
}

// Setup ...
func SetupClient(cfg *config.Config) {
	logger.SetupLogging(cfg.ClientLogging)
}

// Load pretends that the global is more functional in nature ...
func Load(cfg *config.Config) *log.Logger {
	Setup(cfg)
	return log.StandardLogger()
}

func LoadClient(cfg *config.Config) *log.Logger {
	SetupClient(cfg)
	return log.StandardLogger()
}
