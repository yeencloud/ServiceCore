package config

import (
	"github.com/yeencloud/ServiceCore/src/domain/types"
)

// Config is a struct that holds the configuration of the service
type Config struct {
	Version      types.Version
	Repository   *Repository
	Database     *Database
	Metrics      *Database
	GalaxyServer *GalaxyServer
	RpcServer    *RpcServer
}

func NewConfig() *Config {
	conf := Config{}

	conf.getVersion()
	conf.getRepository()
	conf.getDatabase()
	conf.getRPC()
	conf.getGalaxy()

	return &conf
}
