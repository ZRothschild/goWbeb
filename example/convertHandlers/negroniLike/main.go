package main

import (
	"net/http"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	//FromStd将原生http.Handler,http.HandlerFunc和func(w，r，next)转换为context.Handler
	irisMiddleware := iris.FromStd(negronilikeTestMiddleware)
	app.Use(irisMiddleware)
	// 请求 GET: http://localhost:8080/
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1> Home </h1>")
		//这会打印错误
		//这个路由的处理程序永远不会被执行，因为中间件的标准没有通过
	})
	// 请求 GET: http://localhost:8080/ok
	app.Get("/ok", func(ctx iris.Context) {
		ctx.Writef("Hello world!")
		// 这将打印"OK. Hello world!"
	})
	// http://localhost:8080
	// http://localhost:8080/ok
	app.Run(iris.Addr(":8080"))
}

func negronilikeTestMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == "/ok" && r.Method == "GET" {
		w.Write([]byte("OK. "))
		next(w, r) // 中间件
		return
	}
	//否则打印错误
	w.WriteHeader(iris.StatusBadRequest)
	w.Write([]byte("Bad request"))
}
//如果要使用自定义上下文转换自定义处理程序，请查看routing/custom-context
//一个context.Handler