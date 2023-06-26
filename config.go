package ServiceCore

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	environment Environment
}

func (config *Config) loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found.")
	}
}

func (c *Config) GetEnvString(name string) (string, error) {
	envVariable, exists := os.LookupEnv(name)

	if exists == false {
		log.Println("Env variable: ", name, " is not found")
		return "", errors.New("env variable not found")
	}

	if strings.HasPrefix(envVariable, "$") {
		varNameWithoutPrefix := strings.TrimPrefix(envVariable, "$")
		log.Println("Env variable: ", name, " refers to: ", varNameWithoutPrefix)
		return c.GetEnvString(varNameWithoutPrefix)
	}

	if strings.Contains(name, "SECRET") || strings.Contains(name, "KEY") || strings.Contains(name, "TOKEN") || strings.Contains(name, "PASSWORD") {
		log.Println("[", name, "] = ****")
	} else {
		log.Println("[", name, "] =", envVariable)
	}

	return envVariable, nil
}

func (c *Config) GetEnvStringOrDefault(name string, defaultValue string) string {
	envVariable, err := c.GetEnvString(name)

	if err != nil {
		return defaultValue
	}

	return envVariable
}

func (c *Config) GetEnvInt(name string) (int, error) {
	value, err := c.GetEnvString(name)

	if err != nil {
		return 0, err
	}

	envVariable, err := strconv.Atoi(value)

	if err != nil {
		log.Println(name, " value is invalid, should be integer")
		return 0, errors.New("value is not int")
	}

	return envVariable, nil
}

func (c *Config) GetEnvIntOrDefault(name string, defaultValue int) int {
	envVariable, err := c.GetEnvInt(name)

	if err != nil {
		return defaultValue
	}

	return envVariable
}

func (c *Config) GetEnvBool(name string) (bool, error) {
	value, err := c.GetEnvString(name)

	if err != nil {
		return false, err
	}

	return value == "true" || value == "TRUE" || value == "1" || value == "yes" || value == "YES", nil
}

func (c *Config) GetEnvBoolOrDefault(name string, defaultValue bool) bool {
	envVariable, err := c.GetEnvBool(name)

	if err != nil {
		return defaultValue
	}

	return envVariable
}

func (c *Config) RequireEnvString(name string) string {
	value, err := c.GetEnvString(name)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return value
}

func (c *Config) RequireEnvInt(name string) int {
	envVariable, err := c.GetEnvInt(name)

	if err != nil {
		log.Fatalln(name, " value is invalid")
	}

	return envVariable
}

func (c *Config) RequireEnvBool(name string) bool {
	value, err := c.GetEnvBool(name)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return value
}

func newConfig() *Config {
	conf := Config{}

	conf.loadEnv()

	return &conf
}