#### 腾讯云SCF Go web 框架

简单的go web框架

* demo
```go
package main

import (
	"github.com/jinfeijie/scf"
)

func main() {
	s := scf.New()
	s.SetTrafficMode(scf.TrafficModeGW)
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

