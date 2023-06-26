package servicecore

import (
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/rs/zerolog/log"
)

// GetGalaxyAddress GetGalaxyPort : If we are running in kubernetes, galaxy's service (galaxy) should be exposed as an environment variable by kubernetes
// Otherwise you can set the environment variables GALAXY_HOST and GALAXY_PORT
// If you don't set them, it will default to localhost:3000
func (config *Config) GetGalaxyAddress() string {
	if KubernetesUtil.IsRunningInKubernetes() {
		str, err := config.GetEnvString("GALAXY_SERVICE_HOST")

		if err != nil {
			log.Err(err).Msg("Failed to get address for galaxy service from kubernetes")
		}

		return str
	}
	return config.GetEnvStringOrDefault("GALAXY_HOST", "localhost")
}

func (config *Config) GetGalaxyPort() int {
	if KubernetesUtil.IsRunningInKubernetes() {
		port, err := config.GetEnvInt("GALAXY_SERVICE_PORT")

		if err != nil {
			log.Err(err).Msg("Failed to get address for galaxy service from kubernetes")
		}

		return port
	}
	return config.GetEnvIntOrDefault("GALAXY_PORT", 3000)
}

func (config *Config) GetRPCPort() int {
	if KubernetesUtil.IsRunningInKubernetes() {
		return KubernetesUtil.GetInternalServicePort()
	} else {
		return config.GetEnvIntOrDefault("RPC_PORT", 0)
	}
}