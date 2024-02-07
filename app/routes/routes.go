package routes

import (
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/models/render"
	"github.com/tigorlazuardi/redmage/app/routes/htmx"
)

type Routes struct {
	Config    *config.Config
	HotReload chan struct{}
}

func (r *Routes) Register(e *echo.Echo) {
	e.GET("/", r.HomePage)
	e.GET("/config", r.ConfigPage)
	if r.Config.HotReload {
		e.GET("/hot_reload", r.createHotReloadRoute())
	}

	htmxGroup := e.Group("/htmx")
	htmxRoutes := htmx.Routes{
		Config: r.Config,
		Schema: schema.NewDecoder(),
	}
	htmxRoutes.RegisterV1(htmxGroup)
}

func (r *Routes) renderContext(c echo.Context) render.Context {
	return render.Context{
		Echo:   c,
		Config: r.Config,
	}
}
