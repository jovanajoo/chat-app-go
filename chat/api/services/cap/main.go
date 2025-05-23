package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/jovanajoo/chat-app-go/chat/foundation/logger"
)

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return "" // TODO: NEED TRACE IDs
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "CAP", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------

	log.Info(ctx, "startup", "status", "started")
	defer log.Info(ctx, "startup", "status", "shutting down")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil
}
