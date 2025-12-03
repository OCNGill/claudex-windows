package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Doc         []string `toml:"doc"`
	NoOverwrite bool     `toml:"no_overwrite"`
}

func loadConfig() (*Config, error) {
	config := &Config{
		Doc:         []string{},
		NoOverwrite: false,
	}

	if _, err := os.Stat(".claudex.toml"); err == nil {
		if _, err := toml.DecodeFile(".claudex.toml", config); err != nil {
			return nil, err
		}
	}
	return config, nil
}
