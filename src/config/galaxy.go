package config

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
	types2 "github.com/yeencloud/ServiceCore/src/domain/types"
)

type GalaxyServer struct {
	Host types2.Host
	Port types2.Port
}

// GetGalaxyAddress GetGalaxyPort : If we are running in kubernetes, galaxy's service (galaxy) should be exposed as an environment variable by kubernetes
// Otherwise you can set the environment variables GALAXY_HOST and GALAXY_PORT
// If you don't set them, it will default to localhost:3000
func (config *Config) getGalaxy() {
	if config.GalaxyServer != nil {
		return
	}

	glx := &GalaxyServer{}

	if KubernetesUtil.IsRunningInKubernetes() {
		str, err := config.GetEnvString("GALAXY_SERVICE_HOST")

		if err != nil {
			log.Err(err).Msg("Failed to get host for galaxy service from kubernetes")
		}

		glx.Host = types2.Host(str)
	} else {
		str, err := config.GetEnvString("GALAXY_HOST")

		if err != nil {
			log.Err(err).Msg("Failed to get host for galaxy service")
		}

		glx.Host = types2.Host(str)
	}

	if KubernetesUtil.IsRunningInKubernetes() {
		port, err := config.GetEnvInt("GALAXY_SERVICE_PORT")

		if err != nil {
			log.Err(err).Msg("Failed to get address for galaxy service from kubernetes")
		}

		glx.Port = types2.Port(port)
	} else {
		port, err := config.GetEnvInt("GALAXY_PORT")

		if err != nil {
			log.Err(err).Msg("Failed to get address for galaxy service")
		}

		glx.Port = types2.Port(port)
	}

	config.GalaxyServer = glx
}
