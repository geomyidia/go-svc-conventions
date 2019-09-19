package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/geomyidia/go-svc-conventions/app"
	pb "github.com/geomyidia/go-svc-conventions/app/grpc"
	"github.com/geomyidia/go-svc-conventions/cfg"
	logger "github.com/geomyidia/go-svc-conventions/components/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	a := new(app.Application)
	a.Config = cfg.NewConfig()
	a.Logger = logger.Load()
	connectionOpts := fmt.Sprintf("%s:%d", a.Config.GRPCD.Host, a.Config.GRPCD.Port)

	// Set up a connection to the server.
	conn, err := grpc.Dial(connectionOpts, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to gRPC server: %v", err)
	}
	defer conn.Close()
	c := pb.NewServiceExampleClient(conn)

	// Contact the server and print out its response.
	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch command {
	case "echo":
		data := fmt.Sprintf("%s", os.Args[2:])
		r, err := c.Echo(ctx, &pb.GenericData{Data: data})
		if err != nil {
			log.Fatalf("could not get echo reply: %v", err)
		}
		log.Printf("Echo: %s", r.GetData())
	case "health":
		r, err := c.Health(ctx, &pb.HealthRequest{})
		if err != nil {
			log.Fatalf("could not get health reply: %v", err)
		}
		log.Printf("Services: %s", r.GetServices())
		log.Printf("Errors: %s", r.GetErrors())
	default:
		r, err := c.Ping(ctx, &pb.PingRequest{})
		if err != nil {
			log.Fatalf("could not get ping reply: %v", err)
		}
		log.Printf("Reply: %s", r.GetData())
	}
	
}
