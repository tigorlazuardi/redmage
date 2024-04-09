package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (r *Routes) SubredditsListAPI(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	params := parseSubredditListQuery(req)

	result, err := r.API.ListSubreddits(ctx, params)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := json.NewEncoder(rw).Encode(result); err != nil {
		log.New(ctx).Err(err).Error("failed to encode subreddits into json")
	}
}

func (r *Routes) SubredditsListPage(rw http.ResponseWriter, req *http.Request) {
}

func parseSubredditListQuery(req *http.Request) (params api.ListSubredditsParams) {
	params.Name = req.FormValue("name")
	params.Limit, _ = strconv.ParseInt(req.FormValue("limit"), 10, 64)
	if params.Limit < 1 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
	}
	params.Offset, _ = strconv.ParseInt(req.FormValue("offset"), 10, 64)
	return params
}
