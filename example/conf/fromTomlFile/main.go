package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})
	// [...]
	//当你有两个配置时很好，一个用于开发，另一个用于生产用途。
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./configs/iris.tml")))

	// 会run之前加载:
	// app.Configure(iris.WithConfiguration(iris.TOML("./configs/iris.tml")))
	// app.Run(iris.Addr(":8080"))
}
