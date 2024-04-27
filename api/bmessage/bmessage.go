package bmessage

import (
	"encoding/json"

	"github.com/alecthomas/units"
)

type ImageMetadata struct {
	Kind   ImageKind
	URL    string
	Height int64
	Width  int64
}

type ImageKind int

func (im ImageKind) String() string {
	switch im {
	case KindThumbnail:
		return "Thumbnail"
	default:
		return "Image"
	}
}

func (im ImageKind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + im.String() + `"`), nil
}

const (
	KindImage ImageKind = iota
	KindThumbnail
)

type DownloadEvent int

func (do DownloadEvent) MarshalJSON() ([]byte, error) {
	return []byte(`"` + do.String() + `"`), nil
}

func (do DownloadEvent) String() string {
	switch do {
	case DownloadStart:
		return "DownloadStart"
	case DownloadProgress:
		return "DownloadProgress"
	case DownloadEnd:
		return "DownloadEnd"
	case DownloadError:
		return "DownloadError"
	default:
		return "Unknown"
	}
}

const (
	DownloadStart DownloadEvent = iota
	DownloadProgress
	DownloadEnd
	DownloadError
)

type ImageDownloadMessage struct {
	Event         DownloadEvent     `json:"event"`
	Metadata      ImageMetadata     `json:"metadata"`
	ContentLength units.MetricBytes `json:"content_length"`
	Downloaded    units.MetricBytes `json:"downloaded"`
	Subreddit     string            `json:"subreddit"`
	PostURL       string            `json:"post_url"`
	PostID        string            `json:"post_id"`
	Error         error             `json:"error"`
}

func (im ImageDownloadMessage) MarshalJSON() ([]byte, error) {
	type Alias ImageDownloadMessage
	type W struct {
		Alias
		Error string `json:"error"`
	}

	errMsg := ""
	if im.Error != nil {
		errMsg = im.Error.Error()
	}

	w := W{
		Alias: Alias(im),
		Error: errMsg,
	}

	return json.Marshal(w)
}
