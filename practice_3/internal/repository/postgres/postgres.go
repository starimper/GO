package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"  
	_ "github.com/golang-migrate/migrate/v4/source/file"       
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"practice_3/pkg/modules"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPGXDialect(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	if err := runMigrations(cfg); err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Dialect{DB: db}
}

func runMigrations(cfg *modules.PostgreConfig) error {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
