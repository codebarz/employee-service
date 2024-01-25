package role

import (
	"context"

	"github.com/codebarz/employee-service/entities/roles"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, r *roles.NewRole) (*roles.Role, error) {
	//using uuid because the context doesn't have traceID yet
	traceID := uuid.New().String()
	return s.roleRepo.Create(ctx, traceID, r)
}

func (s *service) Query(ctx context.Context) ([]roles.Role, error) {
	traceID := uuid.New().String()
	return s.roleRepo.Query(ctx, traceID)
}
func (s *service) QueryByID(ctx context.Context, roleID string) (*roles.Role, error) {
	traceID := uuid.New().String()
	return s.roleRepo.QueryByID(ctx, traceID, roleID)
}

func (s *service) Update(ctx context.Context, updatedRole roles.UpdateRole) (*roles.Role, error) {
	traceID := uuid.New().String()
	return s.roleRepo.Update(ctx, traceID, updatedRole.Id, updatedRole)
}

func (s *service) Delete(ctx context.Context, roleID string) error {
	traceID := uuid.New().String()

	return s.roleRepo.Delete(ctx, traceID, roleID)
}
