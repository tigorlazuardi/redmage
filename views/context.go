package views

import (
	"net/http"
	"net/url"

	"github.com/tigorlazuardi/redmage/config"
)

type Context struct {
	Config  *config.Config
	Request *http.Request
	Query   url.Values
}

func (c *Context) AppendQuery(keyValue ...string) string {
	query := c.Request.URL.Query()
	for i := 0; i < len(keyValue); i += 2 {
		query.Add(keyValue[i], keyValue[i+1])
	}
	return query.Encode()
}

func NewContext(config *config.Config, request *http.Request) *Context {
	return &Context{
		Config:  config,
		Request: request,
		Query:   request.URL.Query(),
	}
}
