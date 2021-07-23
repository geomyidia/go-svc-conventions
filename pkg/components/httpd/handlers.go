package httpd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
	"github.com/geomyidia/go-svc-conventions/pkg/version"
)

// Echo ...
func (s *HTTPServer) Echo(ctx *gin.Context) {
	echoData, _ := ioutil.ReadAll(ctx.Request.Body)
	log.Debugf("Received HTTP echo request: %s", echoData)
	ctx.String(http.StatusOK, fmt.Sprintf("%s\n", echoData))
}

// Health ...
func (s *HTTPServer) Health(ctx *gin.Context) {
	log.Debug("Received HTTP health request")
	event := msgbus.NewEvent("status:health", "DATA")
	s.Bus.Publish(event)
	ctx.String(http.StatusOK, fmt.Sprintf("Services: OK\nErrors: NULL\n"))
}

// Ping ...
func (s *HTTPServer) Ping(ctx *gin.Context) {
	log.Debug("Received HTTP ping request")
	log.Tracef("Available topics: %+v", s.Bus.Topics())
	event := msgbus.NewEvent("status:ping", "DATA")
	s.Bus.Publish(event)
	ctx.String(http.StatusOK, "pong\n")
}

// Version ...
func (s *HTTPServer) Version(ctx *gin.Context) {
	log.Debug("Received HTTP version request")
	vsn := version.VersionData()
	ctx.String(http.StatusOK, fmt.Sprintf(
		"Version: %s\nBuild Date: %s\nGit Commit: %s\nGit Branch: %s\nGit Summary: %s\n",
		vsn.Semantic, vsn.BuildDate, vsn.GitCommit, vsn.GitBranch, vsn.GitSummary))
}
