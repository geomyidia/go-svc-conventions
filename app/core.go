package app

import (
	"fmt"

	"github.com/labstack/echo/v4/middleware"
	"github.com/oubiwann/go-svc-conventions/components"
)

// Application ...
type Application struct {
	components.Default
}

// SetRoutes ...
func (a *Application) SetRoutes() {
	a.HTTP.POST("/echo", Echo)
	a.HTTP.GET("/health", Health)
	a.HTTP.GET("/ping", Ping)
}

// SetMiddleware ...
func (a *Application) SetMiddleware() {
	a.HTTP.Use(middleware.Logger())
	a.HTTP.Use(middleware.Recover())
}

// Start ...
func (a *Application) Start() {
	serverOpts := fmt.Sprintf("%s:%d", a.Config.HTTPD.Host, a.Config.HTTPD.Port)
	server := a.HTTP.Start(serverOpts)
	a.HTTP.Logger.Fatal(server)
}
