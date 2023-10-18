package config

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
)

func (config *Config) GetEnvString(name string) (string, error) {
	envVariable, exists := os.LookupEnv(name)

	if exists == false {
		log.Warn().Str("envvar", name).Msg("Env variable is not found")
		return "", errors.New("env variable not found")
	}

	if strings.HasPrefix(envVariable, "$") {
		varNameWithoutPrefix := strings.TrimPrefix(envVariable, "$")
		log.Info().Str("envvar", name).Str("refers", varNameWithoutPrefix).Msg("Env variable refers to another env variable")
		return config.GetEnvString(varNameWithoutPrefix)
	}

	if strings.Contains(name, "SECRET") || strings.Contains(name, "KEY") || strings.Contains(name, "TOKEN") || strings.Contains(name, "PASSWORD") {
		log.Info().Str("envvar", name).Msg("****")
	} else {
		log.Info().Str("envvar", name).Msg(envVariable)
	}

	return envVariable, nil
}

func (config *Config) GetEnvStringOrDefault(name string, defaultValue string) string {
	envVariable, err := config.GetEnvString(name)

	if err != nil {
		return defaultValue
	}

	return envVariable
}

func (config *Config) GetEnvInt(name string) (int, error) {
	value, err := config.GetEnvString(name)

	if err != nil {
		return 0, err
	}

	envVariable, err := strconv.Atoi(value)

	if err != nil {
		log.Err(err).Str("name", name).Msg(" value is invalid, should be integer")
		return 0, errors.New("value is not int")
	}

	return envVariable, nil
}

func (config *Config) GetEnvIntOrDefault(name string, defaultValue int) int {
	envVariable, err := config.GetEnvInt(name)

	if err != nil {
		return defaultValue
	}

	return envVariable
}

func (config *Config) GetEnvBool(name string) (bool, error) {
	value, err := config.GetEnvString(name)

	if err != nil {
		return false, err
	}

	return value == "true" || value == "TRUE" || value == "1" || value == "yes" || value == "YES", nil
}

func (config *Config) GetEnvBoolOrDefault(name string, defaultValue bool) bool {
	envVariable, err := config.GetEnvBool(name)

	if err != nil {
		return defaultValue
	}

	return envVariable
}

func (config *Config) RequireEnvString(name string) string {
	value, err := config.GetEnvString(name)

	if err != nil {
		log.Err(err).Str("name", name).Msg("Env variable is not found")
		os.Exit(1)
	}

	return value
}

func (config *Config) RequireEnvInt(name string) int {
	envVariable, err := config.GetEnvInt(name)

	if err != nil {
		log.Err(err).Str("name", name).Msg("Env variable is not found")
		os.Exit(1)
	}

	return envVariable
}

func (config *Config) RequireEnvBool(name string) bool {
	value, err := config.GetEnvBool(name)

	if err != nil {
		log.Err(err).Str("name", name).Msg("Env variable is not found")
		os.Exit(1)
	}

	return value
}
