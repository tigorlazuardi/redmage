package events

import (
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/redmage/api/events"
	"github.com/tigorlazuardi/redmage/config"
)

type Handler struct {
	Config    *config.Config
	Broadcast *broadcast.Relay[events.Event]
}

func (handler *Handler) Subscribe() (<-chan events.Event, func()) {
	listener := handler.Broadcast.Listener(10)
	return listener.Ch(), listener.Close
}
