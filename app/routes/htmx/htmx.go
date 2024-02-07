package htmx

import (
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
)

type Routes struct {
	Config *config.Config
	Schema *schema.Decoder
}

func (r *Routes) RegisterV1(group *echo.Group) {
	v1Group := group.Group("/v1")
	v1Group.POST("/config/update/download", r.ConfigUpdateDownload)
}
