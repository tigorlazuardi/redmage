package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

func Load() (*Config, error) {
	cfg := Default()

	fp := file.Provider(cfg.ConfigFile)
	err := cfg.Koanf.Load(fp, yaml.Parser())
	if err != nil {
		return cfg, err
	}

	_ = cfg.Koanf.Load(env.Provider("REDMAGE_", ".", func(s string) string {
		key := strings.TrimPrefix(s, "REDMAGE_")
		key = strings.ToLower(key)
		key = strings.ReplaceAll(key, "__", ".")
		return key
	}), nil)

	_ = cfg.Koanf.Unmarshal("", cfg)

	return cfg, nil
}
