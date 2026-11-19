package config

import (
	"github.com/DeleteElf/goframework/entities"
)

type ServerConfig struct {
	entities.BaseConfig
	Port  int
	Token string
}
