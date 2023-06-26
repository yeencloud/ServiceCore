package ServiceCore

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
)

// If we are running in kubernetes, galaxy's service (galaxy) should be exposed as an environment variable by kubernetes
// Otherwise you can set the environment variables GALAXY_HOST and GALAXY_PORT
// If you don't set them, it will default to localhost:3000
func (conf *Config) GetGalaxyAddress() string {
	if KubernetesUtil.IsRunningInKubernetes() {
		str, err := conf.GetEnvString("GALAXY_SERVICE_HOST")

		if err != nil {
			log.Err(err).Msg("Failed to get address for galaxy service from kubernetes")
		}

		return str
	} else {
		return conf.GetEnvStringOrDefault("GALAXY_HOST", "localhost")
	}
}

func (conf *Config) GetGalaxyPort() int {
	if KubernetesUtil.IsRunningInKubernetes() {
		port, err := conf.GetEnvInt("GALAXY_SERVICE_PORT")

		if err != nil {
			log.Err(err).Msg("Failed to get address for galaxy service from kubernetes")
		}

		return port
	} else {
		return conf.GetEnvIntOrDefault("GALAXY_PORT", 3000)
	}
}

func (c *Config) GetRPCPort() int {
	if KubernetesUtil.IsRunningInKubernetes() {
		return KubernetesUtil.GetInternalServicePort()
	} else {
		return c.GetEnvIntOrDefault("RPC_PORT", 0)
	}
}