package types

type Environment string

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
)

func (env Environment) IsDevelopment() bool {
	return env == EnvironmentDevelopment
}
