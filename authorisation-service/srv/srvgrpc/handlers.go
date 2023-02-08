package srvgrpc

import (
	"context"

	"github.com/peccancy/authorisation-service/api"
	"github.com/pkg/errors"
)

// GetServiceToken builds service token based on id_token
func (s *Srv) GetServiceToken(_ context.Context, req *api.GetServiceTokenRequest) (*api.GetServiceTokenResponse, error) {
	token, err := s.service.BuildServiceToken(req.IdToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build service token")
	}

	return &api.GetServiceTokenResponse{
		ServiceToken: token,
	}, nil
}
