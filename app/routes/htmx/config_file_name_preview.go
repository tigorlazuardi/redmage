package htmx

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/models/render"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func (r *Routes) ConfigFileNamePreview(c echo.Context) error {
	ctx := render.Context{
		Echo:   c,
		Config: r.Config,
	}
	format := c.FormValue("naming_format")
	if format != "" && strings.Count(format, "{{") == 0 && strings.Count(format, "}}") == 0 {
		c.Response().Header().Set("HX-Trigger", "invalid-naming-format")
		return pages.
			ConfigFilenameError("Error: File Naming Format must contain at least one variable, otherwise all images in the device profile will be overwritten.\nIf using given sets in the dropdown when typing, click on the Dropdown selection to expand them.").
			Render(c.Request().Context(), c.Response())
	}
	sub := c.FormValue("subreddit")
	if sub == "" {
		for _, s := range r.Config.Subreddits {
			sub = s.Name
			break
		}
		if sub == "" {
			sub = "wallpapers"
		}
	}

	return pages.ConfigFilenamePreviewOverride(ctx,
		c.FormValue("name"),
		sub,
		format,
	).Render(c.Request().Context(), c.Response())
}
