package main

import (
	"context"
	"net/http"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/app"
	"github.com/geomyidia/go-svc-conventions/internal/util"
	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/httpd"
	"github.com/geomyidia/go-svc-conventions/pkg/components/logging"
)

func main() {
	// Create the application object and assign components to it
	a := new(app.Application)
	a.Config = config.NewConfig()
	a.Logger = logging.Load(a.Config)
	//a.GRPCD = reverb.New()

	// Create context that listens for the interrupt signal from the OS.
	ctx, cancel := util.SignalWithContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var wg sync.WaitGroup

	httpDaemon := httpd.SetupServer(a.Config)

	// Initializing the HTTP server in its own goroutine and wire to wait group
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Infof("HTTP daemon listening on %s ...", a.Config.HTTPConnectionString())
		if err := httpDaemon.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	cancel()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpDaemon.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")

	wg.Wait()
}
