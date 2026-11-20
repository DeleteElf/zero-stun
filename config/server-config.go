package config

import (
	"github.com/DeleteElf/zero-core/entities"
)

type ServerConfig struct {
	entities.BaseConfig
	Port  int
	Token string
}
