package gcp

import (
	"context"
	db "duplicates-finder/db/generated"
	"log"
	"runtime"

	"cloud.google.com/go/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Client struct {
	Client *storage.Client
	Store  *db.SqlStore
}

func NewClient() (*Client, error) {
	config, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/upload_manager?sslmode=disable")
	if err != nil {
		return nil, err
	}

	numWorkers := runtime.NumCPU()
	if numWorkers > 20 {
		log.Println("Number of workers is too high, setting to 20")
		config.MaxConns = 20
	} else {
		log.Println("Number of workers is", numWorkers)
		config.MaxConns = int32(numWorkers)
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), config)
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
