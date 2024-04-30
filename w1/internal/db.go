package internal

import "github.com/jmoiron/sqlx"

var db *sqlx.DB

// GetDB returns the global database connection instance
func GetDB() *sqlx.DB {
	return db
}

func SetDB(dbConn *sqlx.DB) {
	db = dbConn
}
