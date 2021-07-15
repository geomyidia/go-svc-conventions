package grpcd

import (
	"github.com/geomyidia/reverb"
	"github.com/labstack/gommon/log"
)

// SetupgRPCImplementation ...
func (a *Application) SetupgRPCImplementation(r *reverb.Reverb) {
	log.Debug("Setting up gRPC implementation ...")
	s := NewGRPCHandlerServer()
	s.RegisterServer(r.GRPCServer)
	log.Info("gRPC implementation set up.")
}

// StartgRPCD ...
func (a *Application) StartgRPCD() {
	log.Debug("Starting gRPC daemon ...")
	serverOpts := a.Config.GRPCConnectionString()
	server := a.GRPCD.Start(serverOpts)
	a.SetupgRPCImplementation(server)
	go server.Serve()
	log.Infof("gRPC daemon started on %s.", serverOpts)
}
