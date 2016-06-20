// Package config contains functions to read Env variables (String, Int). It reads config variables
// for this project and stores them in the structure
package config

import (
	"os"
	"strconv"
)

// Config stores environment variables
type Config struct {
	DbName   string
	DbUser   string
	DbHost   string
	DbPass   string
	DbPort   int
	HttpPort int
	Secret   []byte
	ExpDays  int
}

var Cfg Config

// Init extracts all environment variables for further use
func Init() {
	cfg := Config{
		GetEnvStr("PROJ_DB_NAME"),
		GetEnvStr("PROJ_DB_USER"),
		GetEnvStr("PROJ_DB_HOST"),
		GetEnvStr("PROJ_DB_PWD"),
		GetEnvInt("PROJ_DB_PORT"),
		GetEnvInt("PROJ_HTTP_PORT"),
		[]byte(GetEnvStr("PROJ_SECRET")),
		GetEnvInt("PROJ_JWT_EXP_DAYS"),
	}
	Cfg = cfg
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
