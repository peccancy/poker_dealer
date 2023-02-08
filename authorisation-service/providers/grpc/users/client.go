package users

import (
	"context"
	grpc_prometheus "pkg/mod/github.com/grpc-ecosystem/go-grpc-prometheus@v1.2.0"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	"github.com/peccancy/authorisation-service/models"
	core "github.com/peccancy/authorisation-service/providers/grpc"
	api "github.com/peccancy/user/api/unauthorized"
)

func New(ctx context.Context, cfg core.Config) (*Client, error) {
	c := Client{
		ctx: ctx,
	}

	// init connection to multitenancy-dashboard-domain-grpcsrv
	var err error
	c.conn, err = grpc.DialContext(
		ctx,
		cfg.Target,
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithInsecure(),
		//grpc.WithBlock(),
	)
	if err != nil {
		log.Info("providers.grpc.bo-srv init err:", err)
		return nil, errors.Wrap(err, "failed to init providers.grpc.bo-srv grpc connection")
	}

	// handle shutdown
	go func() {
		<-ctx.Done()
		if c.conn != nil {
			log.Info("providers.grpc.bo-srv: close connection")
			c.conn.Close()
		}
	}()

	c.client = api.NewUnauthorizedUserServiceClient(c.conn)

	return &c, nil
}

type Client struct {
	// bo srv
	conn   *grpc.ClientConn
	client api.UnauthorizedUserServiceClient
	// context
	ctx context.Context
}

func (s *Client) Check() error {
	if s.conn == nil {
		return core.ErrInit("users service client")
	}

	if s.conn.GetState() != connectivity.Ready {
		return core.ErrConn("users service client")
	}

	return nil
}

func (s *Client) Name() string {
	return "Users Service Client"
}

func (s *Client) GetUserPermissions(accountID, userID, email string) (models.Permissions, string, error) {
	resp, err := s.client.GetUserPermissionsAndID(s.ctx, &api.GetUserPermissionsAndIDRequest{
		AccountId: accountID,
		UserId:    userID,
		Email:     email,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return nil, "", models.ErrUserDenied("user_service_client.user_denied")
		}

		return nil, "", errors.Wrap(err, "user_service_client.read_user_permissions")
	}
	return resp.Permissions, resp.Id, nil
}

func (s *Client) Upsert(accountID, userID, givenName, familyName, email, source string) error {
	req := api.UpsertUserRequest{
		UserId:     userID,
		AccountId:  accountID,
		GivenName:  givenName,
		FamilyName: familyName,
		Email:      email,
		Source:     source,
	}

	resp, err := s.client.UpsertUser(s.ctx, &req)
	if err != nil {
		return errors.Wrap(err, "user_service_client.upsert_user")
	}

	switch resp.Status {
	case api.UpsertUserResponse_ACCEPTED:
		return nil
	case api.UpsertUserResponse_IGNORED:
		log.WithFields(log.Fields{
			"accountID": accountID,
			"userID":    userID,
		}).Debug("user_service_client.user_ignored")
		return nil
	case api.UpsertUserResponse_USER_DENIED:
		return models.ErrUserDenied("user_service_client.user_denied")
	default:
		return errors.New("user_service_client.unsupported_resp_status")
	}
}

func (s *Client) GetTokenContext(accountID, userID, email string) (string, error) {
	resp, err := s.client.GetUserContext(s.ctx, &api.GetUserContextRequest{
		AccountId: accountID,
		UserId:    userID,
		Email:     email,
	})
	if err != nil {
		return "", errors.Wrap(err, "user_service_client.read_user_context")
	}

	switch resp.GetStatus() {
	case api.GetUserContextResponse_OK:
		return resp.Context, nil
	case api.GetUserContextResponse_NOT_FOUND:
		return "", nil
	default:
		return "", errors.New("user_service_client.invalid_response")
	}
}
