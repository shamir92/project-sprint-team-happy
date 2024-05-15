package configuration

import (
	"os"
)

type databaseWriter struct {
	user     string
	password string
	host     string
	port     string
	name     string
	param    string
}

type IDatabaseWriter interface {
	GetUser() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetName() string
	GetDBParam() string
}

func NewDatabaseWriter() *databaseWriter {
	return &databaseWriter{
		user:     os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		name:     os.Getenv("DB_NAME"),
		param:    os.Getenv("DB_PARAMS"),
	}
}

func (dw *databaseWriter) GetUser() string {
	return dw.user
}

func (dw *databaseWriter) GetPassword() string {
	return dw.password
}

func (dw *databaseWriter) GetHost() string {
	return dw.host
}

func (dw *databaseWriter) GetPort() string {
	return dw.port
}

func (dw *databaseWriter) GetName() string {
	return dw.name
}

func (dw *databaseWriter) GetDBParam() string {
	return dw.param
}
