package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
	"gopkg.in/yaml.v3"
)

func (r *Routes) ConfigPage(c echo.Context) error {
	cfg, _ := yaml.Marshal(r.Config)

	return pages.
		Config(r.renderContext(c), string(cfg)).
		Render(c.Request().Context(), c.Response())
}
