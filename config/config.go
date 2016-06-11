package config

import (
	"os"
	"strconv"
)

// Config stores environment variables
type Config struct {
	DbName     string
	DBUser     string
	DbHost     string
	DBPassword string
	DbPort     int
}

// GetConfig extracts all environment variables for further use
func GetConfig() Config {
	return Config{
		GetEnvStr("PROJ_DB_NAME"),
		GetEnvStr("PROJ_DB_USER"),
		GetEnvStr("PROJ_DB_HOST"),
		GetEnvStr("PROJ_DB_PWD"),
		GetEnvInt("PROJ_DB_PORT"),
	}
}

// GetEnvStr returns a environment variable as a string. Panics if it does not exist
func GetEnvStr(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("")
	}
	return val
}

// GetEnvInt returns a environment variable as an integer. Panics if it is not an integer
func GetEnvInt(key string) int {
	val, err := strconv.Atoi(GetEnvStr(key))
	if err != nil {
		panic("")
	}
	return val
}