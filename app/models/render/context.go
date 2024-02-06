package render

import (
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
)

type Context struct {
	Echo   echo.Context
	Config *config.Config
}
