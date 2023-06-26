package ServiceCore

type Environment string

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
)

func (c *Config) GetEnvironment() Environment {
	if c.environment == "" {
		envVar := c.GetEnvStringOrDefault("ENV", EnvironmentDevelopment)

		c.environment = EnvironmentDevelopment
		if envVar == "production" || envVar == "prod" {
			c.environment = EnvironmentProduction
		}
	}

	return c.environment
}

func (c *Config) IsDevelopment() bool {
	return c.GetEnvironment() == EnvironmentDevelopment
}