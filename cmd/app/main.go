package main

import (
	"github.com/labstack/echo/v4"

	"github.com/oubiwann/go-svc-conventions/app"
	"github.com/oubiwann/go-svc-conventions/cfg"
)

func main() {
	a := new(app.Application)
	a.Config = cfg.NewConfig()
	a.HTTP = echo.New()

	a.SetRoutes()
	a.SetMiddleware()
	a.Start()
}
