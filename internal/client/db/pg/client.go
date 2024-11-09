package pg

import (
	"context"
	"log"

	"github.com/vbulash/chat-server/internal/client/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgClient struct {
	masterDBC db.DB
}

// New ...
func New(ctx context.Context, dsn string) (db.Client, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Ошибка конфигурации pgxpool: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Ошибка коннекта к БД: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: pool},
	}, nil
}

// DB ...
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close ...
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
