package reddit

import (
	"io"
	"net/http"
)

type ProgressReader struct {
	OnProgress func(downloaded int64, contentLength int64, err error)
	OnClose    func(closeErr error)

	reader        io.ReadCloser
	contentLength int64
	downloaded    int64
}

func (progressReader *ProgressReader) WrapHTTPResponse(resp *http.Response) *http.Response {
	progressReader.reader = resp.Body
	progressReader.contentLength = resp.ContentLength
	resp.Body = progressReader
	return resp
}

func (progressReader *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = progressReader.reader.Read(p)
	progressReader.downloaded += int64(n)
	if progressReader.OnProgress != nil {
		progressReader.OnProgress(progressReader.downloaded, progressReader.contentLength, err)
	}
	return n, err
}

func (progressReader *ProgressReader) Close() error {
	err := progressReader.reader.Close()
	if progressReader.OnClose != nil {
		progressReader.OnClose(err)
	}
	return err
}
