package htmx

import (
	"bytes"
	"time"

	"github.com/inhies/go-bytesize"
	"github.com/labstack/echo/v5"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/templates/components"
)

type byteSizeParser struct {
	bytesize.ByteSize
}

type secondDuration struct {
	time.Duration
}

func (s *secondDuration) UnmarshalText(text []byte) error {
	if len(bytes.TrimSpace(text)) == 0 {
		*s = secondDuration{0}
		return nil
	}
	dur, err := time.ParseDuration(string(text) + "s")
	if err != nil {
		return err
	}
	s.Duration = dur
	return nil
}

func (b *byteSizeParser) UnmarshalText(text []byte) error {
	if len(bytes.TrimSpace(text)) == 0 {
		*b = byteSizeParser{0}
		return nil
	}
	text = append(text, 'K', 'B')
	return b.ByteSize.UnmarshalText(text)
}

func (r *Routes) ConfigUpdateDownload(c echo.Context) error {
	type UpdateDownloadSchema struct {
		config.Download
		ConnectionTimeout     secondDuration `schema:"connection_timeout"`
		DownloadIdleTimeout   secondDuration `schema:"download_idle_timeout"`
		DownloadIdleThreshold byteSizeParser `schema:"download_idle_threshold"`
	}

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := c.Request().ParseForm(); err != nil {
		c.Response().WriteHeader(400)
		return components.
			ErrorToast("Error reading form data: %s", err.Error()).
			Render(c.Request().Context(), c.Response())
	}

	s := UpdateDownloadSchema{}
	if err := r.Schema.Decode(&s, c.Request().PostForm); err != nil {
		c.Response().WriteHeader(400)
		return components.
			ErrorToast("Error parsing form data: %s", err.Error()).
			Render(c.Request().Context(), c.Response())
	}

	s.Download.DownloadIdleThreshold = s.DownloadIdleThreshold.ByteSize
	s.Download.ConnectionTimeout = s.ConnectionTimeout.Duration
	s.Download.DownloadIdleTimeout = s.DownloadIdleTimeout.Duration

	r.Config.Download = s.Download
	if err := r.Config.Sync(); err != nil {
		c.Response().WriteHeader(500)
		return components.
			ErrorToast(err.Error()).
			Render(c.Request().Context(), c.Response())
	}

	c.Response().WriteHeader(200)
	return components.
		SuccessToast("Config updated successfully").
		Render(c.Request().Context(), c.Response())
}
