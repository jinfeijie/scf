package scf

import (
	"bytes"
	"encoding/json"
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
	traceId        string
	serviceId      string
	path           string
	Body           io.ReadCloser
	Query          url.Values
	queryMap       map[string]string
	Param          url.Values
	paramMap       map[string]string
	StageVariables map[string]string
	RealIp         string

	// RESPONSE
	responseHeader http.Header
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

	if requestId, ok := r.Headers["x-api-requestid"]; ok {
		c.requestId = requestId
	}

	if traceId, ok := r.Headers["x-b3-traceid"]; ok {
		c.traceId = traceId
	}

	if serviceId, ok := r.RequestContext["serviceId"].(string); ok {
		c.serviceId = serviceId
	}

	if realIp, ok := r.RequestContext["sourceIp"].(string); ok {
		c.RealIp = realIp
	}

	c.path = r.Path

	for k, v := range r.RequestContext {
		c.Set(k, v)
	}

	c.Headers = http.Header{}
	for k, v := range r.Headers {
		c.Headers.Add(k, v)
	}

	c.Query = url.Values{}
	c.queryMap = r.QueryString
	for k, v := range r.QueryString {
		c.Query.Set(k, v)
	}

	c.Param = url.Values{}
	c.paramMap = r.QueryStringParameters
	for k, v := range r.QueryStringParameters {
		c.Param.Set(k, v)
	}

	c.responseHeader = http.Header{}
}

//func (c *Context) GWWarp(reply Reply) GWReply {
//	return GWReply{
//		IsBase64Encoded: false,
//		StatusCode:      reply.Code,
//		Headers: map[string]interface{}{
//			"Content-Type": "application/json; charset=utf-8",
//		},
//		Body: Json(reply),
//	}
//}

func (c *Context) NotFound() Reply {
	return c.JSONCode(404, "404 page not found", "["+c.Method+"] "+c.path)
}

func (c *Context) JSON(msg string, data interface{}) Reply {
	return c.JSONCode(200, msg, data)
}

func (c *Context) Error(msg string, data interface{}) Reply {
	return c.JSONCode(500, msg, data)
}

func (c *Context) JSONCode(code int, msg string, data interface{}) Reply {
	return c.Reply(ContentTypeJson, code, msg, data)
}

func (c *Context) HtmlCode(code int, msg string, data string) Reply {
	return c.Reply(ContentTypeHtml, code, msg, data)
}

func (c *Context) Reply(contentType string, code int, msg string, data interface{}) Reply {
	return Reply{
		ContentType: contentType,
		Header:      c.responseHeader,
		Code:        code,
		Msg:         msg,
		Data:        data,
		RequestId:   c.requestId,
		TraceId:     c.traceId,
		UnixTime:    time.Now().Unix(),
	}
}

func (c *Context) ReplyHeader(key string, value string) {
	c.responseHeader.Add(key, value)
}

func (c *Context) BodyBind(data interface{}) error {
	all, err := ioutil.ReadAll(c.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(all, &data); err != nil {
		return err
	}
	return nil
}

func (c *Context) QueryBind(data interface{}) error {
	if err := json.Unmarshal([]byte(Json(c.queryMap)), &data); err != nil {
		return err
	}
	return nil
}

func (c *Context) ParamsBind(data interface{}) error {
	if err := json.Unmarshal([]byte(Json(c.paramMap)), &data); err != nil {
		return err
	}
	return nil
}
