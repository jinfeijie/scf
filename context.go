package scf

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Context struct {
	m              sync.Map
	Headers        http.Header
	request        map[string]interface{}
	Method         string
	requestId      string
	serviceId      string
	path           string
	Body           io.ReadCloser
	Query          url.Values
	Param          url.Values
	StageVariables map[string]string
}

func NewContext(r *Req) *Context {
	c := &Context{}
	c.BuildCtx(r)
	return c
}

func (c *Context) Value(key interface{}) interface{} {
	value, exist := c.m.Load(key)
	if exist {
		return value
	}
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Set(key interface{}, value interface{}) {
	c.m.Store(key, value)
}

func (c *Context) Get(key interface{}) (interface{}, bool) {
	return c.m.Load(key)
}

func (c *Context) Del(key interface{}) {
	c.m.Delete(key)
}

func (c *Context) AllHeaders() http.Header {
	return c.Headers
}
func (c *Context) Header(key string) string {
	return c.Headers.Get(key)
}

func (c *Context) Request(key string) interface{} {
	if val, exist := c.request[key]; exist {
		return val
	}
	return nil
}

func (c *Context) RequestId() string {
	return c.requestId
}

func (c *Context) ServiceId() string {
	return c.serviceId
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) GetBody() io.ReadCloser {
	return c.Body
}

func (c *Context) GetQuery(key string) string {
	return c.Query.Get(key)
}

func (c *Context) BuildCtx(r *Req) {
	c.Method = r.HttpMethod
	c.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(r.Body)))
	c.StageVariables = r.StageVariables

	if requestId, ok := r.RequestContext["requestId"].(string); ok {
		c.requestId = requestId
	}
	c.serviceId = r.StageVariables["stage"]
	c.path = r.Path

	for k, v := range r.RequestContext {
		c.Set(k, v)
	}

	c.Headers = http.Header{}
	for k, v := range r.Headers {
		c.Headers.Add(k, v)
	}

	c.Query = url.Values{}
	for k, v := range r.QueryString {
		c.Query.Set(k, v)
	}

	c.Param = url.Values{}
	for k, v := range r.QueryStringParameters {
		c.Param.Set(k, v)
	}
}

func (c *Context) NotFound() Reply {
	return Reply{
		Msg:       "404 page not found",
		Data:      nil,
		RequestId: c.requestId,
		UnixTime:  time.Now().Unix(),
	}
}

func (c *Context) JSON(msg string, data interface{}) Reply {
	return Reply{
		Msg:       msg,
		Data:      data,
		RequestId: c.requestId,
		UnixTime:  time.Now().Unix(),
	}
}
