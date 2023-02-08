package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"

	"github.com/peccancy/authorisation-service/providers/grpc/bo"
	"github.com/peccancy/authorisation-service/providers/grpc/users"

	"github.com/peccancy/authorisation-service/srv/srvgql"

	"github.com/peccancy/authorisation-service/srv/srvgrpc"
	"github.com/peccancy/authorisation-service/srv/srvhttp"

	"github.com/peccancy/authorisation-service/options"
	"github.com/peccancy/authorisation-service/service"
	"github.com/peccancy/authorisation-service/servicetoken"
	"github.com/peccancy/authorisation-service/srv/healthcheck"
	"github.com/peccancy/authorisation-service/srv/prometheus"
	_ "github.com/peccancy/chassi/gracefulshutdown"
)

func main() {
	// read grpcsrv config from os env
	config := options.ReadEnv()

	// init logger
	initLogger(config)

	log.Info("begin...")

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	if err := run(ctx, config); err != nil {
		cancel()
		log.WithError(err).Error("service startup")
	}
}

func run(ctx context.Context, config options.Config) error {
	// load private key
	keyData, _ := pem.Decode(config.PrivateKeyBytes)
	if keyData == nil {
		return errors.New("private key pem decode empty result")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyData.Bytes)
	if err != nil {
		return errors.Wrap(err, "private key parse")
	}

	// build bo srv client
	boClient, err := bo.New(ctx, config.MDServiceConf)
	if err != nil {
		return errors.Wrap(err, "back-office service connection")
	}

	// build users srv client
	usersClient, err := users.New(ctx, config.UserServiceConf)
	if err != nil {
		return errors.Wrap(err, "user service connection")
	}

	// exchanger to exchange from id tokent to service token
	exchanger := servicetoken.NewExchanger(boClient, usersClient, config.Issuer, config.DefaultTokenTTL,
		&privateKey.PublicKey, config.SkipUpsertUser)

	// build main service
	srv := service.New(exchanger, config.Issuer, privateKey, config.KeyID, config.DefaultTokenTTL)

	// build http service
	httpSrv, err := srvhttp.New(config.HTTPPort, srv)
	if err != nil {
		return errors.Wrap(err, "http service init error")
	}

	// build graphql service
	gqlSrv, err := srvgql.New(config.GraphqlPort, srv)
	if err != nil {
		return errors.Wrap(err, "graphql service init error")
	}

	// build grpc service
	grpcSrv, err := srvgrpc.New(config.GRPCPort, srv)
	if err != nil {
		return errors.Wrap(err, "grpc service init error")
	}

	// build prometheus service
	prometheusSrv := prometheus.New(config.PrometheusPort)

	// build healthcheck service
	healthSrv := healthcheck.New(
		config.HealthCheckPort,
		prometheusSrv.HealthCheck,
		httpSrv.HealthCheck,
		gqlSrv.HealthCheck,
		grpcSrv.HealthCheck,
		boClient.Check,
		usersClient.Check,
	)

	var wg = &sync.WaitGroup{}

	// run srv
	httpSrv.Run(ctx, wg)
	grpcSrv.Run(ctx, wg)
	gqlSrv.Run(ctx, wg)
	healthSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)

	// wait while working
	wg.Wait()
	log.Info("end")

	return nil
}

func initLogger(config options.Config) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	switch strings.ToLower(config.LogLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
