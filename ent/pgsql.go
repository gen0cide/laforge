package ent

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib" //
)

// PGOpen Open new PostGres connection
func PGOpen(databaseURL string) *Client {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	maxConnections := 10
	if limit := os.Getenv("PG_CONN_LIMIT"); limit != "" {
		connLimit, err := strconv.Atoi(limit)
		if err != nil {
			return nil
		} else {
			maxConnections = connLimit
		}
	}

	maxIdle := 10
	if limit := os.Getenv("PG_IDLE_LIMIT"); limit != "" {
		idleLimit, err := strconv.Atoi(limit)
		if err != nil {
			return nil
		} else {
			maxConnections = idleLimit
		}
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxConnections)
	db.SetConnMaxLifetime(time.Hour)

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return NewClient(Driver(drv))
}
