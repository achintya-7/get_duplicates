package gcp

import (
	"context"
	db "duplicates-finder/db/generated"

	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Client struct {
	Client *storage.Client
	Store  *db.SqlStore
}

func NewClient() (*Client, error) {
	connPool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/upload_manager?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Check if the database is reachable
	err = connPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	store := db.NewSqlStore(connPool)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Client{Client: client, Store: store}, nil
}
