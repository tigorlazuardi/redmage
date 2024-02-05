package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/tmpl"
)

type Routes struct{}

func (r *Routes) Home(c echo.Context) error {
	return tmpl.Hello("Tigor").Render(c.Request().Context(), c.Response())
}
