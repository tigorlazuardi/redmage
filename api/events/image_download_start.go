package events

import (
	"context"
	"encoding/json"
	"io"
)

type ImageDownloadEvent string

const (
	ImageDownloadStart    ImageDownloadEvent = "image.download.start"
	ImageDownloadEnd      ImageDownloadEvent = "image.download.end"
	ImageDownloadError    ImageDownloadEvent = "image.download.error"
	ImageDownloadProgress ImageDownloadEvent = "image.download.progress"
)

type ImageDownload struct {
	EventKind     ImageDownloadEvent `json:"event"`
	ImageURL      string             `json:"image_url"`
	ImageHeight   int64              `json:"image_height"`
	ImageWidth    int64              `json:"image_width"`
	ContentLength int64              `json:"content_length"`
	Downloaded    int64              `json:"downloaded"`
	Subreddit     string             `json:"subreddit"`
	PostURL       string             `json:"post_url"`
	PostName      string             `json:"post_name"`
	PostTitle     string             `json:"post_title"`
	Error         error              `json:"error"`
}

// Render the template.
func (im ImageDownload) Render(ctx context.Context, w io.Writer) error {
	panic("not implemented") // TODO: Implement
}

// Event returns the event name
func (im ImageDownload) Event() string {
	return "image.download"
}

// SerializeTo writes the event data to the writer.
//
// SerializeTo must not write multiple linebreaks (single linebreak is fine)
// in succession to the writer since it will mess up SSE events.
func (im ImageDownload) SerializeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(im)
}

type ImageDownloadSubreddit struct {
	ImageDownload
}

func (im ImageDownloadSubreddit) Event() string {
	return string(im.EventKind) + "." + im.Subreddit
}
