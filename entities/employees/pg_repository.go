package employees

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/codebarz/employee-service/database"
	"github.com/codebarz/employee-service/entities/roles"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		return nil, ErrInvalidUUID
	}

	q, args, err := p.psql.Select("*").From("roles").Where(sq.Eq{"id": roleUUID}).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "role.QueryByID", database.Log(q, args...))

	role := &roles.Role{}

	if err := p.db.GetContext(ctx, role, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
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
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "employee.Query", database.Log(q, args...))

	if err := p.db.QueryRowxContext(ctx, q, args...).StructScan(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (p *PgRepository) Query(ctx context.Context, traceID string) ([]Employee, error) {
	// TODO: Populate role data
	// Not sure if should be done here or in the gateway that exposes the service to the client
	// q, args, err := p.psql.Select("*, role as role").From(TABLE_NAME).Join("roles ON employees.role = roles.id").ToSql()
	q, args, err := p.psql.Select("*").From(TABLE_NAME).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "employee.Query", database.Log(q, args...))

	employees := make([]Employee, 0)

	if err := p.db.SelectContext(ctx, &employees, q, args...); err != nil {
		return nil, err
	}

	return employees, nil
}

func (p *PgRepository) QueryByID(ctx context.Context, traceID string, employeeID string) (*Employee, error) {
	// ensure a valid uuid is passed
	eUUID, err := uuid.Parse(employeeID)

	if err != nil {
		return nil, ErrInvalidUUID
	}

	q, args, err := p.psql.Select("*").From(TABLE_NAME).Where(sq.Eq{"id": eUUID}).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "employee.QueryByID", database.Log(q, args...))

	employee := &Employee{}

	if err := p.db.GetContext(ctx, employee, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return employee, nil
}

func (p *PgRepository) Update(ctx context.Context, traceID string, employeeID string, updatedEmployee UpdateEmployee) (*Employee, error) {
	if _, err := uuid.Parse(employeeID); err != nil {
		return nil, ErrInvalidUUID
	}

	e, err := p.QueryByID(ctx, traceID, employeeID)

	if err != nil {
		return nil, ErrNotFound
	}

	if updatedEmployee.FirstName != nil {
		e.FirstName = *updatedEmployee.FirstName
	}

	if updatedEmployee.LastName != nil {
		e.LastName = *updatedEmployee.LastName
	}

	if updatedEmployee.Email != nil {
		e.Email = *updatedEmployee.Email
	}

	if updatedEmployee.Role != nil {
		e.Role = *updatedEmployee.Role
	}

	q, args, err := p.psql.Update("Employees").SetMap(map[string]interface{}{
		"first_name": e.FirstName,
		"last_name":  e.LastName,
		"email":      e.Email,
		"role":       e.Role,
	}).Where(sq.Eq{"id": employeeID}).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	if _, err = p.db.ExecContext(ctx, q, args...); err != nil {
		return nil, err
	}

	return e, nil
}

func (p *PgRepository) Delete(ctx context.Context, traceID string, employeeID string) error {
	eUUID, err := uuid.Parse(employeeID)

	if err != nil {
		return ErrInvalidUUID
	}

	q, args, err := p.psql.Delete("employees").Where(sq.Eq{"id": eUUID}).ToSql()

	if err != nil {
		return ErrBuildingQuery
	}

	level.Info(p.log).Log("%s : %s : query : %s", traceID, "employee.Delete",
		database.Log(q, args...),
	)

	if _, err = p.db.ExecContext(ctx, q, args...); err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "22P02" {
				return ErrBuildingQuery
			}
		}
		return err
	}

	return nil
}
