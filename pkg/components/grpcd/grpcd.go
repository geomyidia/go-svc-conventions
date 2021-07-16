package grpcd

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
)

// SetupServer ...
func SetupServer(cfg *config.Config) *GRPCHandlerServer {
	log.Debug("Setting up gRPC daemon ...")
	s := NewGRPCHandlerServer(cfg.GRPCD)
	log.Debug("gRPC implementation set up.")
	return s
}
