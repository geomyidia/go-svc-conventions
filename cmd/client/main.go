package main

import (
	"os"

	"github.com/geomyidia/go-svc-conventions/cfg"
	"github.com/geomyidia/go-svc-conventions/client"
	logger "github.com/geomyidia/go-svc-conventions/components/logging"
)

func main() {
	// Create the client objects and assign components to it
	c := new(client.Client)
	c.Config = cfg.NewConfig()
	c.Logger = logger.Load(c.Config)

	// Perform client setup and then issue the parsed command
	c.SetupConnection()
	defer c.Close()
	c.ParseArgs(os.Args)
	c.RunCommand()
}
