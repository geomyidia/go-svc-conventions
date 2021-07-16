package httpd

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components"
)

// SetupServer ...
func SetupServer(app *components.Application) *HTTPHandlerServer {
	log.Debug("Setting up HTTP daemon ...")
	s := NewHTTPHandlerServer(app)
	log.Debug("HTTP daemon set up.")
	return s
}
