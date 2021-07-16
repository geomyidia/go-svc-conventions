package httpd

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
)

// SetupServer ...
func SetupServer(cfg *config.Config) *HTTPHandlerServer {
	log.Debug("Setting up HTTP daemon ...")
	s := NewHTTPHandlerServer(cfg.HTTPD)
	log.Debug("HTTP daemon set up.")
	return s
}
