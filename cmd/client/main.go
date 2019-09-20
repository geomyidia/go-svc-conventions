package main

import (
	"os"

	"github.com/geomyidia/go-svc-conventions/client"
	"github.com/geomyidia/go-svc-conventions/components/config"
	"github.com/geomyidia/go-svc-conventions/components/logging"
)

func main() {
	// Create the client object and assign components to it
	c := new(client.Client)
	c.Config = config.NewConfig()
	c.Logger = logging.Load(c.Config)

	// Perform client setup and then issue the parsed command
	c.SetupConnection()
	defer c.Close()
	c.ParseArgs(os.Args)
	c.RunCommand()
}
