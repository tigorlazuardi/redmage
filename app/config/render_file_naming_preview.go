package config

import (
	"strings"
	"text/template"
	"time"
)

const defaultNamingFormat = "{{ .Device.Name }}/{{ .Subreddit.Name }}/{{ .Image.DownloadedAt.Unix }}_{{ .Image.ID }}.{{ .Image.Extension }}"

func (c *Config) RenderFileNamingPreview() string {
	namingFormat := defaultNamingFormat
	if len(c.Profiles) > 0 {
		for _, profile := range c.Profiles {
			if profile.NamingFormat != "" {
				namingFormat = profile.NamingFormat
				break
			}
		}
	}
	tmpl, err := template.New("").Parse(c.Download.Directory + "/" + namingFormat)
	if err != nil {
		return "Error: bad template format: " + err.Error()
	}

	vars := map[string]any{
		"Device": map[string]any{
			"Name": "My-Laptop",
		},
		"Subreddit": map[string]any{
			"Name": "aww",
		},
		"Image": map[string]any{
			"DownloadedAt": time.Now(),
			"ID":           "AAAAAAA",
			"Extension":    "jpg",
			"PostedAt":     time.Now().Add(-time.Hour * 24 * 7),
		},
	}

	s := &strings.Builder{}
	err = tmpl.Execute(s, vars)
	if err != nil {
		return "Error: failed executing template: " + err.Error()
	}

	return s.String()
}
