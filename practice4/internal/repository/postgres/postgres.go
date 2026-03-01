package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"practice_4/pkg/modules"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPGXDialect(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var db *sqlx.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		log.Println("DB not ready yet, retrying...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		panic(err)
	}

	return &Dialect{DB: db}
}
