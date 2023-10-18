package environment

import (
	"github.com/yeencloud/ServiceCore/src/domain/types"
	"os"
)

// Environment could have been a part of config since it basically read environment variables but we need that variable earlier for logging purposes

// IsProduction will be assumed to be true by default unless the environment variable SERVICE_ENVIRONMENT is set to "development"
// this is because we want to make sure that the default is production and that we don't accidentally run in development mode in production
func isProduction() bool {
	env, present := os.LookupEnv("SERVICE_ENVIRONMENT")

	if !present {
		return true
	}

	return env != "development"
}

func GetEnvironment() types.Environment {
	if isProduction() {
		return types.EnvironmentProduction
	}

	return types.EnvironmentDevelopment
}
