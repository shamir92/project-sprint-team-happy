package postgres

import (
	"database/sql"
	"log"
)

func NewDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	return db
}
