package bo

import (
	"context"

	core "github.com/peccancy/authorisation-service/providers/grpc"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/peccancy/authorisation-service/models"

	log "github.com/sirupsen/logrus"
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

	c.client = api.NewBackOfficeServiceClient(c.conn)

	return &c, nil
}

type Client struct {
	// bo srv
	conn   *grpc.ClientConn
	client api.BackOfficeServiceClient
	// context
	ctx context.Context
}

func (s *Client) Check() error {
	if s.conn == nil {
		return core.ErrInit("back office service client")
	}

	if s.conn.GetState() != connectivity.Ready {
		return core.ErrConn("back office service client")
	}

	return nil
}

func (s *Client) GetTenantMachineTokenPermissions(accountID, clientID string) (models.Permissions, error) {
	resp, err := s.client.GetTenantMachineTokenPermissions(s.ctx, &api.GetTenantMachineTokenPermissionsRequest{
		AccountId: accountID,
		ClientId:  clientID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "bo_client.read_tenant_machine_permissions")
	}
	return resp.Permissions, nil
}

func (s *Client) GetInternalMachineTokenPermissions(clientID string) (models.Permissions, error) {
	resp, err := s.client.GetInternalMachineTokenPermissions(s.ctx, &api.GetInternalMachineTokenPermissionsRequest{
		ClientId: clientID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "bo_client.read_internal_machine_permissions")
	}
	return resp.Permissions, nil
}
