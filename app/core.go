package app

import (
	"fmt"

	"github.com/labstack/echo/v4/middleware"
	"github.com/geomyidia/go-svc-conventions/components"
	"github.com/geomyidia/reverb"
	log "github.com/sirupsen/logrus"
)

// Application ...
type Application struct {
	components.Default
}

// SetHTTPDRoutes ...
func (a *Application) SetHTTPDRoutes() {
	log.Debug("Setting up HTTPD routes ...")
	a.HTTPD.POST("/rest/echo", Echo)
	a.HTTPD.GET("/rest/health", Health)
	a.HTTPD.GET("/rest/ping", Ping)
	log.Info("HTTPD routes set up.")
}

// SetHTTPDMiddleware ...
func (a *Application) SetHTTPDMiddleware() {
	log.Debug("Setting up HTTPD middleware ...")
	a.HTTPD.Pre(middleware.RemoveTrailingSlash())
	a.HTTPD.Use(middleware.Logger())
	a.HTTPD.Use(middleware.Recover())
	log.Info("HTTPD middleware set up.")
}

// SetupgRPCImplementation ...
func (a *Application) SetupgRPCImplementation(r *reverb.Reverb) {
	log.Debug("Setting up gRPC implementation ...")
	s := NewExampleServer()
	s.RegisterServer(r.GRPCServer)
	log.Info("gRPC implementation set up.")
}

// StartgRPCD ...
func (a *Application) StartgRPCD() {
	log.Debug("Starting gRPC daemon ...")
	serverOpts := fmt.Sprintf("%s:%d", a.Config.GRPCD.Host, a.Config.GRPCD.Port)
	server := a.GRPCD.Start(serverOpts)
	a.SetupgRPCImplementation(server)
	go server.Serve()
	log.Infof("gRPC daemon started on %s.", serverOpts)
}

// StartHTTPD ...
func (a *Application) StartHTTPD() {
	log.Debug("Starting HTTP daemon ...")
	serverOpts := fmt.Sprintf("%s:%d", a.Config.HTTPD.Host, a.Config.HTTPD.Port)
	server := a.HTTPD.Start(serverOpts)
	a.HTTPD.Logger.Fatal(server)
}

// Start ...
func (a *Application) Start() {
	a.StartgRPCD()
	a.StartHTTPD()
	log.Info("System started.")
}
