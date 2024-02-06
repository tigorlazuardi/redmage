package routes

import (
	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/tmpl"
)

func (r *Routes) Config(c echo.Context) error {
	cfg, _ := json.MarshalIndent(config.Default(), "", "  ")

	return tmpl.ConfigPage(r.renderContext(c), string(cfg)).
		Render(c.Request().Context(), c.Response())
}
