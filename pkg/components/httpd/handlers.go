package httpd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Echo ...
func Echo(ctx *gin.Context) {
	echoed, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.String(http.StatusOK, fmt.Sprintf("%s\n", echoed))
}

// Health ...
func Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, fmt.Sprintf("Services: OK\nErrors: NULL\n"))
}

// Ping ...
func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong\n")
}
