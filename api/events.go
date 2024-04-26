package api

import (
	"github.com/tigorlazuardi/redmage/api/bmessage"
)

func (api *API) SubscribeImageDownloadEvent() (<-chan bmessage.ImageDownloadMessage, func()) {
	listener := api.downloadBroadcast.Listener(10)
	return listener.Ch(), listener.Close
}
