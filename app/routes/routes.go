package routes

import (
	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/tmpl"
)

type Routes struct{}

func (r *Routes) Home(c echo.Context) error {
	cfg, _ := json.MarshalIndent(config.Default(), "", "    ")
	render := tmpl.RenderContext{
		Echo:   c,
		Config: config.Default(),
	}

	return tmpl.Home(render, string(cfg)).Render(c.Request().Context(), c.Response())
}
