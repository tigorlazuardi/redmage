package config

import (
	"strings"
	"text/template"
	"time"
)

const defaultNamingFormat = "{{ .Device.Name }}/{{ .Subreddit.Name }}/{{ .Image.DownloadedAt.Unix }}_{{ .Image.ID }}.{{ .Image.Extension }}"

type PreviewOptions struct {
	Device    Device
	Subreddit SubredditConfig
	Image     PreviewImageOption
}

type Date struct {
	time.Time
}

func (t *Date) String() string {
	return t.Format(time.DateOnly)
}

type PreviewImageOption struct {
	ID           string
	DownloadedAt Date
	PostedAt     Date
	Extension    string
}

func (c *Config) RenderFileNamingPreview() string {
	return c.RenderFileNamingPreviewWithOverride("", "", "")
}

func (c *Config) RenderFileNamingPreviewWithOverride(device string, subreddit string, format string) string {
	format = strings.TrimLeft(strings.TrimSpace(format), "/")
	if format == "" {
		format = defaultNamingFormat
	}
	if subreddit == "" {
		subreddit = "wallpapers"
	}
	if device == "" {
		device = "My-Laptop"
	}
	if d, ok := c.Devices[device]; ok {
		if d.NamingFormat != "" {
			format = d.NamingFormat
		}
	}
	tmpl, err := template.New("").Parse(c.Download.Directory + "/" + format)
	if err != nil {
		return "Error: bad template format: " + err.Error()
	}

	vars := PreviewOptions{
		Device: Device{
			Name:                 device,
			NSFW:                 true,
			NamingFormat:         format,
			AspectRatioWidth:     1920,
			AspectRatioHeight:    1080,
			AspectRatioTolerance: 0.2,
			MinWidth:             1920,
			MaxWidth:             4096,
			MinHeight:            1080,
			MaxHeight:            2160,
		},
		Subreddit: SubredditConfig{
			Name:        subreddit,
			Schedule:    "0 0 * * *",
			LookupCount: 100,
		},
		Image: PreviewImageOption{
			DownloadedAt: Date{time.Now()},
			Extension:    "jpg",
			PostedAt:     Date{time.Now().Add(-24 * time.Hour)},
			ID:           "0123456789ABCDEF",
		},
	}

	s := &strings.Builder{}
	err = tmpl.Execute(s, vars)
	if err != nil {
		return "Error: failed parsing template:\n" + err.Error()
	}

	return s.String()
}
