package grpcd

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/geomyidia/go-svc-conventions/api"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/db"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
)

// GRPCServer ...
type GRPCServer struct {
	pb.UnimplementedServiceExampleServer
	Addr   string
	Server *grpc.Server
	Bus    *msgbus.MsgBus
	DB     *db.DB
}

// NewGRPCServer ...
func NewGRPCServer(cfg *config.Config, bus *msgbus.MsgBus, db *db.DB) *GRPCServer {
	log.Debug("Setting up gRPC daemon ...")
	s := &GRPCServer{
		Addr: cfg.GRPCD.ConnectionString(),
		Bus:  bus,
		DB:   db,
	}
	gs := grpc.NewServer()
	// Register reflection service on gRPC server.
	reflection.Register(gs)
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
