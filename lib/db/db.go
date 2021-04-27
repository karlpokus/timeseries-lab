package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"timeseries/lib/telemetry"
)

func Connect(url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pgxpool.Connect(ctx, url)
}

func Insert(pool *pgxpool.Pool, rcd telemetry.Record) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // noop if tx is already closed
	stmt := "INSERT INTO telemetry2 (ts, key, value) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, stmt, rcd.Time, rcd.Key, rcd.Value)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
