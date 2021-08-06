package main

import (
	"context"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/internal/util"
	"github.com/geomyidia/go-svc-conventions/pkg/components"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/db"
	"github.com/geomyidia/go-svc-conventions/pkg/components/grpcd"
	"github.com/geomyidia/go-svc-conventions/pkg/components/httpd"
	"github.com/geomyidia/go-svc-conventions/pkg/components/logging"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
)

func main() {
	// Create the application object and assign components to it
	app := new(components.Application)
	app.Config = config.NewConfig()
	app.Logger = logging.Load(app.Config)
	app.Bus = msgbus.NewMsgBus()
	app.DB = db.NewDB(app.Config, app.Bus)
	app.HTTPD = httpd.NewHTTPServer(app.Config, app.Bus, app.DB)
	app.GRPCD = grpcd.NewGRPCServer(app.Config, app.Bus, app.DB)

	// Set up subscriptions
	var handlers []msgbus.Handler
	handlers = append(handlers,
		msgbus.AddHandler("*", msgbus.HandleWildCard),
		//msgbus.AddHandler("status:*", msgbus.HandleStatusWildCard),
		msgbus.AddHandler("status:ping", msgbus.HandlePing),
		msgbus.AddHandler("status:health", msgbus.HandleHealth),
	)
	app.Bus.AddHandlers(handlers)

	// Create context that listens for the interrupt signal from the OS.
	ctx, cancel := util.SignalWithContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup

	// Initialise the message bus/event auditor in its own goroutine and wire to wait group
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Bus.Serve(ctx)
	}()

	// Initialise the database connection in its own goroutine and wire to wait group
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.DB.Connect()
	}()

	// Initialise the HTTP server in its own goroutine and wire to wait group
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.HTTPD.Serve()
	}()

	// Initialise the gRPC server in its own goroutine and wire to wait group
	// XXX this has been dsiable since it blocks (due to gRPC server-shutdown
	//     not using context / cancelation) Is there are way to cancel gRPC
	//     servers with a context?
	//     ticket: https://github.com/geomyidia/go-svc-conventions/issues/11
	//wg.Add(1)
	go func() {
		//defer wg.Done()
		app.GRPCD.Serve()
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	cancel()
	log.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The HTTP server requires a context to shutdown
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Shutdown running components
	app.HTTPD.Shutdown(ctx)
	app.GRPCD.Shutdown()
	app.DB.Shutdown()

	log.Info("Waiting for wait groups to finish ...")
	wg.Wait()
	log.Info("Application shutdown complete.")
}
