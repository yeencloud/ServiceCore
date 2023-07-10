package config

type Environment string

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
)

func (config *Config) GetEnvironment() Environment {
	if config.environment == "" {
		envVar := config.GetEnvStringOrDefault("ENV", EnvironmentDevelopment)

		config.environment = EnvironmentDevelopment
		if envVar == "production" || envVar == "prod" {
			config.environment = EnvironmentProduction
		}
	}

	return config.environment
}

func (config *Config) IsDevelopment() bool {
	return config.GetEnvironment() == EnvironmentDevelopment
}