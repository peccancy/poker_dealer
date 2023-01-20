package service

import (
	"context"
	"github.com/peccancy/chassi/grpc/auth"
	"github.com/peccancy/chassi/health"
	"github.com/peccancy/chassi/prometheus"
	"github.com/peccancy/poker_dealer/logs"
	"github.com/peccancy/poker_dealer/opts"
	"github.com/peccancy/poker_dealer/storage/mongo"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"sync"

	chassisgrpc "github.com/peccancy/chassi/grpc"
)

type tocken struct{}

func Start(ctx context.Context, config *opts.Config) (err error) {
	logs.Info().Write("Starting service")

	repo, err := mongo.New(ctx, config.MongoDB)
	if err != nil {
		logs.StartService().Write(err)
		return errors.Wrap(err, "start mongo connection error")
	}
	defer repo.Close()

	//t := tocken{}
	//
	//rout := router.New(repo, t)

	var wg sync.WaitGroup

	// run main gRPC service, default on 8080
	_, err = chassisgrpc.ConnectServiceSimpleWithContext(
		ctx,
		&wg,
		config.GrpcPort,
		func(server *grpc.Server) {
			//api.RegisterTicketsDomainServiceServer(server, rout)
		},
		config.JwksURL,
		!config.JwtSkipVerifyToken,
		logs.EntryFactory,
		createRoutePermissionMap(),
	)
	if err != nil {
		logs.GrpcServerInitFailure().Write(err)
		return errors.Wrap(err, "start grpc server error")
	}

	// run healthcheck service, default on 8888
	checks := []func() error{repo.HealthCheck}

	health.ConnectHealthCheckServiceWithCtx(ctx, &wg, config.HealthcheckPort, checks, checks)

	// run prometheus service, default on 9100
	prometheus.ConnectPrometheusServiceWithContext(ctx, &wg, config.PrometheusPort)

	logs.Info().Write("Service started")

	wg.Wait()

	return nil
}

func createRoutePermissionMap() auth.RoutePermissionMap {
	return auth.NewRoutePermissionMap() //.
	//WithRoutePermission(api.TicketsDomainService_GetTicketTypes, resources.Tickets, authorisation.All, authorisation.Read).
	//WithRoutePermission(api.TicketsDomainService_GetTicketByID, resources.Tickets, authorisation.Any, authorisation.Read).
	//WithRoutePermission(api.TicketsDomainService_GetTickets, resources.Tickets, authorisation.All, authorisation.Read).
	//WithRoutePermission(api.TicketsDomainService_GetTicketsByLinkedEntities, resources.Tickets, authorisation.All, authorisation.Read).
	//WithRoutePermission(api.TicketsDomainService_CreateTicket, resources.Tickets, authorisation.Any, authorisation.Create).
	//WithRoutePermission(api.TicketsDomainService_UpdateTicketByOwner, resources.Tickets, authorisation.Any, authorisation.Update).
	//WithRoutePermission(api.TicketsDomainService_UpdateTicketByAssignee, resources.Tickets, authorisation.Any, authorisation.Update)
}
