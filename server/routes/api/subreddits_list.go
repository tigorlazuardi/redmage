package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (api *API) ListSubreddits(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	subs, err := api.Subreddits.ListSubreddits(r.Context(), 10)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("failed to list subreddits: %s", err)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": msg})
		return
	}

	if err := json.NewEncoder(rw).Encode(subs); err != nil {
		log.New(r.Context()).Err(err).Error("failed to list subreddits")
	}
}
