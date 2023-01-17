package service

import (
	"github.com/pkg/errors"
	"context"
	"fmt"
	"github.com/peccancy/chassi/authorisation"
	"github.com/peccancy/chassi/health"
	"github.com/peccancy/chassi/http/auth"
	"github.com/peccancy/chassi/prometheus"
	"net/http"
	"poker_dealer/logs"
	"poker_dealer/opts"
	"runtime/debug"
	"sync"

	chassishttp "github.com/peccancy/chassi/http"
)

func Start(ctx context.Context, config *opts.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("panic %s: %s", r, string(debug.Stack()))
			logs.ServicePanicked().Write(msg)
			err = errors.Errorf(msg)
		}
	}()

	// Configure service
	logs.Info().Write("Starting service")

	jwkProvider := authorisation.NewURLJwkProvider(config.JwksURL)


	var wg sync.WaitGroup

	chassishttp.ConnectServiceWithContext(ctx, &wg, config.ServerPort, handler)

	prometheus.ConnectPrometheusServiceWithContext(ctx, &wg, config.PrometheusPort)

	var checks []func() error
	go func(ctx context.Context) {
		logs.Info().Write("Starting health check")
		src := health.ConnectHealthCheckServiceWithCtx(ctx, &wg, config.HealthcheckPort, checks, checks)
		<-ctx.Done()
		err := src.Close()
		if err != nil {
			logs.ServicePanicked.WithCorrelationID(ctx).Write(err)
		}
	}(ctx)

	logs.Info().Write("Service started")

	wg.Wait()

	return nil
}
