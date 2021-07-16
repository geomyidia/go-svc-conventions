package grpcd

import (
	"context"

	log "github.com/sirupsen/logrus"

	pb "github.com/geomyidia/go-svc-conventions/api"
	"github.com/geomyidia/go-svc-conventions/pkg/version"
)

// Echo ...
func (s *GRPCServer) Echo(ctx context.Context, in *pb.GenericData) (*pb.GenericData, error) {
	log.Debugf("Received echo %v", in)
	return &pb.GenericData{Data: in.GetData()}, nil
}

// Health ...
func (s *GRPCServer) Health(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
	log.Debugf("Received health request")
	return &pb.HealthReply{Services: "OK", Errors: "NULL"}, nil
}

// Ping ...
func (s *GRPCServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Debug("Received ping")
	return &pb.PingReply{Data: "PONG"}, nil
}

// Version ...
func (s *GRPCServer) Version(
	_ context.Context, in *pb.VersionRequest) (*pb.VersionReply, error) {
	log.Debugf("Received version request")
	vsn := version.VersionData()
	return &pb.VersionReply{
		Version:    vsn.Semantic,
		BuildDate:  vsn.BuildDate,
		GitCommit:  vsn.GitCommit,
		GitBranch:  vsn.GitBranch,
		GitSummary: vsn.GitSummary,
	}, nil
}
