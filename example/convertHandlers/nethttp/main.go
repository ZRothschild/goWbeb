package main

import (
	"net/http"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	//FromStd将原生http.Handler,http.HandlerFunc和func(w，r，next)转换为context.Handler
	irisMiddleware := iris.FromStd(nativeTestMiddleware)
	app.Use(irisMiddleware)
	// 请求 GET: http://localhost:8080/
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("Home")
	})
	// 请求 GET: http://localhost:8080/ok
	app.Get("/ok", func(ctx iris.Context) {
		ctx.HTML("<b>Hello world!</b>")
	})
	// http://localhost:8080
	// http://localhost:8080/ok
	app.Run(iris.Addr(":8080"))
}

func nativeTestMiddleware(w http.ResponseWriter, r *http.Request) {
	println("Request path: " + r.URL.Path)
}
//如果要使用自定义上下文转换自定义处理程序，请查看routing/custom-context
//一个context.Handler。