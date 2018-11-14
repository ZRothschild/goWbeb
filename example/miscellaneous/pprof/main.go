//go中有pprof包来做代码的性能监控
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/pprof"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1> Please click <a href='/debug/pprof'>here</a>")
	})
	app.Any("/debug/pprof/{action:path}", pprof.New())
	app.Run(iris.Addr(":8080"))
}