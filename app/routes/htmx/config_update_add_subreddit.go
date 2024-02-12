package htmx

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/models/render"
	"github.com/tigorlazuardi/redmage/app/templates/components"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func (r *Routes) ConfigUpdateAddSubreddit(c echo.Context) error {
	const defaultNamingFormat = "{{ .Device.Name }}/{{ .Subreddit.Name }}/{{ .Image.DownloadedAt.Unix }}_{{ .Image.ID }}.{{ .Image.Extension }}"
	type AddSubredditSchema struct {
		Name                 string  `schema:"name"`
		NSFW                 bool    `schema:"nsfw"`
		NamingFormat         string  `schema:"naming_format"`
		AspectRatioWidth     float64 `schema:"aspect_ratio_width"`
		AspectRatioHeight    float64 `schema:"aspect_ratio_height"`
		AspectRatioTolerance float64 `schema:"aspect_ratio_tolerance"`
		MinWidth             float64 `schema:"min_width"`
		MaxWidth             float64 `schema:"max_width"`
		MinHeight            float64 `schema:"min_height"`
		MaxHeight            float64 `schema:"max_height"`
	}

	if err := c.Request().ParseForm(); err != nil {
		c.Response().WriteHeader(400)
		return components.ErrorToast("Error reading request data data: %s", err.Error()).Render(c.Request().Context(), c.Response())
	}

	parsed := AddSubredditSchema{}
	if err := r.Schema.Decode(&parsed, c.Request().Form); err != nil {
		c.Response().WriteHeader(400)
		return components.ErrorToast("Error parsing form data: %s", err.Error()).Render(c.Request().Context(), c.Response())
	}

	parsed.NamingFormat = strings.TrimLeft(strings.TrimSpace(parsed.NamingFormat), "/")
	if parsed.NamingFormat == "" {
		parsed.NamingFormat = defaultNamingFormat
	} else if parsed.NamingFormat != "" && strings.Count(parsed.NamingFormat, "{{") == 0 && strings.Count(parsed.NamingFormat, "}}") == 0 {
		c.Response().WriteHeader(400)
		return components.
			ErrorToast("File Naming Format must contain at least one variable").
			Render(c.Request().Context(), c.Response())
	}

	dev := config.Device{
		Name:                 parsed.Name,
		NSFW:                 parsed.NSFW,
		NamingFormat:         defaultNamingFormat,
		AspectRatioWidth:     parsed.AspectRatioWidth,
		AspectRatioHeight:    parsed.AspectRatioHeight,
		AspectRatioTolerance: parsed.AspectRatioTolerance,
		MinWidth:             parsed.MinWidth,
		MaxWidth:             parsed.MaxWidth,
		MinHeight:            parsed.MinHeight,
		MaxHeight:            parsed.MaxHeight,
	}

	r.Config.Devices[parsed.Name] = dev
	if err := r.Config.Sync(); err != nil {
		c.Response().WriteHeader(500)
		return components.ErrorToast("Error saving config: %s", err.Error()).Render(c.Request().Context(), c.Response())
	}

	ctx := render.Context{
		Echo:   c,
		Config: r.Config,
	}

	return pages.ConfigDeviceSection(ctx).Render(c.Request().Context(), c.Response())
}
