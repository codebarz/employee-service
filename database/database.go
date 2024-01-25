package database

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kit/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseURL string
	l           log.Logger
}

//go:embed migrations/*.sql
var fs embed.FS

func (c *Config) OpenConnection(cfg Config) (*sqlx.DB, error) {
	if cfg.DatabaseURL == "" {
		c.l.Log("invalid postgres db URL passed")
		return nil, errors.New("invalid postgres db URL passed")
	}

	return sqlx.Open("postgres", cfg.DatabaseURL)
}

func (c *Config) Migrate(cfg Config) error {
	d, err := iofs.New(fs, "migrations")

	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, cfg.DatabaseURL)

	if err != nil {
		c.l.Log("Migration error:", err)
		return err
	}

	if err := m.Migrate(1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			c.l.Log("migration: no database change", err)
		} else {
			return err
		}
	}

	return nil
}

func NewDatabase(l log.Logger) *Config {
	return &Config{l: l}
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
