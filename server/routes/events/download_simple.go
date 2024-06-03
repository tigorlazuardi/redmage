package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (handler *Handler) SimpleDownloadEvent(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.EventsAPI")
	defer span.End()

	flush, ok := rw.(http.Flusher)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": "response writer does not support streaming"})
		return
	}

	log.New(ctx).Info("new simple event stream connection", "user_agent", r.UserAgent())

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.WriteHeader(200)
	flush.Flush()

	ev, close := handler.Subscribe()
	defer close()

	for {
		select {
		case <-r.Context().Done():
			log.New(ctx).Info("simple event stream connection closed", "user_agent", r.UserAgent())
			return
		case event := <-ev:
			msg := event.Event()
			if _, err := fmt.Fprintf(rw, "event: %s\ndata: %s\n\n", msg, msg); err != nil {
				return
			}
			flush.Flush()
		}
	}
}
