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
	//如果iris.YAML的输入字符串参数为“〜”，则它从主目录加载配置
	//并且可以在许多iris实例之间共享。
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
	// 在run之前加载:
	// app.Configure(iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
	// app.Run(iris.Addr(":8080"))
}