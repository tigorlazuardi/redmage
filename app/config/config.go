package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/inhies/go-bytesize"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Profiles   map[string]Profile         `yaml:"profiles" koanf:"profiles" json:"profiles"`
	Subreddits map[string]SubredditConfig `yaml:"subreddits" koanf:"subreddits" json:"subreddits"`
	Download   Download                   `yaml:"download" koanf:"download" json:"download"`
	HotReload  bool                       `yaml:"hot_reload" koanf:"hot_reload" json:"hot_reload"`

	Koanf *koanf.Koanf `json:"-" yaml:"-" koanf:"-"`

	ConfigFile string `json:"-" yaml:"-" koanf:"-"`
	mu         sync.Mutex
}

func (c *Config) Sync() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Koanf = koanf.New(".")
	err := c.Koanf.Load(structs.Provider(c, "koanf"), nil)
	if err != nil {
		return fmt.Errorf("error syncing config: %w", err)
	}

	b, err := c.Koanf.Marshal(yaml.Parser())
	if err != nil {
		return fmt.Errorf("error syncing config: fail to serialize data to yaml: %w", err)
	}

	if err := os.WriteFile(c.ConfigFile, b, 0644); err != nil {
		return fmt.Errorf("error syncing config: fail to write config to file: %w", err)
	}

	return nil
}

type Profile struct {
	AspectRatioX         float64 `yaml:"aspect_ratio_x" koanf:"aspect_ratio_x" json:"aspect_ratio_x"`
	AspectRatioY         float64 `yaml:"aspect_ratio_y" koanf:"aspect_ratio_y" json:"aspect_ratio_y"`
	AspectRatioTolerance float64 `yaml:"aspect_ratio_tolerance" koanf:"aspect_ratio_tolerance" json:"aspect_ratio_tolerance"`

	MinX float64 `yaml:"min_x" koanf:"min_x" json:"min_x"`
	MaxX float64 `yaml:"max_x" koanf:"max_x" json:"max_x"`

	MinY float64 `yaml:"min_y" koanf:"min_y" json:"min_y"`
	MaxY float64 `yaml:"max_y" koanf:"max_y" json:"max_y"`

	NSFW         bool   `yaml:"nsfw" koanf:"nsfw" json:"nsfw"`
	NamingFormat string `yaml:"naming_format" koanf:"naming_format" json:"naming_format"`
}

type SubredditConfig struct {
	Schedule    string `yaml:"schedule" koanf:"schedule" json:"schedule"`
	LookupCount int    `yaml:"lookup_count" koanf:"lookup_count" json:"lookup_count"`
}

func (s SubredditConfig) Count() int {
	if s.LookupCount == 0 {
		return 100
	}
	return s.LookupCount
}

type Download struct {
	Directory             string            `yaml:"directory" koanf:"directory" json:"directory" schema:"directory"`
	Concurrency           int               `yaml:"concurrency" koanf:"concurrency" json:"concurrency" schema:"concurrency"`
	ConnectionTimeout     time.Duration     `yaml:"connection_timeout" koanf:"connection_timeout" json:"connection_timeout" schema:"connection_timeout"`
	DownloadIdleTimeout   time.Duration     `yaml:"download_idle_timeout" koanf:"download_idle_timeout" json:"download_idle_timeout" schema:"download_idle_timeout"`
	DownloadIdleThreshold bytesize.ByteSize `yaml:"download_idle_threshold" koanf:"download_idle_threshold" json:"download_idle_threshold" schema:"download_idle_threshold"`
}
