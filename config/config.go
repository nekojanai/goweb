package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port         int
	PortAsString string
	DataPath     string
}

func (config *Config) Read(path string) error {
	_, err := toml.DecodeFile(path, &config)

	// special cases
	config.PortAsString = fmt.Sprintf(":%v", config.Port)

	return err
}
