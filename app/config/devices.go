package config

import (
	"sort"
	"strings"
)

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

func (c *Config) GetSortedDevices() []Device {
	c.mu.Lock()
	defer c.mu.Unlock()

	devices := make([]Device, 0, len(c.Devices))
	for _, device := range c.Devices {
		devices = append(devices, device)
	}
	sort.Slice(devices, func(i, j int) bool {
		return strings.ToLower(devices[i].Name) < strings.ToLower(devices[j].Name)
	})
	return devices
}
