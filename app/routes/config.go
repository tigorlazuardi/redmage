package routes

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func (r *Routes) ConfigPage(c echo.Context) error {
	cfg, _ := r.Config.Koanf.Marshal(yaml.Parser())

	return pages.
		Config(r.renderContext(c), string(cfg)).
		Render(c.Request().Context(), c.Response())
}
