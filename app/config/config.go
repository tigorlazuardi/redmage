package config

import (
	"time"

	"github.com/inhies/go-bytesize"
)

type Config struct {
	Profiles   []Profile   `yaml:"profiles" koanf:"profiles"`
	Subreddits []Subreddit `yaml:"subreddits" koanf:"subreddits"`
	Download   Download    `yaml:"download" koanf:"download"`
}

type Profile struct {
	Name string `yaml:"name" koanf:"name"`

	AspectRatioX         float64 `yaml:"aspect_ratio_x" koanf:"aspect_ratio_x"`
	AspectRatioY         float64 `yaml:"aspect_ratio_y" koanf:"aspect_ratio_y"`
	AspectRatioTolerance float64 `yaml:"aspect_ratio_tolerance" koanf:"aspect_ratio_tolerance"`

	MinX float64 `yaml:"min_x" koanf:"min_x"`
	MaxX float64 `yaml:"max_x" koanf:"max_x"`

	MinY float64 `yaml:"min_y" koanf:"min_y"`
	MaxY float64 `yaml:"max_y" koanf:"max_y"`

	NSFW         bool   `yaml:"nsfw" koanf:"nsfw"`
	NamingFormat string `yaml:"naming_format" koanf:"naming_format"`
}

type Subreddit struct {
	Name        string `yaml:"name" koanf:"name"`
	Schedule    string `yaml:"schedule" koanf:"schedule"`
	LookupCount int    `yaml:"lookup_count" koanf:"lookup_count"`
}

func (s Subreddit) Count() int {
	if s.LookupCount == 0 {
		return 100
	}
	return s.LookupCount
}

type Download struct {
	Directory             string            `yaml:"directory" koanf:"directory"`
	Concurrency           int               `yaml:"concurrency" koanf:"concurrency"`
	ConnectionTimeout     time.Duration     `yaml:"connection_timeout" koanf:"connection_timeout"`
	DownloadIdleTimeout   time.Duration     `yaml:"download_idle_timeout" koanf:"download_idle_timeout"`
	DownloadIdleThreshold bytesize.ByteSize `yaml:"download_idle_threshold" koanf:"download_idle_threshold"`
}
