package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	*koanf.Koanf
}

func EmptyConfig() *Config {
	return NewConfigBuilder().Build()
}

type ConfigBuilder struct {
	koanf *koanf.Koanf
	err   error
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{koanf: koanf.New(".")}
}

func (builder *ConfigBuilder) Build() *Config {
	return &Config{Koanf: builder.koanf}
}

func (builder *ConfigBuilder) BuildHandle() (*Config, error) {
	return &Config{Koanf: builder.koanf}, builder.err
}

func (builder *ConfigBuilder) LoadDefault() *ConfigBuilder {
	provider := confmap.Provider(map[string]any{
		"log.enable": true,
		"log.source": true,
		"log.format": "pretty",
		"log.level":  "info",
		"log.output": "stderr",
	}, ".")

	_ = builder.koanf.Load(provider, nil)
	return builder
}
