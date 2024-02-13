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
	Server Server `yaml:"server" koanf:"server" json:"server" schema:"server"`

	Devices    map[string]Device          `yaml:"profiles" koanf:"profiles" json:"profiles"`
	Subreddits map[string]SubredditConfig `yaml:"subreddits" koanf:"subreddits" json:"subreddits"`
	Download   Download                   `yaml:"download" koanf:"download" json:"download"`
	HotReload  bool                       `yaml:"hot_reload" koanf:"hot_reload" json:"hot_reload"`

	Koanf *koanf.Koanf `json:"-" yaml:"-" koanf:"-"`

	ConfigFile string `json:"-" yaml:"-" koanf:"-"`
	mu         sync.Mutex
}

type Server struct {
	Host string `yaml:"host" koanf:"host" json:"host" schema:"host"`
	Port int    `yaml:"port" koanf:"port" json:"port" schema:"port"`
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

type Device struct {
	Name         string `yaml:"name" koanf:"name" json:"name" schema:"name,required"`
	NSFW         bool   `yaml:"nsfw" koanf:"nsfw" json:"nsfw" schema:"nsfw"`
	NamingFormat string `yaml:"naming_format" koanf:"naming_format" json:"naming_format" schema:"naming_format"`

	AspectRatioWidth     int     `yaml:"aspect_ratio_width" koanf:"aspect_ratio_width" json:"aspect_ratio_width" schema:"aspect_ratio_width"`
	AspectRatioHeight    int     `yaml:"aspect_ratio_height" koanf:"aspect_ratio_height" json:"aspect_ratio_height" schema:"aspect_ratio_height"`
	AspectRatioTolerance float64 `yaml:"aspect_ratio_tolerance" koanf:"aspect_ratio_tolerance" json:"aspect_ratio_tolerance" schema:"aspect_ratio_tolerance"`

	MinWidth  float64 `yaml:"min_width" koanf:"min_width" json:"min_width" schema:"min_width"`
	MinHeight float64 `yaml:"min_height" koanf:"min_height" json:"min_height" schema:"min_height"`

	MaxWidth  float64 `yaml:"max_width" koanf:"max_width" json:"max_width" schema:"max_width"`
	MaxHeight float64 `yaml:"max_height" koanf:"max_height" json:"max_height" schema:"max_height"`
}

type SubredditConfig struct {
	Name        string `yaml:"name" koanf:"name" json:"name"`
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
