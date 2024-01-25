package employees

import (
	"context"
	"time"
)

type Employee struct {
	ID        string     `json:"id" db:"id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Role      string     `json:"role" db:"role"`
	Email     string     `json:"email" db:"email"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type NewEmployee struct {
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
	Role      string `json:"role" db:"role" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required"`
}

type UpdateEmployee struct {
	Id        string  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      *string `json:"role"`
	Email     *string `json:"email"`
}

type Repository interface {
	Create(ctx context.Context, traceID string, e *NewEmployee) (*Employee, error)
	Query(ctx context.Context, traceID string) ([]Employee, error)
	QueryByID(ctx context.Context, traceID string, employeeID string) (*Employee, error)
	Update(ctx context.Context, traceID string, employeeID string, updatedRole UpdateEmployee) (*Employee, error)
	Delete(ctx context.Context, traceID string, employeeID string) error
}
