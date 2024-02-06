package routes

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (r *Routes) createHotReloadRoute() func(c echo.Context) error {
	sseChans := make(map[*http.Request]chan struct{})
	return func(c echo.Context) error {
		ch := make(chan struct{}, 1)
		sseChans[c.Request()] = ch
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Content-Type", "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")
		c.Response().WriteHeader(200)
		for {
			select {
			case <-c.Request().Context().Done():
				delete(sseChans, c.Request())
				return nil
			case <-r.HotReload:
				// broadcast to all connected clients
				// that hot reload is triggered.
				//
				// The sender only send one signal,
				// and the receiver is only one, and chosen at random, so
				// we have to propagate the signal to all
				// connected clients.
				for _, ch := range sseChans {
					ch <- struct{}{}
				}
			case <-ch:
				_, err := io.WriteString(c.Response(), "data: Hot reload triggered\n\n")
				if err != nil {
					return err
				}
				c.Response().Flush()
			}
		}
	}
}
