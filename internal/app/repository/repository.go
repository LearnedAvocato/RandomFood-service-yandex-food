package repository

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

const (
	RequestLogsTableName = "request_logs"
)

type Repository struct {
	client *pgx.Conn
}

func NewRepository(ctx context.Context) (*Repository, error) {
	client, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_DSN"))
	if err != nil {
		return nil, err
	}

	return &Repository{client: client}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.client.Close(ctx)
}
