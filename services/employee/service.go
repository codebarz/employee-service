package employee

import (
	"context"

	"github.com/codebarz/employee-service/entities/employees"
	"github.com/go-kit/log"
)

type service struct {
	l            log.Logger
	employeeRepo employees.Repository
}

type Service interface {
	Create(ctx context.Context, e *employees.NewEmployee) (*employees.Employee, error)
	Query(ctx context.Context) ([]employees.Employee, error)
	QueryByID(ctx context.Context, employeeID string) (*employees.Employee, error)
	Update(ctx context.Context, updatedEmployee employees.UpdateEmployee) (*employees.Employee, error)
	Delete(ctx context.Context, employeeID string) error
}

func NewService(log log.Logger, e employees.Repository) Service {
	return &service{
		l:            log,
		employeeRepo: e,
	}
}
