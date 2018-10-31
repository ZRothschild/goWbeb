package main

import (
	"github.com/kataras/iris"
)
//与嵌入式单页面应用程序相同但没有go-bindata，文件是"原始"存储在
//当前系统目录
var page = struct {
	Title string
}{"Welcome"}

func newApp() *iris.Application {
	app := iris.New()
	app.RegisterView(iris.HTML("./public", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Page", page)
		ctx.View("index.html")
	})
	//或者只是按原样提供index.html：
	// app.Get("/{f:path}", func(ctx iris.Context) {
	// 	ctx.ServeFile("index.html", false)
	// })
	assetHandler := app.StaticHandler("./public", false, false)
	//作为SPA的替代方案，您可以查看/routing/dynamic-path/root-wildcard
	app.SPA(assetHandler)
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	// http://localhost:8080/index.html
	// http://localhost:8080/app.js
	// http://localhost:8080/css/main.css
	app.Run(iris.Addr(":8080"))
}
