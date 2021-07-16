package grpcd

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/geomyidia/reverb"

	pb "github.com/geomyidia/go-svc-conventions/api"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/version"
)

// GRPCHandlerServer ...
type GRPCHandlerServer struct {
	pb.UnimplementedServiceExampleServer
	Addr   string
	Server *reverb.Reverb
}

// NewGRPCHandlerServer ...
func NewGRPCHandlerServer(cfg *config.GRPCDConfig) *GRPCHandlerServer {
	s := &GRPCHandlerServer{Addr: cfg.ConnectionString()}
	r := reverb.New()
	s.RegisterServer(r.GRPCServer)
	s.Server = r
	return s
}

// RegisterServer ...
func (s *GRPCHandlerServer) RegisterServer(grpcServer *grpc.Server) {
	pb.RegisterServiceExampleServer(grpcServer, s)
}

// Echo ...
func (s *GRPCHandlerServer) Echo(ctx context.Context, in *pb.GenericData) (*pb.GenericData, error) {
	log.Debugf("Received: %v", in)
	return &pb.GenericData{Data: in.GetData()}, nil
}

// Health ...
func (s *GRPCHandlerServer) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
	log.Debugf("Received: %v", in)
	return &pb.HealthReply{Services: "OK", Errors: "NULL"}, nil
}

// Ping ...
func (s *GRPCHandlerServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Debugf("Received: %v", in)
	return &pb.PingReply{Data: "PONG"}, nil
}

// Version ...
func (s *GRPCHandlerServer) Version(
	_ context.Context, in *pb.VersionRequest) (*pb.VersionReply, error) {
	log.Debugf("Received: %v", in)
	vsn := version.VersionData()
	return &pb.VersionReply{
		Version:    vsn.Semantic,
		BuildDate:  vsn.BuildDate,
		GitCommit:  vsn.GitCommit,
		GitBranch:  vsn.GitBranch,
		GitSummary: vsn.GitSummary,
	}, nil
}

// Serve
func (s *GRPCHandlerServer) Serve() {
	log.Infof("gRPC daemon listening on %s ...", s.Addr)
	s.Server.Start(s.Addr).Serve()
	log.Info("gRPC daemon is quitting ...")
}

// Shutdown ...
func (s *GRPCHandlerServer) Shutdown() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if s.Server != nil {
			s.Server.GRPCServer.GracefulStop()
			//s.Server.GRPCServer.Stop()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if s.Server != nil {
			s.Server.Close()
		}
	}()

	wg.Wait()
	log.Debugf("gRPC Daemon has been shutdown.")
}
