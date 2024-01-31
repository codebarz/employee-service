package client

import (
	"context"

	"github.com/codebarz/employee-service/rpc/proto/employeepb"
	"github.com/go-kit/log"
	"github.com/greyfinance/grey-go-libs/monitoring/depstatus"
	"github.com/greyfinance/grey-go-libs/monitoring/depstatus/grpcdep"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Client is the event API gRPC client
type Client struct {
	l    log.Logger
	conn *grpc.ClientConn
	c    employeepb.EmployeeServiceClient
}

// New creates a new gRPC client
func NewEmployeeServiceClient(l log.Logger, conn *grpc.ClientConn) *Client {
	if l == nil {
		l = log.NewNopLogger()
	}
	c := &Client{
		l:    l,
		conn: conn,
		c:    employeepb.NewEmployeeServiceClient(conn),
	}
	depstatus.Register(grpcdep.Wrap("employeeservice", conn))
	return c
}

// Close closes the connection of the gRPC client
func (c *Client) Close() {
	c.conn.Close()
}

// CreateEvent validated the token and returns the login id contained
func (c *Client) CreateEmployee(ctx context.Context, in *employeepb.CreateEmployeeRequest) (*employeepb.Employee, error) {
	resp, err := c.c.CreateEmployee(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) QueryEmployee(ctx context.Context, in *employeepb.QueryEmployeesRequest) (*employeepb.Employees, error) {
	resp, err := c.c.QueryEmployees(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) QueryEmployeeByID(ctx context.Context, in *employeepb.QueryEmployeeByIDRequest) (*employeepb.Employee, error) {
	resp, err := c.c.QueryEmployeeByID(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) UpdateEmployee(ctx context.Context, in *employeepb.UpdateEmployeeRequest) (*employeepb.Employee, error) {
	resp, err := c.c.UpdateEmployee(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteEmployee(ctx context.Context, in *employeepb.DeleteEmployeeRequest) (*emptypb.Empty, error) {
	resp, err := c.c.DeleteEmployee(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
