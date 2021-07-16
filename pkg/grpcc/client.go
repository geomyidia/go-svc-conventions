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
	log.Debugf("Parsed cmd: %s", cmd)
	log.Debugf("Parsed args: %v", args)
}

// Close the gRPC connection
func (c *Client) Close() {
	c.GRPCConn.Close()
	log.Debug("Closed gRPC client connection.")
}

// RunCommand ...
func (c *Client) RunCommand() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Debugf("Running command: %s", c.Command)
	switch c.Command {
	case "echo":
		data := fmt.Sprintf("%s", c.Args)
		r, err := c.GRPCClient.Echo(ctx, &pb.GenericData{Data: data})
		if err != nil {
			log.Fatalf("could not get echo reply: %v", err)
		}
		fmt.Printf("Echo: %s\n", r.GetData())
	case "health":
		r, err := c.GRPCClient.Health(ctx, &pb.HealthRequest{})
		if err != nil {
			log.Fatalf("could not get health reply: %v", err)
		}
		fmt.Printf("Services: %s\n", r.GetServices())
		fmt.Printf("Errors: %s\n", r.GetErrors())
	case "version":
		r, err := c.GRPCClient.Version(ctx, &pb.VersionRequest{})
		if err != nil {
			log.Fatalf("could not get version reply: %v", err)
		}
		fmt.Printf("Version: %s\n", r.GetVersion())
		fmt.Printf("BuildDate: %s\n", r.GetBuildDate())
		fmt.Printf("GitCommit: %s\n", r.GetGitCommit())
		fmt.Printf("GitBranch: %s\n", r.GetGitBranch())
		fmt.Printf("GitSummary: %s\n", r.GetGitSummary())
	default:
		r, err := c.GRPCClient.Ping(ctx, &pb.PingRequest{})
		if err != nil {
			log.Fatalf("could not get ping reply: %v", err)
		}
		fmt.Printf("%s\n", r.GetData())
	}
}
