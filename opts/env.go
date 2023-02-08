package opts

import (
	"github.com/peccancy/chassi/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ReadOS() (*Config, error) {

	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("APP_NAME", "pocker_dealer")
	viper.SetDefault("PRETTY_LOG_OUTPUT", false)
	viper.SetDefault("LOG_LEVEL", "ERROR")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("JWKS_URL", "http://authorisation-service/.well-known/jwks.json")
	viper.SetDefault("JWKS_EXPIRATION", "1m")
	viper.SetDefault("JWT_VERIFY_TOKEN", false)


	logLevel, err := log.ParseLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}

	return &Config{
		AppName:         viper.GetString("APP_NAME"),
		PrettyLogOutput: viper.GetBool("PRETTY_LOG_OUTPUT"),
		LogLevel:        *logLevel,
		ServerPort:      viper.GetInt("SERVER_PORT"),
		HealthcheckPort: viper.GetInt("HEALTHCHECK_PORT"),
		PrometheusPort:  viper.GetInt("PROMETHEUS_PORT"),
		JwksURL:         viper.GetString("JWKS_URL"),
		JwtVerifyToken:  viper.GetBool("JWT_VERIFY_TOKEN"),
		SentryURL:       viper.GetString("SENTRY_URL"),
		SentryENV:       viper.GetString("SENTRY_ENV"),
	}, nil
}
