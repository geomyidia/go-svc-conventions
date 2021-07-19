package httpd

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/db"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
)

// HTTPServer ...
type HTTPServer struct {
	Addr   string
	Routes *gin.Engine
	Server *http.Server
	Bus    *msgbus.MsgBus
	DB     *db.DB
}

func NewHTTPServer(cfg *config.Config, bus *msgbus.MsgBus, db *db.DB) *HTTPServer {
	log.Debug("Setting up HTTP daemon ...")
	s := &HTTPServer{
		Bus: bus,
		DB:  db,
	}
	s.SetupRoutes(cfg.HTTPD)
	s.Addr = cfg.HTTPD.ConnectionString()
	s.Server = &http.Server{
		Addr:    s.Addr,
		Handler: s.Routes,
	}
	log.Debug("HTTP daemon set up.")
	return s
}

// SetupRoutes ...
func (s *HTTPServer) SetupRoutes(cfg *config.HTTPDConfig) {
	log.Debug("Setting up HTTPD routes ...")
	var router *gin.Engine
	if cfg.RequestLogging {
		gin.ForceConsoleColor()
		router = gin.Default()
	} else {
		router = gin.New()
	}
	router.POST("/echo", s.Echo)
	router.GET("/health", s.Health)
	router.GET("/ping", s.Ping)
	router.GET("/version", s.Version)
	s.Routes = router
}

// Serve ...
func (s *HTTPServer) Serve() {
	log.Infof("HTTP daemon listening on %s ...", s.Addr)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Info("HTTP daemon is quitting ...")
}

// Shutdown ...
func (s *HTTPServer) Shutdown(ctx context.Context) {
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Debugf("HTTP Daemon has been shutdown.")
}
