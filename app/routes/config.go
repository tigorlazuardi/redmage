package routes

import (
	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func (r *Routes) Config(c echo.Context) error {
	cfg, _ := json.MarshalIndent(config.Default(), "", "  ")

	return pages.
		Config(r.renderContext(c), string(cfg)).
		Render(c.Request().Context(), c.Response())
}
