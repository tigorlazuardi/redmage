package reddit

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alecthomas/units"
)

var ErrIdleTimeoutReached = errors.New("download idle timeout reached")

type ImageDownloadReader struct {
	OnProgress         func(downloaded int64, contentLength int64, err error)
	OnClose            func(closeErr error)
	IdleTimeout        time.Duration
	IdleSpeedThreshold units.MetricBytes

	errCancel      error
	cancelDebounce *time.Timer
	reader         io.ReadCloser
	contentLength  int64

	downloaded atomic.Int64

	deltastart time.Time
	deltavalue atomic.Int64

	end time.Time

	exit chan struct{}

	mu sync.Mutex
}

func (idr *ImageDownloadReader) WrapHTTPResponse(resp *http.Response) *http.Response {
	idr.reader = resp.Body
	idr.contentLength = resp.ContentLength
	idr.exit = make(chan struct{}, 1)
	resp.Body = idr
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-idr.exit:
				return
			case <-ticker.C:
				idr.checkSpeed()
			}
		}
	}()
	return resp
}

func (idr *ImageDownloadReader) checkSpeed() {
	now := time.Now()
	if idr.deltastart.IsZero() {
		idr.deltastart = now
	}

	if idr.cancelDebounce == nil {
		idr.cancelDebounce = time.AfterFunc(idr.IdleTimeout, func() {
			idr.mu.Lock()
			defer idr.mu.Unlock()
			idr.errCancel = ErrIdleTimeoutReached
		})
	}

	if now.Sub(idr.deltastart) < time.Second {
		return
	}
	idr.deltastart = now

	delta := idr.deltavalue.Load()

	if delta >= idr.IdleSpeedThreshold {
		idr.deltavalue.Store(0)
		idr.cancelDebounce.Stop()
		idr.cancelDebounce = nil
	}
}

func (idr *ImageDownloadReader) Read(p []byte) (n int, err error) {
	n, err = idr.reader.Read(p)

	idr.deltavalue.Add(int64(n))
	newd := idr.downloaded.Add(int64(n))
	if idr.OnProgress != nil {
		idr.OnProgress(newd, idr.contentLength, err)
	}

	idr.mu.Lock()
	if idr.errCancel != nil {
		idr.mu.Unlock()
		idr.OnProgress(newd, idr.contentLength, idr.errCancel)
		return n, idr.errCancel
	}
	idr.mu.Unlock()
	return n, err
}

func (idr *ImageDownloadReader) Close() error {
	idr.exit <- struct{}{}
	err := idr.reader.Close()
	if idr.OnClose != nil {
		idr.OnClose(err)
	}
	return err
}
