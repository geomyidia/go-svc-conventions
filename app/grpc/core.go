package grpc

import (
	"context"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

// ExampleServerImpl ...
type ExampleServerImpl struct {}

// New ...
func New() *ExampleServerImpl {
	return &ExampleServerImpl{}
}

// Echo ...
func (s *ExampleServerImpl) Echo(ctx context.Context, in *GenericData) (*GenericData, error) {
	log.Debugf("Received: %v", in)
	return &GenericData{Data: in.GetData()}, nil
}

// Health ...
func (s *ExampleServerImpl) Health(ctx context.Context, in *HealthRequest) (*HealthReply, error) {
	log.Debugf("Received: %v", in)
	return &HealthReply{Services: "OK", Errors: "NULL"}, nil
}

// Ping ...
func (s *ExampleServerImpl) Ping(ctx context.Context, in *PingRequest) (*PingReply, error) {
	log.Debugf("Received: %v", in)
	return &PingReply{Data: "PONG"}, nil
}

// RegisterServer ...
func (s *ExampleServerImpl) RegisterServer(grpcServer *grpc.Server) {
	RegisterServiceExampleServer(grpcServer, s)

}
