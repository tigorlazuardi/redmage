package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/tmpl"
)

type Routes struct{}

func (r *Routes) Register(e *echo.Echo) {
	e.GET("/", r.Home)
	e.GET("/config", r.Config)
}

func (r *Routes) renderContext(c echo.Context) tmpl.RenderContext {
	return tmpl.RenderContext{
		Echo: c,
		// TODO: find some way to pass latest config
		Config: config.Default(),
	}
}
