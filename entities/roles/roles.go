package roles

import (
	"context"
	"time"
)

type Role struct {
	ID        string     `json:"id" db:"id"`
	Title     string     `json:"title" validate:"required" db:"title"`
	Level     int32      `json:"level" db:"level"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type NewRole struct {
	Title string `json:"title" validate:"required"`
	Level int    `json:"level"`
}

type UpdateRole struct {
	Id    string  `json:"id"`
	Title *string `json:"title"`
	Level *int32  `json:"level"`
}

type Repository interface {
	Create(ctx context.Context, traceID string, r *NewRole) (*Role, error)
	Query(ctx context.Context, traceID string) ([]Role, error)
	QueryByID(ctx context.Context, traceID string, roleID string) (*Role, error)
	Update(ctx context.Context, traceID string, roleID string, updatedRole UpdateRole) (*Role, error)
	Delete(ctx context.Context, traceID string, roleID string) error
}
