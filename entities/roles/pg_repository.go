package roles

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/codebarz/employee-service/database"
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

const TABLE_NAME = "roles"

func NewPgRepository(l log.Logger, db *sqlx.DB) Repository {
	return &PgRepository{
		log:  l,
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (p *PgRepository) Create(ctx context.Context, traceID string, r *NewRole) (*Role, error) {
	role := &Role{}

	q, args, err := p.psql.Insert(TABLE_NAME).SetMap(map[string]interface{}{
		"title": r.Title,
		"level": r.Level,
	}).Suffix("RETURNING *").ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "role.Query", database.Log(q, args...))

	if err := p.db.QueryRowxContext(ctx, q, args...).StructScan(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (p *PgRepository) Query(ctx context.Context, traceID string) ([]Role, error) {
	q, args, err := p.psql.Select("*").From(TABLE_NAME).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "role.Query", database.Log(q, args...))

	role := make([]Role, 0)

	if err := p.db.SelectContext(ctx, &role, q, args...); err != nil {
		return nil, err
	}

	return role, nil
}

func (p *PgRepository) QueryByID(ctx context.Context, traceID string, roleID string) (*Role, error) {
	// ensure a valid uuid is passed
	roleUUID, err := uuid.Parse(roleID)

	if err != nil {
		return nil, ErrInvalidUUID
	}

	q, args, err := p.psql.Select("*").From(TABLE_NAME).Where(sq.Eq{"id": roleUUID}).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	level.Info(p.log).Log("traceID", traceID, "role.QueryByID", database.Log(q, args...))

	role := &Role{}

	if err := p.db.GetContext(ctx, role, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return role, nil
}

func (p *PgRepository) Update(ctx context.Context, traceID string, roleID string, updatedRole UpdateRole) (*Role, error) {
	if _, err := uuid.Parse(roleID); err != nil {
		return nil, ErrInvalidUUID
	}

	role, err := p.QueryByID(ctx, traceID, roleID)

	if err != nil {
		return nil, err
	}

	if updatedRole.Title != nil {
		role.Title = *updatedRole.Title
	}

	if updatedRole.Level != nil {
		role.Level = *updatedRole.Level
	}

	q, args, err := p.psql.Update("Roles").SetMap(map[string]interface{}{
		"title": role.Title,
		"level": role.Level,
	}).Where(sq.Eq{"id": roleID}).ToSql()

	if err != nil {
		return nil, ErrBuildingQuery
	}

	if _, err = p.db.ExecContext(ctx, q, args...); err != nil {
		return nil, err
	}

	return role, nil
}

func (p *PgRepository) Delete(ctx context.Context, traceID string, roleID string) error {
	id, err := uuid.Parse(roleID)

	if err != nil {
		return ErrBuildingQuery
	}

	q, args, err := p.psql.Delete("roles").Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return ErrBuildingQuery
	}

	level.Info(p.log).Log("%s : %s : query : %s", traceID, "role.Delete",
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
