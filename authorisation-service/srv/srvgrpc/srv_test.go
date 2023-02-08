package srvgrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Srv

	srv.readiness = false
	assert.Equal(t, "grpc service is't ready yet", srv.HealthCheck().Error())
}
