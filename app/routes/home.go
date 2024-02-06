package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/tmpl"
)

func (r *Routes) Home(c echo.Context) error {
	return tmpl.
		Home(r.renderContext(c)).
		Render(c.Request().Context(), c.Response())
}
