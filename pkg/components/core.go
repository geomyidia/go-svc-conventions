package components

import (
	logger "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
)

// Base component collection
type Base struct {
	Config *config.Config
	Logger *logger.Logger
}

// TestBase component that keeps stdout clean
type TestBase struct {
	Config *config.Config
}

// Default component collection
type Default struct {
	Base
}

// Add more components here that have more or less than what's done above. This
// is useful for testing or runnning in different binaries/executables, etc.

// Application ...
type Application struct {
	Default
	Bus *msgbus.MsgBus
}
