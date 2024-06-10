package events

import (
	"github.com/teivah/broadcast"
	apievents "github.com/tigorlazuardi/redmage/api/events"
	"github.com/tigorlazuardi/redmage/config"
)

type Handler struct {
	Config    *config.Config
	Broadcast *broadcast.Relay[apievents.Event]
}

func NewHandler(cfg *config.Config, broadcast *broadcast.Relay[apievents.Event]) *Handler {
	return &Handler{
		Config:    cfg,
		Broadcast: broadcast,
	}
}

func (handler *Handler) Subscribe() (<-chan apievents.Event, func()) {
	listener := handler.Broadcast.Listener(10)
	return listener.Ch(), listener.Close
}
