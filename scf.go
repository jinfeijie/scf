package scf

import (
	"context"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"net/http"
)

type Handler func(ctx *Context) Reply

var router = make(map[string]map[string]Handler)

type Scf struct {
}

func New() *Scf {
	return &Scf{}
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
	cloudfunction.Start(scf.Server)
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

func (scf *Scf) Server(_ context.Context, r *Req) (Reply, error) {
	ctx := NewContext(r)
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
