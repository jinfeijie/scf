package scf

import "net/http"

type Reply struct {
	Code        int         `json:"code"`
	Msg         string      `json:"msg"`
	Data        interface{} `json:"data"`
	ContentType string      `json:"content_type,omitempty"`
	Header      http.Header `json:"header,omitempty"`
	RequestId   string      `json:"request_id"`
	TraceId     string      `json:"trace_id"`
	UnixTime    int64       `json:"unix_time"`
}

type GWReply struct {
	IsBase64Encoded bool                   `json:"isBase64Encoded"`
	StatusCode      int                    `json:"statusCode"`
	Headers         map[string]interface{} `json:"headers"`
	Body            string                 `json:"body"`
}

const (
	ContentTypeHtml = "text/html; charset=utf-8"
	ContentTypeJson = "application/json; charset=utf-8"
)
