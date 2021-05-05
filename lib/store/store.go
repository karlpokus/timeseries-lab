package store

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"timeseries/lib/telemetry"
)

type Store interface {
	Insert(telemetry.Record) error
	Keys() ([]string, error)
	Query(telemetry.Query) ([]telemetry.Record, error)
	Close()
}

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) Insert(rcd telemetry.Record) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // noop if tx is already closed
	stmt := "INSERT INTO telemetry3 (ts, key, value) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, stmt, rcd.Time, rcd.Key, rcd.Value)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (db *Database) Keys() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.Pool.Query(ctx, "SELECT DISTINCT key FROM telemetry3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (db *Database) Query(q telemetry.Query) ([]telemetry.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// TODO fix this mess
	stmt := "SELECT * FROM telemetry3 WHERE ts >= $1 AND ts < $2"
	target := q.Targets[0].Target
	if target != "" {
		stmt += " AND key = '" + target + "'" // sql injection vuln?
	}
	stmt += " LIMIT $3"

	log.Println(stmt)

	rows, err := db.Pool.Query(ctx, stmt, q.From, q.To, q.MaxDataPoints)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rcds []telemetry.Record
	for rows.Next() {
		var (
			Time  time.Time
			Key   string
			Value float64
		)
		err = rows.Scan(&Time, &Key, &Value)
		if err != nil {
			return nil, err
		}
		rcds = append(rcds, telemetry.Record{
			Time, Key, Value,
		})
	}
	return rcds, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}

func New(url string) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Database{
		Pool: pool,
	}, nil
}
