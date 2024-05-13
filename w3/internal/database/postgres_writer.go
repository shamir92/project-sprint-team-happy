package database

import (
	"database/sql"
	"fmt"
	"halosuster/configuration"
	"log"
	"strings"
)

type postgresWriter struct {
	db *sql.DB
}

type IPostgresWriter interface {
	Close() error
}

func NewPostgresWriter(configDB configuration.IDatabaseWriter) *postgresWriter {
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?%s",
		strings.TrimSpace(configDB.GetUser()),
		strings.TrimSpace(configDB.GetPassword()),
		strings.TrimSpace(configDB.GetProtocol()),
		strings.TrimSpace(configDB.GetHost()),
		strings.TrimSpace(configDB.GetPort()),
		strings.TrimSpace(configDB.GetName()),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err) // TODO: handle this properly!
		return nil
	}
	return &postgresWriter{db: db}
}

func (mw *postgresWriter) GetDB() *sql.DB {
	return mw.db
}

func (mw *postgresWriter) Close() error {
	return mw.db.Close()
}
