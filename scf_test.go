package scf

import (
	"context"
	"fmt"
	"testing"
)

func TestServerWarp(t *testing.T) {
	scf := New()
	scf.GET("/abababa/:c", func(ctx *Context) Reply {
		fmt.Println(ctx.path)
		fmt.Println(ctx.GetParams("c"))
		return ctx.HtmlCode(200, "123", "wwwwwww")
	})
	resp, err := scf.ServerWarp(context.Background(), &Req{
		Body:             "",
		HeaderParameters: nil,
		Headers: map[string]string{
			"x-api-requestid": "dsdsds",
		},
		HttpMethod:            "GET",
		Path:                  "/abababa/abc",
		PathParameters:        nil,
		QueryString:           nil,
		QueryStringParameters: nil,
		RequestContext:        nil,
		StageVariables:        nil,
	})
	t.Log(Json(resp), err)

}

func BenchmarkServerWarp(b *testing.B) {
	scf := New()
	scf.GET("/abababa/:c", func(ctx *Context) Reply {
		return ctx.HtmlCode(200, "123", "wwwwwww")
	})

	for i := 0; i < b.N; i++ {
		_, err := scf.ServerWarp(context.Background(), &Req{
			Body:             "",
			HeaderParameters: nil,
			Headers: map[string]string{
				"x-api-requestid": "dsdsds",
			},
			HttpMethod:            "GET",
			Path:                  "/abababa/abc",
			PathParameters:        nil,
			QueryString:           nil,
			QueryStringParameters: nil,
			RequestContext:        nil,
			StageVariables:        nil,
		})
		if err != nil {
			b.Error(err)
		}
		b.ReportAllocs()
	}
}
