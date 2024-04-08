package views

import (
	"net/http"

	"github.com/tigorlazuardi/redmage/config"
)

type Context struct {
	Config  *config.Config
	Request *http.Request
}

func NewContext(config *config.Config, request *http.Request) *Context {
	return &Context{
		Config:  config,
		Request: request,
	}
}
