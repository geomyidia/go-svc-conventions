package app

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/geomyidia/go-svc-conventions/api"
	log "github.com/sirupsen/logrus"
)

// ExampleServerImpl ...
type ExampleServerImpl struct {
	pb.UnimplementedServiceExampleServer
}

// New ...
func NewExampleServer() *ExampleServerImpl {
	return &ExampleServerImpl{}
}

// Echo ...
func (s *ExampleServerImpl) Echo(ctx context.Context, in *pb.GenericData) (*pb.GenericData, error) {
	log.Debugf("Received: %v", in)
	return &pb.GenericData{Data: in.GetData()}, nil
}

// Health ...
func (s *ExampleServerImpl) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
	log.Debugf("Received: %v", in)
	return &pb.HealthReply{Services: "OK", Errors: "NULL"}, nil
}

// Ping ...
func (s *ExampleServerImpl) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Debugf("Received: %v", in)
	return &pb.PingReply{Data: "PONG"}, nil
}

// RegisterServer ...
func (s *ExampleServerImpl) RegisterServer(grpcServer *grpc.Server) {
	pb.RegisterServiceExampleServer(grpcServer, s)
}
