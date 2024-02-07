package routes

import (
	"io"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

func (r *Routes) createHotReloadRoute() func(c echo.Context) error {
	type hotReloadClient struct {
		cleaner *time.Timer
		ch      chan struct{}
	}
	var mu sync.Mutex
	knownClients := make(map[string]hotReloadClient)
	cleanup := func(id string) *time.Timer {
		return time.AfterFunc(5*time.Minute, func() {
			mu.Lock()
			defer mu.Unlock()
			delete(knownClients, id)
		})
	}
	return func(c echo.Context) error {
		id := c.Request().URL.Query().Get("id")
		var ch chan struct{}
		if oldClient, known := knownClients[id]; known {
			oldClient.cleaner.Stop()
			oldClient.cleaner = cleanup(id)
			ch = oldClient.ch
		} else {
			client := hotReloadClient{
				cleaner: cleanup(id),
				ch:      make(chan struct{}, 1),
			}
			ch = client.ch
			ch <- struct{}{}
			knownClients[id] = client
		}

		c.Response().Header().Set("Content-Type", "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")
		c.Response().WriteHeader(200)

		for {
			select {
			case <-c.Request().Context().Done():
				return nil
			case <-r.HotReload:
				// broadcast to all connected clients
				// that hot reload is triggered.
				//
				// The sender only send one signal,
				// and the receiver is only one, and chosen at random, so
				// we have to propagate the signal to all
				// connected clients.
				mu.Lock()
				for _, ch := range knownClients {
					ch.ch <- struct{}{}
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
