package configuration

import (
	"os"
)

type appConfiguration struct {
	port string
}

func NewAppConfiguration() *appConfiguration {
	return &appConfiguration{
		port: os.Getenv("APP_PORT"),
	}
}

type IAppConfiguration interface {
	GetPort() string
}

func (ac *appConfiguration) GetPort() string {
	return ac.port
}
