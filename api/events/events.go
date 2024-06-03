package events

import (
	"io"

	"github.com/a-h/templ"
)

type Event interface {
	templ.Component
	// Event returns the event name
	Event() string
	// SerializeTo writes the event data to the writer.
	//
	// SerializeTo must not write multiple linebreaks (single linebreak is fine)
	// in succession to the writer since it will mess up SSE events.
	SerializeTo(w io.Writer) error
}
