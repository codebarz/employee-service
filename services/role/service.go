package role

import (
	"context"

	"github.com/codebarz/employee-service/entities/roles"
	"github.com/go-kit/log"
)

type service struct {
	l        log.Logger
	roleRepo roles.Repository
}

type Service interface {
	Create(ctx context.Context, r *roles.NewRole) (*roles.Role, error)
	Query(ctx context.Context) ([]roles.Role, error)
	QueryByID(ctx context.Context, roleID string) (*roles.Role, error)
	Update(ctx context.Context, updatedRole roles.UpdateRole) (*roles.Role, error)
	Delete(ctx context.Context, roleID string) error
}

func NewService(log log.Logger, r roles.Repository) Service {
	return &service{
		l:        log,
		roleRepo: r,
	}
}
