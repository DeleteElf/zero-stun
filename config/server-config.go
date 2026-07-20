package config

import (
	"github.com/deleteelf/goframework/entities"
)

type ServerConfig struct {
	entities.BaseConfig
	Port int
}
