# New Relic

### [New Relic](https://github.com/newrelic/go-agent)介绍

`New Relic Go Agent`允许您使用`New Relic`监控`Go`应用程序。 它可以帮助您跟踪事务，出站请求，数据库调用以及`Go`应用程序行为的其他部分，
并提供垃圾收集，`goroutine`活动和内存使用的运行概述。

### 示例 `main.go`

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/newrelic"
)

func main() {
	app := iris.New()
	config := newrelic.Config("APP_SERVER_NAME", "NEWRELIC_LICENSE_KEY")
	m, err := newrelic.New(config)
	if err != nil {
		app.Logger().Fatal(err)
	}
	app.Use(m.ServeHTTP)
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("success!\n")
	})
	app.Run(iris.Addr(":8080"))
}
```