package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/models/render"
)

type Routes struct {
	Config *config.Config
}

func (r *Routes) Register(e *echo.Echo) {
	e.GET("/", r.HomePage)
	e.GET("/config", r.ConfigPage)
}

func (r *Routes) renderContext(c echo.Context) render.Context {
	return render.Context{
		Echo:   c,
		Config: r.Config,
	}
}
