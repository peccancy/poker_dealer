package srvgql

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/peccancy/authorisation-service/models"
	core "github.com/peccancy/authorisation-service/providers/grpc"
	gqlmodels "github.com/peccancy/authorisation-service/srv/srvgql/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	service Service
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) ServiceToken(_ context.Context, idToken string) (*gqlmodels.ServiceToken, error) {
	token, err := r.service.BuildServiceToken(idToken)
	if err == nil {
		return &gqlmodels.ServiceToken{Token: token}, nil
	}

	log.Error("gqlsrv: failed to build service token:", err)
	return &gqlmodels.ServiceToken{
		Error: &gqlmodels.Error{
			Code:    errCode(err),
			Message: err.Error(),
		},
	}, nil
}

func errCode(err error) gqlmodels.Code {
	if errors.Is(err, core.ErrNotSupportedTokenType) {
		return gqlmodels.CodeInvalidIDToken
	}
	if e := models.ErrNotValidIDToken(""); errors.As(err, &e) {
		return gqlmodels.CodeInvalidIDToken
	}
	if e := models.ErrUserDenied(""); errors.As(err, &e) {
		return gqlmodels.CodeUserDenied
	}
	return gqlmodels.CodeInternal
}
