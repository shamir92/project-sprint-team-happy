package database

import (
	"database/sql"
	"fmt"
	"halosuster/configuration"
	"log"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type postgresWriter struct {
	db *sql.DB
}

type IPostgresWriter interface {
	GetDB() *sql.DB
	Close() error
}

func NewPostgresWriter(configDB configuration.IDatabaseWriter) (*postgresWriter, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		strings.TrimSpace(configDB.GetUser()),
		strings.TrimSpace(configDB.GetPassword()),
		strings.TrimSpace(configDB.GetHost()),
		strings.TrimSpace(configDB.GetPort()),
		strings.TrimSpace(configDB.GetName()),
		strings.TrimSpace(configDB.GetDBParam()),
	)
	log.Println(dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// log.Fatal(err) // TODO: handle this properly!
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &postgresWriter{db: db}, nil
}

func (mw *postgresWriter) GetDB() *sql.DB {
	return mw.db
}

func (mw *postgresWriter) Close() error {
	return mw.db.Close()
}
