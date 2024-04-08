package www

import (
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (www *WWW) CreateHotReloadRoute() http.HandlerFunc {
	var mu sync.Mutex
	knownClients := make(map[string]chan struct{})
	firstTime := make(chan struct{}, 1)
	firstTime <- struct{}{}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var ch chan struct{}
		if oldChannel, known := knownClients[id]; known {
			ch = oldChannel
		} else {
			ch = make(chan struct{}, 1)
			ch <- struct{}{}
			knownClients[id] = ch
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(200)
		for {
			select {
			case <-r.Context().Done():
				return
			case <-firstTime:
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
				_, err := io.WriteString(w, "data: Hot reload triggered\n\n")
				if err != nil {
					log.New(r.Context()).Err(err).Error("failed to send hot reload trigger", "channel_id", id)
					return
				}
				if w, ok := w.(http.Flusher); ok {
					w.Flush()
				} else {
					panic(errors.New("HotReload: ResponseWriter does not implement http.Flusher interface"))
				}
			}
		}
	}
}
