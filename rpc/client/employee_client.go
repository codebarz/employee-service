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

// EmployeeClient is the event API gRPC client
type EmployeeClient struct {
	l    log.Logger
	conn *grpc.ClientConn
	c    employeepb.EmployeeServiceClient
}

// New creates a new gRPC client
func NewEmployeeServiceClient(l log.Logger, conn *grpc.ClientConn) *EmployeeClient {
	if l == nil {
		l = log.NewNopLogger()
	}
	c := &EmployeeClient{
		l:    l,
		conn: conn,
		c:    employeepb.NewEmployeeServiceClient(conn),
	}
	depstatus.Register(grpcdep.Wrap("employeeservice", conn))
	return c
}

// Close closes the connection of the gRPC client
func (c *EmployeeClient) Close() {
	c.conn.Close()
}

// CreateEvent validated the token and returns the login id contained
func (c *EmployeeClient) CreateEmployee(ctx context.Context, in *employeepb.CreateEmployeeRequest) (*employeepb.Employee, error) {
	resp, err := c.c.CreateEmployee(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *EmployeeClient) QueryEmployee(ctx context.Context, in *employeepb.QueryEmployeesRequest) (*employeepb.Employees, error) {
	resp, err := c.c.QueryEmployees(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *EmployeeClient) QueryEmployeeByID(ctx context.Context, in *employeepb.QueryEmployeeByIDRequest) (*employeepb.Employee, error) {
	resp, err := c.c.QueryEmployeeByID(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *EmployeeClient) UpdateEmployee(ctx context.Context, in *employeepb.UpdateEmployeeRequest) (*employeepb.Employee, error) {
	resp, err := c.c.UpdateEmployee(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *EmployeeClient) DeleteEmployee(ctx context.Context, in *employeepb.DeleteEmployeeRequest) (*emptypb.Empty, error) {
	resp, err := c.c.DeleteEmployee(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
