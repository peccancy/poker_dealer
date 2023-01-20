package opts

import (
	"github.com/peccancy/poker_dealer/storage/mongo"
	"time"

	"github.com/peccancy/chassi/log"
)

type Config struct {
	AppName            string
	PrettyLogOutput    bool
	GrpcPort           int
	LogLevel           log.Level
	ServerPort         int
	HealthcheckPort    int
	PrometheusPort     int
	JwtVerifyToken     bool
	JwksURL            string
	SentryURL          string
	SentryENV          string
	JwtSkipVerifyToken bool
	MongoDB            *mongo.Config
}

type ClientConnConfig struct {
	Host        string
	Port        int
	InitTimeout time.Duration
}
