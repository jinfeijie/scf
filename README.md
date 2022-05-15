#### 腾讯云SCF Go web 框架

支持通用路由的scf开发框架，降低scf开发过程中风格差异带来的心智成本。
> 使用方式与常规web开发框架无异。

* demo
```go
package main

import (
	"github.com/jinfeijie/scf"
)

func main() {
	s := scf.New()
	s.SetTrafficMode(scf.TrafficModeGW)
	s.SetMode(scf.RunModeDebug)
	{
		s.ANY("/scf", func(ctx *scf.Context) scf.Reply {
			return ctx.JSON("scf service", map[string]string{
				"version": "1.0.0",
			})
		})
	}
	s.Run()
}
```

