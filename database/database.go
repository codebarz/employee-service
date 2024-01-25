package database

import (
	"embed"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/johejo/golang-migrate-extra/source/iofs"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseURL string
	l           *log.Logger
}

var fs embed.FS

func (c *Config) OpenConnection(cfg Config) (*sqlx.DB, error) {
	if cfg.DatabaseURL == "" {
		c.l.Fatal("invalid postgres db URL passed")
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
		c.l.Fatal("Migration error:", err)
		return err
	}

	if err := m.Migrate(1); err != nil {
		c.l.Fatal("Migration error:", err)
		return err
	}

	return nil
}

func NewDatabase(l *log.Logger) *Config {
	return &Config{l: l}
}
