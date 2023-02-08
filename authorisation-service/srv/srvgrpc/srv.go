package srvgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/peccancy/authorisation-service/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func New(port int, service Service) (*Srv, error) {
	// build Service
	srv := Srv{
		service: service,
		addr:    fmt.Sprintf(":%d", port),
		server: grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_prometheus.UnaryServerInterceptor)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpc_prometheus.StreamServerInterceptor)),
		),
	}
	api.RegisterAuthorisationServiceServer(srv.server, &srv)

	return &srv, nil
}

type Srv struct {
	addr      string
	service   Service
	server    *grpc.Server
	runErr    error
	readiness bool
}

func (s *Srv) Run(ctx context.Context, wg *sync.WaitGroup) {
	log.Info("grpc service: begin run")

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.runErr = err
		log.Error("grpc service: run error: ", err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.server.Serve(lis)
		log.Error("grpc service: end run > ", err)
		s.runErr = err
	}()

	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
		log.Info("grpc service: graceful stop")
	}()

	s.readiness = true
}

func (s *Srv) HealthCheck() error {
	if !s.readiness {
		return errors.New("grpc service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("grpc service: run issue")
	}
	return nil
}
