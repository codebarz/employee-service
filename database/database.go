package database

import (
	"embed"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Config struct {
	PostgresDBURL string
	DisableTLS    bool
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

const migrationVersion = 4

func configureDBURL(cfg Config) (string, error) {
	u, err := url.Parse(cfg.PostgresDBURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing db url")
	}
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func Open(cfg Config) (*sqlx.DB, error) {
	dbURL, err := configureDBURL(cfg)
	if err != nil {
		return nil, err
	}

	return sqlx.Open("postgres", dbURL)
}

// Migrate knows how to migrate the database.
func Migrate(cfg Config) error {
	dbURL, err := configureDBURL(cfg)
	if err != nil {
		return err
	}

	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		return errors.Wrap(err, "migration files")
	}

	if err := m.Migrate(migrationVersion); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migration: no database change")

		} else {
			return errors.Wrap(err, "migration up")
		}
	}

	return nil
}

func Log(query string, args ...interface{}) string {
	for i, arg := range args {
		n := fmt.Sprintf("$%d", i+1)

		var a string
		switch v := arg.(type) {
		case string:
			a = fmt.Sprintf("%q", v)
		case []byte:
			a = string(v)
		case []string:
			a = strings.Join(v, ",")
		default:
			a = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, n, a, 1)
	}

	return query
}
