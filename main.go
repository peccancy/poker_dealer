package poker_dealer

import (
	"context"
	"github.com/peccancy/poker_dealer/opts"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	// load config
	config, err := opts.ReadOS()
	if err != nil {
		log.Fatal("service init error:", err)
		os.Exit(1)
	}
	// Configure Sentry
	raven.SetDSN(config.SentryURL)
	raven.SetEnvironment(config.SentryENV)

	// init log
	logging.Init(config.LogLevel, os.Stdout, config.PrettyLogOutput)

	err = service.Start(ctx, config)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
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