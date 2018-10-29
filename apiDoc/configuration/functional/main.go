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
	//当您想要更改某些配置字段时，这很好。
	//前缀：With，代码编辑器将帮助您浏览所有内容
	//配置选项，甚至没有参考文档的类型。
	app.Run(iris.Addr(":8080"), iris.WithoutStartupLog, iris.WithCharset("UTF-8"))

	// 在run之前加载:
	// app.Configure(iris.WithoutStartupLog, iris.WithCharset("UTF-8"))
	// app.Run(iris.Addr(":8080"))
}
