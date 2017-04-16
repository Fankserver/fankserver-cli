package config

import (
	"errors"
	"os"

	toml "gopkg.in/burntsushi/toml.v0"
)

var conf *Config

type Config struct {
	Jwt Jwt
	DB  map[string]Database `toml:"database"`
}

type Jwt struct {
	Secret string
}
type Database struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database string
}

// ReadConfigFile tries to parse the config file
func ReadConfigFile(path string) error {
	if path == "" {
		return errors.New("No config present")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("Config file does not exist: " + path)
	}
	_, err := toml.DecodeFile(path, &conf)
	return err
}

// GetConfig returns the parsed config
func GetConfig() Config {
	return *conf
}
