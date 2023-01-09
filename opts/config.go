package opts

import (
	"time"

	"github.com/peccancy/chassi/log"
)

type Config struct {
	AppName            string
	PrettyLogOutput    bool
	LogLevel           log.Level
	ServerPort         int
	HealthcheckPort    int
	PrometheusPort     int
	JwtVerifyToken     bool
	JwksURL            string
	TicketConnConf     *ClientConnConfig //TODO it's a NP for integrate some local service
	SentryURL          string
	SentryENV          string
}

type ClientConnConfig struct {
	Host        string
	Port        int
	InitTimeout time.Duration
}
