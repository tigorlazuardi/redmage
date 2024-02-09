package routes

import (
	"io"
	"sync"

	"github.com/labstack/echo/v5"
)

func (r *Routes) createHotReloadRoute() func(c echo.Context) error {
	var mu sync.Mutex
	knownClients := make(map[string]chan struct{})
	return func(c echo.Context) error {
		id := c.Request().URL.Query().Get("id")
		var ch chan struct{}
		if oldChannel, known := knownClients[id]; known {
			ch = oldChannel
		} else {
			ch = make(chan struct{}, 1)
			ch <- struct{}{}
			knownClients[id] = ch
		}

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
				// dispose own's channel buffer to prevent deadlock.
				// because it will be filled with new signal below.
				if len(ch) > 0 {
					<-ch
				}
				// broadcast to all connected clients
				// that hot reload is triggered.
				//
				// The sender only send one signal,
				// and the receiver is only one, and chosen at random, so
				// we have to propagate the signal to all
				// connected clients.
				mu.Lock()
				for _, ch := range knownClients {
					ch <- struct{}{}
				}
				mu.Unlock()
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
