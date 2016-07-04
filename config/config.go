// Package config contains functions to read Env variables (String, Int, Bool).
// It reads config variables for this project and stores them in the structure for further use
package config

import (
	"os"
	"strconv"
)

// Config stores environment variables
type Config struct {
	DbName       string // name of the psql database
	DbUser       string // user of the psql database
	DbHost       string // psql host
	DbPass       string // psql password
	DbPort       int    // psql port
	HttpPort     int    // http server port
	Secret       []byte // a key with which JWT token is signed
	ExpDays      int    // for how long is JWT token valid
	SaltLen      int    // the length of the salt of user password (hashed with scrypt)
	MailDomain   string // domain name of the mailgun
	MailPrivate  string // private key for the mailgun
	MailPublic   string // public key for the mailgun
	MaxImgSizeKb int64  // maximum possible size of uploaded image in kilobytes
	IsTest       bool   // whether this is a testing environment. Some functions behave differently
	TestEmail    string // all mail to all users will be sent to this address in test environments
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
		GetEnvInt("PROJ_SALT_LEN_BYTE"),
		GetEnvStr("PROJ_MAILGUN_DOMAIN"),
		GetEnvStr("PROJ_MAILGUN_PRIVATE"),
		GetEnvStr("PROJ_MAILGUN_PUBLIC"),
		int64(GetEnvInt("PROJ_MAX_IMG_KB")) * 1024,
		GetEnvBool("PROJ_IS_TEST"),
		GetEnvStr("PROJ_TEST_EMAIL"),
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

// GetEnvBool returns a environment variable as a boolean. Accepts only true, false. Panics otherwise
func GetEnvBool(key string) bool {
	val, err := strconv.ParseBool(GetEnvStr(key))
	if err != nil {
		panic("")
	}
	return val
}
