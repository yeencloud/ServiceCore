package config

type Version string

func (config *Config) GetVersion() Version {
	if config.version == "" {
		envVar := config.GetEnvStringOrDefault("GITHUB_SHA", "unknown")

		config.version = Version(envVar)
	}

	return config.version
}