package grpcd

import (
	"context"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/geomyidia/go-svc-conventions/api"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/version"
)

// GRPCHandlerServer ...
type GRPCHandlerServer struct {
	pb.UnimplementedServiceExampleServer
	Addr   string
	Server *grpc.Server
}

// NewGRPCHandlerServer ...
func NewGRPCHandlerServer(cfg *config.GRPCDConfig) *GRPCHandlerServer {
	s := &GRPCHandlerServer{Addr: cfg.ConnectionString()}
	r := grpc.NewServer()
	s.RegisterServer(r)
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

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Server.Serve(lis)
	log.Info("gRPC daemon is quitting ...")
}

// Shutdown ...
func (s *GRPCHandlerServer) Shutdown() {
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
