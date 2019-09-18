package components

import (
	"github.com/labstack/echo/v4"
	"github.com/oubiwann/go-svc-conventions/cfg"
)

// Default component collection
type Default struct {
	Config *cfg.Config
	HTTP   *echo.Echo
}

// Add more components here that have more or less than what's done above. This
// is useful for testing or runnning in different binaries/executables, etc.
