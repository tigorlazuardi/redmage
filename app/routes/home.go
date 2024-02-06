package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/templates/pages"
)

func (r *Routes) Home(c echo.Context) error {
	return pages.
		Home(r.renderContext(c)).
		Render(c.Request().Context(), c.Response())
}
