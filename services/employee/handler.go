package employee

import (
	"context"

	"github.com/codebarz/employee-service/entities/employees"
	"github.com/codebarz/employee-service/rpc/proto/employeepb"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	l       log.Logger
	service Service
	employeepb.UnimplementedEmployeeServiceServer
}

func NewGRPCHandler(l log.Logger, s Service) *Handler {
	return &Handler{l: l, service: s}
}

func (h *Handler) CreateEmployee(ctx context.Context, eb *employeepb.CreateEmployeeRequest) (*employeepb.Employee, error) {
	newEmployee, err := employees.PBToNewEmployee(eb)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	employee, err := h.service.Create(ctx, newEmployee)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return employees.EmployeeToPB(employee)
}

func (h *Handler) QueryEmployees(ctx context.Context, eb *employeepb.QueryEmployeesRequest) (*employeepb.Employees, error) {
	e, err := h.service.Query(ctx)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return employees.EmployeesToPB(e)
}

func (h *Handler) QueryEmployeeByID(ctx context.Context, eb *employeepb.QueryEmployeeByIDRequest) (*employeepb.Employee, error) {
	e, err := h.service.QueryByID(ctx, eb.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return employees.EmployeeToPB(e)
}

func (h *Handler) DeleteEmployee(ctx context.Context, eb *employeepb.DeleteEmployeeRequest) (*emptypb.Empty, error) {
	err := h.service.Delete(ctx, eb.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil

}

func (h *Handler) UpdateEmployee(ctx context.Context, eb *employeepb.UpdateEmployeeRequest) (*employeepb.Employee, error) {
	firstname := eb.GetFirstName()
	lastname := eb.GetLastName()
	email := eb.GetEmail()
	role := eb.GetRole()

	e, err := h.service.Update(ctx, employees.UpdateEmployee{
		Id:        eb.GetId(),
		FirstName: &firstname,
		LastName:  &lastname,
		Email:     &email,
		Role:      &role,
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return employees.EmployeeToPB(e)
}
