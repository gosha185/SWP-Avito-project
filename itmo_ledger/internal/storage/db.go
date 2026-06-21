package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// NewDB creates a new PostgreSQL connection pool with sane defaults.
// Returns an error if the connection cannot be established.
// TODO(post-mvp): Add read replicas support for GET requests.
// TODO(post-mvp): Add connection pool metrics (Prometheus) to monitor WaitCount.
func NewDB(dsn string) (*sql.DB, error) {
	// default values if there is no dsn
	if dsn == "" {
		dsn = "host=localhost port=5432 user=admin password=admin dbname=admindb sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
