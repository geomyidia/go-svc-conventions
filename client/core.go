package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/geomyidia/go-svc-conventions/components"
	pb "github.com/geomyidia/go-svc-conventions/app/grpc"
	log "github.com/sirupsen/logrus"
)

// Client ...
type Client struct {
	components.Base
	GRPCConn *grpc.ClientConn
	GRPCClient pb.ServiceExampleClient
	Command string
	Args []string
}

// SetupConnection ...
func (c *Client) SetupConnection() {
	connectionOpts := fmt.Sprintf("%s:%d", c.Config.GRPCD.Host, c.Config.GRPCD.Port)
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
	default:
		r, err := c.GRPCClient.Ping(ctx, &pb.PingRequest{})
		if err != nil {
			log.Fatalf("could not get ping reply: %v", err)
		}
		log.Printf("Reply: %s", r.GetData())
	}
}