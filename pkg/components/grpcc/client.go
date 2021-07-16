package grpcc

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/geomyidia/go-svc-conventions/api"
	"github.com/geomyidia/go-svc-conventions/pkg/components"
)

// Client ...
type Client struct {
	components.Base
	GRPCConn   *grpc.ClientConn
	GRPCClient pb.ServiceExampleClient
	Command    string
	Args       []string
}

// NewClient ...
func NewClient() *Client {
	return &Client{}
}

// SetupConnection ...
func (c *Client) SetupConnection() {
	log.Debug("Setting up client connection ...")
	connectionOpts := c.Config.GRPCD.ConnectionString()
	conn, err := grpc.Dial(connectionOpts, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to gRPC server: %v", err)
	}
	c.GRPCConn = conn
	c.GRPCClient = pb.NewServiceExampleClient(conn)
}

// ParseArgs ...
func (c *Client) ParseArgs(rawArgs []string) {
	cmd := ""
	var args []string
	if len(rawArgs) > 1 {
		cmd = rawArgs[1]
		args = rawArgs[2:]
	}
	c.Command = cmd
	c.Args = args
}

// Close the gRPC connection
func (c *Client) Close() {
	c.GRPCConn.Close()
}

// RunCommand ...
func (c *Client) RunCommand() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch c.Command {
	case "echo":
		data := fmt.Sprintf("%s", c.Args)
		r, err := c.GRPCClient.Echo(ctx, &pb.GenericData{Data: data})
		if err != nil {
			log.Fatalf("could not get echo reply: %v", err)
		}
		log.Printf("Echo: %s", r.GetData())
	case "health":
		r, err := c.GRPCClient.Health(ctx, &pb.HealthRequest{})
		if err != nil {
			log.Fatalf("could not get health reply: %v", err)
		}
		log.Printf("Services: %s", r.GetServices())
		log.Printf("Errors: %s", r.GetErrors())
	case "version":
		r, err := c.GRPCClient.Version(ctx, &pb.VersionRequest{})
		if err != nil {
			log.Fatalf("could not get version reply: %v", err)
		}
		log.Printf("Version: %s", r.GetVersion())
		log.Printf("BuildDate: %s", r.GetBuildDate())
		log.Printf("GitCommit: %s", r.GetGitCommit())
		log.Printf("GitBranch: %s", r.GetGitBranch())
		log.Printf("GitSummary: %s", r.GetGitSummary())
	default:
		r, err := c.GRPCClient.Ping(ctx, &pb.PingRequest{})
		if err != nil {
			log.Fatalf("could not get ping reply: %v", err)
		}
		log.Printf("Reply: %s", r.GetData())
	}
}
