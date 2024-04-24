package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (routes *Routes) APIDeviceList(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.APIDeviceList")
	defer span.End()

	query := parseApiDeviceListQueries(r)

	result, err := routes.API.DevicesList(ctx, query)
	if err != nil {
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = json.NewEncoder(rw).Encode(map[string]string{"error": message})
		return
	}

	if err := json.NewEncoder(rw).Encode(result); err != nil {
		log.New(ctx).Err(err).Error("failed to marshal json api devices")
	}
}

func parseApiDeviceListQueries(req *http.Request) (params api.DevicesListParams) {
	params.All, _ = strconv.ParseBool(req.FormValue("all"))
	params.Offset, _ = strconv.ParseInt(req.FormValue("offset"), 10, 64)
	params.Limit, _ = strconv.ParseInt(req.FormValue("limit"), 10, 64)
	params.Q = req.FormValue("q")
	params.OrderBy = req.FormValue("order")
	params.Sort = strings.ToLower(req.FormValue("sort"))

	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.OrderBy == "" {
		params.OrderBy = "name"
	}

	return params
}
