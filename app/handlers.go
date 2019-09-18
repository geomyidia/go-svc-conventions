package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthData ...
type HealthData struct {
	Services string `json:"services"`
	Errors   string `json:"errors"`
}

// Echo ...
func Echo(ctx echo.Context) (err error) {
	echoed, _ := ioutil.ReadAll(ctx.Request().Body)
	return ctx.String(http.StatusOK, fmt.Sprintf("%s", echoed))
}

// Health ...
func Health(ctx echo.Context) (err error) {
	h := &HealthData{
		Services: "OK",
		Errors:   "NULL",
	}
	return ctx.JSON(http.StatusOK, h)
}

// Ping ...
func Ping(ctx echo.Context) (err error) {
	return ctx.String(http.StatusOK, "pong")
}
