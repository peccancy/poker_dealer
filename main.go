package poker_dealer

import (
	"context"
	"github.com/peccancy/poker_dealer/logs"
	"github.com/peccancy/poker_dealer/opts"
	logging "github.com/peccancy/chassi/log"
	"log"
	"os"
	"os/signal"
	"github.com/peccancy/poker_dealer/service"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if logs.ServicePanickedError != nil {
				logs.ServicePanickedError().Write("%s: %s", r, string(debug.Stack()))
			}
			log.Fatalf("panic: %s: %s", r, string(debug.Stack()))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	// load config
	config, err := opts.ReadOS()
	if err != nil {
		log.Fatal("service init error:", err)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// init log
	logging.Init(config.LogLevel, os.Stdout, config.PrettyLogOutput)

	err = service.Start(ctx, config)
	if err != nil {
		logs.Error().Write(err)
		os.Exit(1)
	}

}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("Got Interrupt signal")
		stop()
	}()
}