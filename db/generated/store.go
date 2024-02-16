package db

import "github.com/jackc/pgx/v5/pgxpool"

type SqlStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewSqlStore(connPool *pgxpool.Pool) *SqlStore {
	return &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
