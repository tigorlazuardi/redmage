package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (routes *Routes) EventsAPI(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.EventsAPI")
	defer span.End()

	flush, ok := rw.(http.Flusher)
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": "server does not support streaming"})
		return
	}

	log.New(ctx).Info("new event stream connection", "user_agent", r.UserAgent())

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.WriteHeader(200)
	flush.Flush()

	downloadEvent, closeDownloadEvent := routes.API.SubscribeImageDownloadEvent()
	defer closeDownloadEvent()

	enc := json.NewEncoder(rw)

	for {
		select {
		case <-r.Context().Done():
			log.New(ctx).Info("event stream connection closed", "user_agent", r.UserAgent())
			return
		case event := <-downloadEvent:
			if _, err := io.WriteString(rw, "event: image_download\n"); err != nil {
				return
			}
			if _, err := io.WriteString(rw, "data: "); err != nil {
				return
			}
			if err := enc.Encode(event); err != nil {
				return
			}
			// Single '\n' because enc.Encode already append '\n'
			if _, err := io.WriteString(rw, "\n"); err != nil {
				return
			}
		}
		flush.Flush()
	}
}
