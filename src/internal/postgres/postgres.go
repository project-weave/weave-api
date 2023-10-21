package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	db  *sql.DB
	DSN string
}

func NewDB(dsn string) *DB {
	return &DB{
		DSN: dsn,
	}
}

func (db *DB) Open() error {
	if db.DSN == "" {
		return fmt.Errorf("DSN is empty, could not open db pool")
	}

	postgresDB, err := sql.Open("postgres", db.DSN)
	if err != nil {
		return err
	}
	postgresDB.SetMaxOpenConns(25)
	postgresDB.SetMaxIdleConns(25)
	postgresDB.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = postgresDB.PingContext(ctx)
	if err != nil {
		return err
	}

	db.db = postgresDB
	return nil
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
