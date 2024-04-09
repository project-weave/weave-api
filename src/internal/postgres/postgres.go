package postgres

import (
	"context"

	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	dsn  string
}

func NewDB(dsn string) *DB {
	return &DB{
		dsn: dsn,
	}
}

func (db *DB) Open() error {
	if db.dsn == "" {
		return fmt.Errorf("DSN is empty, could not open db pool")
	}

	pgConfig, err := pgxpool.ParseConfig(db.dsn)
	if err != nil {
		return fmt.Errorf("could not parse DSN: %w", err)
	}
	pgConfig.MaxConns = 25
	pgConfig.MaxConnIdleTime = 15 * time.Minute

	pgpool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pgpool.Ping(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Connected to database")

	db.pool = pgpool
	return nil
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
