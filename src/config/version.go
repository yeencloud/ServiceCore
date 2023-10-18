package config

import (
	"github.com/yeencloud/ServiceCore/src/domain/types"
)

func (config *Config) getVersion() {
	if config.Version != "" {
		return
	}

	envVar := config.GetEnvStringOrDefault("GITHUB_SHA", "unknown")
	config.Version = types.Version(envVar)
}
