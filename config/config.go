package config

import (
	"github.com/knadh/koanf/v2"
)

type Config struct {
	*koanf.Koanf
}

func EmptyConfig() *Config {
	return NewConfigBuilder().Build()
}
