package options

import (
	"time"

	grpcCore "github.com/peccancy/authorisation-service/providers/grpc"
)

type Config struct {
	LogLevel        string
	HTTPPort        int
	GraphqlPort     int
	GRPCPort        int
	HealthCheckPort int
	PrometheusPort  int
	Issuer          string
	KeyID           string
	PrivateKeyBytes []byte
	DefaultTokenTTL time.Duration
	MDServiceConf   grpcCore.Config
	UserServiceConf grpcCore.Config
	SkipUpsertUser  bool
}
