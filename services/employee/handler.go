package employee

import (
	"context"

	"github.com/codebarz/employee-service/entities/employees"
	"github.com/codebarz/employee-service/rpc/proto/employeepb"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
