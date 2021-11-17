package ent

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3" //
)

// SQLLiteOpen Open new SQLLite connection
func SQLLiteOpen(databaseURL string) *Client {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	maxConnections := 1
	if limit := os.Getenv("SQLLITE_CONN_LIMIT"); limit != "" {
		connLimit, err := strconv.Atoi(limit)
		if err != nil {
			return nil
		} else {
			maxConnections = connLimit
		}
	}

	maxIdle := 10
	if limit := os.Getenv("SQLLITE_IDLE_LIMIT"); limit != "" {
		idleLimit, err := strconv.Atoi(limit)
		if err != nil {
			return nil
		} else {
			maxConnections = idleLimit
		}
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxConnections)
	// db.SetConnMaxLifetime(time.Hour)

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.SQLite, db)
	return NewClient(Driver(drv))
}
