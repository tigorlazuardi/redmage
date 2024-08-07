package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/a-h/templ"
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

// URLWithExtraQuery creates a query based from baseUrl with queries joined between
// current context and extraQueries.
//
// extraKeyValues is an alternating key-value pair.
func (c *Context) URLWithExtraQuery(baseUrl string, extraKeyValues ...string) templ.SafeURL {
	query := c.Request.URL.Query()
	for k := range query {
		if query.Get(k) == "" {
			delete(query, k)
		}
	}
	for i := 0; i < len(extraKeyValues); i += 2 {
		query.Set(extraKeyValues[i], extraKeyValues[i+1])
	}
	return templ.SafeURL(fmt.Sprintf("%s?%s", baseUrl, query.Encode()))
}

func (c *Context) JSONQuery() string {
	m := make(map[string]string, len(c.Query))
	for k := range c.Query {
		v := c.Query.Get(k)
		if v != "" {
			m[k] = v
		}
	}
	v, _ := json.Marshal(m)
	return string(v)
}

func NewContext(config *config.Config, request *http.Request) *Context {
	return &Context{
		Config:  config,
		Request: request,
		Query:   request.URL.Query(),
	}
}
