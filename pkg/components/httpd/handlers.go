package httpd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
	"github.com/geomyidia/go-svc-conventions/pkg/version"
)

// HTTPHandlerServer ...
type HTTPHandlerServer struct {
	Addr   string
	Routes *gin.Engine
	Server *http.Server
	Bus    *msgbus.MsgBus
}

func NewHTTPHandlerServer(a *components.Application) *HTTPHandlerServer {
	s := &HTTPHandlerServer{
		Bus: a.Bus,
	}
	cfg := a.Config.HTTPD
	s.SetupRoutes(cfg)
	s.SetupSubscriptions()
	s.Addr = cfg.ConnectionString()
	s.Server = &http.Server{
		Addr:    s.Addr,
		Handler: s.Routes,
	}
	return s
}

// SetupRoutes ...
func (s *HTTPHandlerServer) SetupRoutes(cfg *config.HTTPDConfig) {
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

// SetupSubscriptions ...
func (s *HTTPHandlerServer) SetupSubscriptions() {
	s.Bus.Subscribe("ping", func(event *msgbus.Event) { log.Warnf("Got event: %#v", event) })
}

// Echo ...
func (s *HTTPHandlerServer) Echo(ctx *gin.Context) {
	echoData, _ := ioutil.ReadAll(ctx.Request.Body)
	log.Debugf("Got echo request: %s", echoData)
	ctx.String(http.StatusOK, fmt.Sprintf("%s\n", echoData))
}

// Health ...
func (s *HTTPHandlerServer) Health(ctx *gin.Context) {
	log.Debug("Got health request")
	ctx.String(http.StatusOK, fmt.Sprintf("Services: OK\nErrors: NULL\n"))
}

// Ping ...
func (s *HTTPHandlerServer) Ping(ctx *gin.Context) {
	log.Debug("Got ping request")
	s.Bus.Publish(msgbus.NewEvent("ping", "DATA"))
	ctx.String(http.StatusOK, "pong\n")
}

// Version ...
func (s *HTTPHandlerServer) Version(ctx *gin.Context) {
	log.Debug("Got version request")
	vsn := version.VersionData()
	ctx.String(http.StatusOK, fmt.Sprintf(
		"Version: %s\nBuild Date: %s\nGit Commit: %s\nGit Branch: %s\nGit Summary: %s\n",
		vsn.Semantic, vsn.BuildDate, vsn.GitCommit, vsn.GitBranch, vsn.GitSummary))
}

// Serve ...
func (s *HTTPHandlerServer) Serve() {
	log.Infof("HTTP daemon listening on %s ...", s.Addr)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Info("HTTP daemon is quitting ...")
}

// Shutdown ...
func (s *HTTPHandlerServer) Shutdown(ctx context.Context) {
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Debugf("HTTP Daemon has been shutdown.")
}
