package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	i := 0
	//让我们在下一个请求时模拟panic
	app.Get("/", func(ctx iris.Context) {
		i++
		if i%2 == 0 {
			panic("a panic here")
		}
		ctx.Writef("Hello, refresh one time more to get panic!")
	})
	// http://localhost:8080, 刷新5-6次
	app.Run(iris.Addr(":8080"))
}
//注意： app := iris.Default()而不是iris.New()自动使用恢复中间件。