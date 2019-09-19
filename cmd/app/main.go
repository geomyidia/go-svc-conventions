package main

import (
	"github.com/labstack/echo/v4"

	"github.com/geomyidia/go-svc-conventions/app"
	"github.com/geomyidia/go-svc-conventions/cfg"
	logger "github.com/geomyidia/go-svc-conventions/components/logging"
	"github.com/geomyidia/reverb"
)

func main() {
	// Create the application objects and assign components to it
	a := new(app.Application)
	a.Config = cfg.NewConfig()
	a.Logger = logger.Load(a.Config)
	a.HTTPD = echo.New()
	a.GRPCD = reverb.New()

	// Perform application setup and then start the services
	a.SetHTTPDRoutes()
	a.SetHTTPDMiddleware()
	a.Start()
}
