# `Tollbooth`

## `Tollbooth`介绍

这是一个限制`HTTP`请求次数的中间件，该库被认为已完成。主要版本更改是向后不兼容的.`v2.0.0`简化了旧`API`

1. `v1.0.0`:此版本维护旧`API`，但所有第三方模块都移动到他们自己的`repo`。
2. `v2.x.x`:全新的`API`，用于代码清理，线程安全和自动过期的数据结构。
3. `v3.x.x`:显然我们一直在使用`golang.org/x/time/rate`。 见问题`＃48`。它始终限制每`1`秒的`X`数。持续时间不可更改，因此将
`TTL`传递给`tollbooth`是没有意义的。

## 示例代码

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/didip/tollbooth"
	"github.com/iris-contrib/middleware/tollboothic"
)

// $ go get github.com/didip/tollbooth
// $ go run main.go

func main() {
	app := iris.New()
	limiter := tollbooth.NewLimiter(1, nil)
	//或使用可过期的token buckets创建限制器
	//此设置表示：
	//创建1 request/second限制器和
	//其中的每个token bucket将在最初设置后1小时到期。
	// limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	app.Get("/", tollboothic.LimitHandler(limiter), func(ctx iris.Context) {
		ctx.HTML("<b>Hello, world!</b>")
	})
	app.Run(iris.Addr(":8080"))
}
//阅读更多信息：https://github.com/didip/tollbooth
```