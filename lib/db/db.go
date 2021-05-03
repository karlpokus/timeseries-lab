package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"timeseries/lib/telemetry"
)

type Range struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Target struct {
	Data   interface{} `json:"data"` // empty string or {key: "", operator: "=", value: x}
	Target string      `json:"target"`
	Type   string      `json:"type"`
}

type QueryRequest struct {
	Range         `json:"range"`
	Targets       []Target `json:"targets"`
	AdhocFilters  []string `json:"adhocFilters"`
	MaxDataPoints int      `json:"maxDataPoints"`
}

type SearchRequest struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

// Connect returns a connection pool to the db
func Connect(url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pgxpool.Connect(ctx, url)
}

// Insert inserts the record in the db
func Insert(pool *pgxpool.Pool, rcd telemetry.Record) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := pool.Begin(ctx)
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

// Keys return unique keys from db
func Keys(pool *pgxpool.Pool) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := pool.Query(ctx, "SELECT DISTINCT key FROM telemetry3")
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

func Query(pool *pgxpool.Pool, q QueryRequest) ([]telemetry.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := "SELECT * FROM telemetry3 WHERE ts >= $1 AND ts < $2"
	target := q.Targets[0].Target
	if target != "" {
		stmt += " AND key = '" + target + "'" // sql injection vuln?
	}
	stmt += " LIMIT $3"

	log.Println(stmt)

	rows, err := pool.Query(ctx, stmt, q.From, q.To, q.MaxDataPoints)
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

/*
	store interface

	import
		lib/store
		lib/telemetry

	type Store struct {
		Insert(telemetry.Record) error
		Keys() ([]string, error)
		Query(QueryRequest) ([]telemetry.Record, error)
	}

	// returns a Store
	str := store.New()
	str := store.Mock()

*/
