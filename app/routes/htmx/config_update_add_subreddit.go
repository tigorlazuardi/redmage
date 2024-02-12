package htmx

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/models/render"
	"github.com/tigorlazuardi/redmage/app/templates/components"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func (r *Routes) ConfigUpdateAddSubreddit(c echo.Context) error {
	const defaultNamingFormat = "{{ .Device.Name }}/{{ .Subreddit.Name }}/{{ .Image.DownloadedAt.Unix }}_{{ .Image.ID }}.{{ .Image.Extension }}"

	if err := c.Request().ParseForm(); err != nil {
		c.Response().WriteHeader(400)
		return components.ErrorToast("Error reading request data data: %s", err.Error()).Render(c.Request().Context(), c.Response())
	}

	parsed := config.Device{}
	if err := r.Schema.Decode(&parsed, c.Request().Form); err != nil {
		c.Response().WriteHeader(400)
		return components.ErrorToast("Error parsing form data: %s", err.Error()).Render(c.Request().Context(), c.Response())
	}
	originalName := parsed.Name
	loweredName := strings.ToLower(parsed.Name)
	if _, exist := r.Config.Devices[loweredName]; exist {
		c.Response().WriteHeader(409)
		return components.ErrorToast("Device with name '%s' already exists", originalName).Render(c.Request().Context(), c.Response())
	}
	switch {
	case parsed.AspectRatioHeight != 0 && parsed.AspectRatioWidth == 0:
		c.Response().WriteHeader(400)
		return components.
			ErrorToast("Aspect Ratio Width cannot be 0 if Aspect Ratio Height is not 0").
			Render(c.Request().Context(), c.Response())
	case parsed.AspectRatioWidth != 0 && parsed.AspectRatioHeight == 0:
		c.Response().WriteHeader(400)
		return components.
			ErrorToast("Aspect Ratio Height cannot be 0 if Aspect Ratio Width is not 0").
			Render(c.Request().Context(), c.Response())
	case parsed.AspectRatioHeight != 0 && parsed.AspectRatioWidth != 0:
		g := gcd(parsed.AspectRatioWidth, parsed.AspectRatioHeight)
		parsed.AspectRatioWidth /= g
		parsed.AspectRatioHeight /= g
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

	r.Config.Devices[loweredName] = parsed
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
