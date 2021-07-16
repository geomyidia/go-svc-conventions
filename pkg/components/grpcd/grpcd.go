package grpcd

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/geomyidia/go-svc-conventions/api"
	"github.com/geomyidia/go-svc-conventions/pkg/components"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
)

// GRPCServer ...
type GRPCServer struct {
	pb.UnimplementedServiceExampleServer
	Addr   string
	Server *grpc.Server
	Bus    *msgbus.MsgBus
}

// NewGRPCServer ...
func NewGRPCServer(app *components.Application) *GRPCServer {
	log.Debug("Setting up gRPC daemon ...")
	cfg := app.Config.GRPCD
	s := &GRPCServer{
		Addr: cfg.ConnectionString(),
		Bus:  app.Bus,
	}
	gs := grpc.NewServer()
	s.RegisterServer(gs)
	s.Server = gs
	log.Debug("gRPC implementation set up.")
	return s
}

// RegisterServer ...
func (s *GRPCServer) RegisterServer(grpcServer *grpc.Server) {
	pb.RegisterServiceExampleServer(grpcServer, s)
}

// Serve
func (s *GRPCServer) Serve() {
	log.Infof("gRPC daemon listening on %s ...", s.Addr)

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Server.Serve(lis)
	log.Info("gRPC daemon is quitting ...")
}

// Shutdown ...
func (s *GRPCServer) Shutdown() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if s.Server != nil {
			s.Server.GracefulStop()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if s.Server != nil {
			s.Server.Stop()
		}
	}()

	wg.Wait()
	log.Debugf("gRPC Daemon has been shutdown.")
}
