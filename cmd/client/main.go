package main

import (
	"os"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/grpcc"
	"github.com/geomyidia/go-svc-conventions/pkg/components/logging"
)

func main() {
	// Create the client object and assign components to it
	c := grpcc.NewClient()
	c.Config = config.NewConfig()
	c.Logger = logging.LoadClient(c.Config)

	// Perform client setup and then issue the parsed command
	c.SetupConnection()
	defer c.Close()
	c.ParseArgs(os.Args)
	c.RunCommand()
}
