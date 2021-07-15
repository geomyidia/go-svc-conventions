package httpd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
)

// SetupRoutes ...
func SetupRoutes(cfg *config.Config) *gin.Engine {
	log.Debug("Setting up HTTPD routes ...")
	gin.ForceConsoleColor()
	router := gin.Default()
	router.GET("/echo", Echo)
	router.GET("/health", Health)
	router.GET("/ping", Ping)
	return router
}

// SetupServer ...
func SetupServer(cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:    cfg.HTTPConnectionString(),
		Handler: SetupRoutes(cfg),
	}
}
