package employee

import (
	"context"

	"github.com/codebarz/employee-service/entities/employees"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, e *employees.NewEmployee) (*employees.Employee, error) {
	//using uuid because the context doesn't have traceID yet
	traceID := uuid.New().String()
	return s.employeeRepo.Create(ctx, traceID, e)
}

func (s *service) Query(ctx context.Context) ([]employees.Employee, error) {
	traceID := uuid.New().String()
	return s.employeeRepo.Query(ctx, traceID)
}

func (s *service) QueryByID(ctx context.Context, employeeID string) (*employees.Employee, error) {
	traceID := uuid.New().String()
	return s.employeeRepo.QueryByID(ctx, traceID, employeeID)
}

func (s *service) Update(ctx context.Context, updatedEmployee employees.UpdateEmployee) (*employees.Employee, error) {
	traceID := uuid.New().String()
	return s.employeeRepo.Update(ctx, traceID, updatedEmployee.Id, updatedEmployee)
}

func (s *service) Delete(ctx context.Context, employeeID string) error {
	traceID := uuid.New().String()
	return s.employeeRepo.Delete(ctx, traceID, employeeID)
}
