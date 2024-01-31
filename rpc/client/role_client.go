package client

import (
	"context"

	"github.com/codebarz/employee-service/rpc/proto/rolepb"
	"github.com/go-kit/log"
	"github.com/greyfinance/grey-go-libs/monitoring/depstatus"
	"github.com/greyfinance/grey-go-libs/monitoring/depstatus/grpcdep"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoleClient struct {
	l    log.Logger
	conn *grpc.ClientConn
	c    rolepb.RoleServiceClient
}

func NewRoleClient(l log.Logger, conn *grpc.ClientConn) *RoleClient {
	if l == nil {
		l = log.NewNopLogger()
	}
	c := &RoleClient{
		l:    l,
		conn: conn,
		c:    rolepb.NewRoleServiceClient(conn),
	}
	depstatus.Register(grpcdep.Wrap("employeeservice", conn))
	return c
}

// Close closes the connection of the gRPC client
func (c *RoleClient) Close() {
	c.conn.Close()
}

func (c *RoleClient) CreateRole(ctx context.Context, in *rolepb.CreateRoleRequest) (*rolepb.Role, error) {
	resp, err := c.c.CreateRole(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RoleClient) QueryRole(ctx context.Context, in *rolepb.QueryRoleRequest) (*rolepb.Roles, error) {
	resp, err := c.c.QueryRoles(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RoleClient) QueryRoleByID(ctx context.Context, in *rolepb.QueryRoleByIDRequest) (*rolepb.Role, error) {
	resp, err := c.c.QueryRoleByID(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RoleClient) UpdateRole(ctx context.Context, in *rolepb.UpdateRoleRequest) (*rolepb.Role, error) {
	resp, err := c.c.UpdateRole(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *RoleClient) DeleteRole(ctx context.Context, in *rolepb.DeleteRoleRequest) (*emptypb.Empty, error) {
	resp, err := c.c.DeleteRole(ctx, in)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
