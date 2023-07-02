package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	RequestLogsTableName = "request_logs"
)

type Repository struct {
	client *pgxpool.Pool
}

func NewRepository(ctx context.Context) (*Repository, error) {

	connConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_DSN"))
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}

	connPool, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Repository{client: connPool}, nil
}

func (r *Repository) Close() {
	r.client.Close()
}
