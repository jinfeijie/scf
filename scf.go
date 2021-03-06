package scf

import (
	"context"
	"github.com/spf13/cast"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
	"reflect"
	"sync"
)

type Handler func(ctx *Context) Reply

const Author = "jinfeijie"

type TrafficModeType int

const (
	TrafficModeUnknown = TrafficModeType(iota)
	TrafficModeGW
	TrafficModeServe
)

type RunModeType int

const (
	RunModeTypeUnknown = RunModeType(iota)
	RunModeDebug
	RunModeTest
	RunModeRelease
)

type Scf struct {
	trafficMode TrafficModeType
	runMode     RunModeType
	pool        sync.Pool
	*Router
}

func New() *Scf {
	scf := &Scf{
		trafficMode: TrafficModeServe,
		runMode:     RunModeDebug,
		pool:        sync.Pool{},
		Router:      NewRouter(),
	}
	scf.pool.New = func() interface{} {
		return scf.allocateContext()
	}
	return scf
}

func (scf *Scf) allocateContext() *Context {
	v := make(Params, 0, 1)
	return &Context{
		params: &v,
	}
}

// Use mw TODO
func (scf *Scf) Use(handlers ...Handler) {
	// TODO
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
	switch scf.trafficMode {
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

		headers := make(map[string]string)
		for key, value := range Map(Json(header)) {
			valueOf := reflect.ValueOf(value)
			if valueOf.Len() == 1 {
				headers[key] = cast.ToString(valueOf.Index(0).Interface())
			} else {
				headers[key] = cast.ToString(value)
			}
		}

		headers["Content-Type"] = contentType

		if scf.runMode == RunModeDebug {
			headers["X-Powered-By"] = Author
		}

		gw := events.APIGatewayResponse{
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
	scf.trafficMode = modeType
}
func (scf *Scf) GetTrafficMode() TrafficModeType {
	return scf.trafficMode
}

func (scf *Scf) Server(_ context.Context, r *Req) (Reply, error) {
	ctx := scf.pool.Get().(*Context)
	ctx.Reset()
	ctx.BuildCtx(r)
	defer scf.pool.Put(ctx)

	method := r.HttpMethod
	if route, exist := scf.trees[method]; exist {
		handle, ps, _ := route.getValue(r.Path, ctx.params)
		if handle == nil {
			return ctx.NotFound(), nil
		}

		if ps != nil {
			ctx.params = ps
		}
		return handle(ctx), nil
	}
	return ctx.NotFound(), nil
}

func (scf *Scf) SetMode(modeType RunModeType) {
	scf.runMode = modeType
}

func (scf *Scf) GetMode() RunModeType {
	return scf.runMode
}
