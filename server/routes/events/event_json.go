package events

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (handler *Handler) JSONEvents(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.HTMXEvents")
	defer span.End()

	flush, ok := rw.(http.Flusher)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": "response writer does not support streaming"})
		return
	}
	filters := strings.Split(r.URL.Query().Get("filter"), ",")

	log.New(ctx).Info("new htmx event stream connection", "user_agent", r.UserAgent())
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.WriteHeader(200)
	flush.Flush()

	ev, close := handler.Subscribe()
	defer close()

loop:
	for {
		select {
		case <-r.Context().Done():
			log.New(ctx).Info("simple event stream connection closed", "user_agent", r.UserAgent())
			return
		case event := <-ev:
			msg := event.Event()
			for _, filter := range filters {
				if filter != msg {
					continue loop
				}
			}
			if _, err := fmt.Fprintf(rw, "event: %s\ndata: ", msg); err != nil {
				return
			}
			if err := json.NewEncoder(rw).Encode(event); err != nil {
				log.New(ctx).Err(err).Error("failed to send json event", "user_agent", r.UserAgent())
				return
			}
			// single new line because JSON already has new line
			if _, err := io.WriteString(rw, "\n"); err != nil {
				return
			}
			flush.Flush()
		}
	}
}
