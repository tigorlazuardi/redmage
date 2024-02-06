package routes

import (
	"io"

	"github.com/labstack/echo/v5"
)

func (r *Routes) HotReloadApi(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(200)
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-r.HotReload:
			_, err := io.WriteString(c.Response(), "data: Hot reload triggered\n\n")
			if err != nil {
				return err
			}
			c.Response().Flush()
		}
	}
}
