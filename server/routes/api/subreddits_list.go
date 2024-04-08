package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tigorlazuardi/redmage/db/queries/subreddits"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (api *API) ListSubreddits(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	subs, err := api.Subreddits.ListSubreddits(r.Context(), parseListSubredditsQuery(r))
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

func parseListSubredditsQuery(r *http.Request) subreddits.ListSubredditsParams {
	params := subreddits.ListSubredditsParams{}
	params.Limit, _ = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	params.Offset, _ = strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)

	if params.Limit < 1 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
	}

	return params
}
