//包main显示了如何使用`WriteWithExpiration`
//基于“modtime”，如果If-Modified-Since的时间将于之前的对比，如果超出了refreshEvery的范围
//它会刷新内容，否则会让客户端（99.9％的浏览器） 处理缓存机制，它比iris.Cache更快，因为服务器端
//如果没有任何操作，无需将响应存储在内存中。
package main

import (
	"time"
	"github.com/kataras/iris"
)

const refreshEvery = 10 * time.Second

func main() {
	app := iris.New()
	app.Use(iris.Cache304(refreshEvery))
	// 等同于
	// app.Use(func(ctx iris.Context) {
	// 	now := time.Now()
	// 	if modified, err := ctx.CheckIfModifiedSince(now.Add(-refresh)); !modified && err == nil {
	// 		ctx.WriteNotModified()
	// 		return
	// 	}
	// 	ctx.SetLastModified(now)
	// 	ctx.Next()
	// })
	app.Get("/", greet)
	app.Run(iris.Addr(":8080"))
}

func greet(ctx iris.Context) {
	ctx.Header("X-Custom", "my  custom header")
	ctx.Writef("Hello World! %s", time.Now())
}
