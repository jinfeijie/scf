package scf

import (
	"context"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"net/http"
	"reflect"
	"sync"
)

type Handler func(ctx *Context) Reply

var router = make(map[string]map[string]Handler)

type TrafficModeType int

const (
	TrafficModeUnknown = TrafficModeType(iota)
	TrafficModeGW
	TrafficModeServe
)

type Scf struct {
	TrafficMode TrafficModeType
	pool        sync.Pool
}

func New() *Scf {
	scf := &Scf{
		TrafficMode: TrafficModeServe,
		pool:        sync.Pool{},
	}
	scf.pool.New = func() interface{} {
		return scf.allocateContext()
	}
	return scf
}

func (scf *Scf) allocateContext() *Context {
	return &Context{}
}

func (scf *Scf) Use(handlers ...Handler) {

}

func (scf *Scf) Route(method, path string, handler Handler) {
	if router[method] == nil {
		router[method] = make(map[string]Handler)
	}
	router[method][path] = handler
}

func (scf *Scf) GET(path string, handler Handler) {
	scf.Route(http.MethodGet, path, handler)
}

func (scf *Scf) POST(path string, handler Handler) {
	scf.Route(http.MethodPost, path, handler)
}

func (scf *Scf) OPTIONS(path string, handler Handler) {
	scf.Route(http.MethodOptions, path, handler)
}

func (scf *Scf) PUT(path string, handler Handler) {
	scf.Route(http.MethodPut, path, handler)
}

func (scf *Scf) DELETE(path string, handler Handler) {
	scf.Route(http.MethodDelete, path, handler)
}

func (scf *Scf) HEAD(path string, handler Handler) {
	scf.Route(http.MethodHead, path, handler)
}

func (scf *Scf) ANY(path string, handler Handler) {
	scf.GET(path, handler)
	scf.POST(path, handler)
	scf.PUT(path, handler)
	scf.DELETE(path, handler)
	scf.HEAD(path, handler)
}

func (scf *Scf) Run() {
	cloudfunction.Start(scf.ServerWarp)
}

type Req struct {
	Body                  string                 `json:"body"`
	HeaderParameters      map[string]string      `json:"headerParameters"`
	Headers               map[string]string      `json:"headers"`
	HttpMethod            string                 `json:"httpMethod"`
	Path                  string                 `json:"path"`
	PathParameters        map[string]interface{} `json:"pathParameters"`
	QueryString           map[string]string      `json:"queryString"`
	QueryStringParameters map[string]string      `json:"queryStringParameters"`
	RequestContext        map[string]interface{} `json:"requestContext"`
	StageVariables        map[string]string      `json:"stageVariables"`
}

func (scf *Scf) ServerWarp(ctx context.Context, r *Req) (resp map[string]interface{}, err error) {
	rly, err := scf.Server(ctx, r)
	contentType := rly.ContentType
	rly.ContentType = ""
	header := rly.Header
	rly.Header = nil
	resp = make(map[string]interface{})
	resp = Map(Json(rly))
	switch scf.TrafficMode {
	case TrafficModeGW:
		body := ""
		switch contentType {
		case ContentTypeJson:
			body = Json(rly)
		case ContentTypeHtml:
			_body, ok := rly.Data.(string)
			if ok {
				body = _body
			}
		}

		headers := make(map[string]interface{})
		for key, value := range Map(Json(header)) {
			valueOf := reflect.ValueOf(value)
			if valueOf.Len() == 1 {
				headers[key] = valueOf.Index(0).Interface()
			} else {
				headers[key] = value
			}
		}

		headers["Content-Type"] = contentType

		gw := GWReply{
			IsBase64Encoded: false,
			StatusCode:      rly.Code,
			Headers:         headers,
			Body:            body,
		}
		resp = Map(Json(gw))
	}

	return resp, err
}

func (scf *Scf) SetTrafficMode(modeType TrafficModeType) {
	scf.TrafficMode = modeType
}

func (scf *Scf) Server(_ context.Context, r *Req) (Reply, error) {
	ctx := scf.pool.Get().(*Context)
	ctx.Reset()
	ctx.BuildCtx(r)
	defer scf.pool.Put(ctx)

	method := r.HttpMethod
	route, exist := router[method]
	if !exist {
		return ctx.NotFound(), nil
	}

	handler, handlerExist := route[r.Path]
	if !handlerExist {
		return ctx.NotFound(), nil
	}
	return handler(ctx), nil
}
