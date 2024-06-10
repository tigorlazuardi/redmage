package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tigorlazuardi/redmage/pkg/log"
)

// SimpleEvents is a simple event stream for the purpose of
// notification that something did happen. Not what the content of the event is.
//
// Useful for simple notification whose client just need to know that something
// happened and do something that does not require the content of the event,
// like refreshing the list by calling another http request.
func (handler *Handler) SimpleEvents(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.SimpleDownloadEvent")
	defer span.End()

	flush, ok := rw.(http.Flusher)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": "response writer does not support streaming"})
		return
	}

	var filters []string
	if q := r.URL.Query().Get("filter"); q != "" {
		filters = strings.Split(q, ",")
	}

	log.New(ctx).Info("new simple event stream connection", "user_agent", r.UserAgent())

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
			if _, err := fmt.Fprintf(rw, "event: %s\ndata: %s\n\n", msg, msg); err != nil {
				return
			}
			flush.Flush()
		}
	}
}
