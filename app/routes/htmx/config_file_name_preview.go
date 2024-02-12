package htmx

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/models/render"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func (r *Routes) ConfigFileNamePreview(c echo.Context) error {
	ctx := render.Context{
		Echo:   c,
		Config: r.Config,
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
		c.FormValue("naming_format"),
	).Render(c.Request().Context(), c.Response())
}
