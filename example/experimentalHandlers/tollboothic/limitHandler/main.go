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
//阅读更多信息：https：//github.com/didip/tollbooth
