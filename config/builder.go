package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

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
	provider := confmap.Provider(DefaultConfig, ".")

	err := builder.koanf.Load(provider, nil)
	if err != nil {
		panic(err)
	}
	return builder
}

func (builder *ConfigBuilder) LoadJSON(b []byte) *ConfigBuilder {
	provider := rawbytes.Provider(b)
	err := builder.koanf.Load(provider, json.Parser())
	if err != nil {
		builder.err = err
	}
	return builder
}

func (builder *ConfigBuilder) LoadJSONFile(path string) *ConfigBuilder {
	if _, err := os.Stat(path); err != nil {
		return builder
	}
	provider := file.Provider(path)
	err := builder.koanf.Load(provider, json.Parser())
	if err != nil {
		builder.err = err
	}
	return builder
}

func (builder *ConfigBuilder) LoadYaml(b []byte) *ConfigBuilder {
	provider := rawbytes.Provider(b)
	err := builder.koanf.Load(provider, yaml.Parser())
	if err != nil {
		builder.err = err
	}
	return builder
}

func (builder *ConfigBuilder) LoadYamlFile(path string) *ConfigBuilder {
	if _, err := os.Stat(path); err != nil {
		return builder
	}
	provider := file.Provider(path)
	err := builder.koanf.Load(provider, yaml.Parser())
	if err != nil {
		builder.err = err
	}
	return builder
}

func (builder *ConfigBuilder) LoadEnv() *ConfigBuilder {
	provider := env.Provider("REDMAGE_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "REDMAGE_")
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "_", ".")
		return s
	})

	_ = builder.koanf.Load(provider, nil)
	return builder
}

func (builder *ConfigBuilder) LoadFlags(flags *pflag.FlagSet) *ConfigBuilder {
	provider := posflag.Provider(flags, ".", nil)
	if err := builder.koanf.Load(provider, nil); err != nil {
		builder.err = err
	}

	return builder
}
