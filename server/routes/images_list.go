package routes

import (
	"encoding/json"
	"net/http"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
)

func (routes *Routes) ImagesListAPI(rw http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "*Routes.ImagesList")
	defer span.End()

	var params api.ImageListParams

	params.FillFromQuery(r.URL.Query())

	enc := json.NewEncoder(rw)

	result, err := routes.API.ImagesListWithDevicesAndSubreddits(ctx, params)
	if err != nil {
		log.New(ctx).Err(err).Error("failed to list images")
		code, message := errs.HTTPMessage(err)
		rw.WriteHeader(code)
		_ = enc.Encode(map[string]string{"error": message})
		return
	}

	if err := enc.Encode(result); err != nil {
		log.New(ctx).Err(err).Error("failed to encode images")
	}
}
