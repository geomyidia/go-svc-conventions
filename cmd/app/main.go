package main

import (
	"context"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/internal/util"
	"github.com/geomyidia/go-svc-conventions/pkg/components"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/grpcd"
	"github.com/geomyidia/go-svc-conventions/pkg/components/httpd"
	"github.com/geomyidia/go-svc-conventions/pkg/components/logging"
)

// Application ...
type Application struct {
	components.Default
}

func main() {
	// Create the application object and assign components to it
	a := new(Application)
	a.Config = config.NewConfig()
	a.Logger = logging.Load(a.Config)

	// Create context that listens for the interrupt signal from the OS.
	ctx, cancel := util.SignalWithContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup

	httpDaemon := httpd.SetupServer(a.Config)
	grpcDaemon := grpcd.SetupServer(a.Config)

	// Initialise the HTTP server in its own goroutine and wire to wait group
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpDaemon.Serve()
	}()

	// Initialise the gRPC server in its own goroutine and wire to wait group
	// XXX this has been dsiable since it blocks (due to gRPC server-shutdown
	//     not using context / cancelation) Is there are way to cancel gRPC
	//     servers with a context?
	//wg.Add(1)
	go func() {
		//defer wg.Done()
		grpcDaemon.Serve()
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	cancel()
	log.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The HTTP server requires a context to shutdown
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	httpDaemon.Shutdown(ctx)
	grpcDaemon.Shutdown()

	log.Info("Waiting for wait groups to finish ...")
	wg.Wait()
	log.Info("Application shutdown complete.")
}
