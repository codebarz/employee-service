package employees

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/codebarz/employee-service/database"
	"github.com/codebarz/employee-service/entities/roles"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PgRepository struct {
	log  log.Logger
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

const TABLE_NAME = "employees"

func NewPgRepository(l log.Logger, db *sqlx.DB) Repository {
	return &PgRepository{
		log:  l,
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (p *PgRepository) Create(ctx context.Context, traceID string, e *NewEmployee) (*Employee, error) {
	employee := &Employee{}

	roleUUID, err := uuid.Parse(e.Role)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error parsing role uuid get query by id %v", err))
	}

	q, args, err := p.psql.Select("*").From("roles").Where(sq.Eq{"id": roleUUID}).ToSql()

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error building get role query by id %v", err))
	}

	level.Info(p.log).Log("traceID", traceID, "role.QueryByID", database.Log(q, args...))

	role := &roles.Role{}

	if err := p.db.GetContext(ctx, role, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(fmt.Sprintf("error getting role by id %v", err))
		}
		return nil, err
	}

	q, args, err = p.psql.Insert(TABLE_NAME).SetMap(map[string]interface{}{
		"first_name": e.FirstName,
		"last_name":  e.LastName,
		"role":       role.ID,
		"email":      e.Email,
	}).Suffix("RETURNING *").ToSql()

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error building insert query %v", err))
	}

	level.Info(p.log).Log("traceID", traceID, "employee.Query", database.Log(q, args...))

	if err := p.db.QueryRowxContext(ctx, q, args...).StructScan(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (p *PgRepository) Query(ctx context.Context, traceID string) ([]Employee, error) {
	return nil, nil
}

func (p *PgRepository) QueryByID(ctx context.Context, traceID string, employeeID string) (*Employee, error) {
	return nil, nil
}

func (p *PgRepository) Update(ctx context.Context, traceID string, employeeID string, updatedRole UpdateEmployee) (*Employee, error) {
	return nil, nil
}

func (p *PgRepository) Delete(ctx context.Context, traceID string, employeeID string) error {
	return nil
}
