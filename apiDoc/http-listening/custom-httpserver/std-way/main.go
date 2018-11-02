package main

import (
	"net/http"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello from the server")
	})
	app.Get("/mypath", func(ctx iris.Context) {
		ctx.Writef("Hello from %s", ctx.Path())
	})
	//在自定义http.Server上使用'app'作为http.Handler之前调用.Build
	app.Build()
	//创建我们的自定义服务器并分配Handler/Router
	srv := &http.Server{Handler: app, Addr: ":8080"} //你必须设置Handler:app和Addr,请参阅“iris-way”,它会自动执行此操作
	// http://localhost:8080/
	// http://localhost:8080/mypath
	println("Start a server listening on http://localhost:8080")
	srv.ListenAndServe() // 等同于 app.Run(iris.Addr(":8080"))
	//注意：
	//根本不显示所有。 即使应用程序的配置允许，中断处理程序也是如此。
	//`.Run`是唯一一个关心这三者的函数。
	// 更多：
	//如果您需要在同一个应用程序中使用多个服务器，请参阅“multi”。
	//用于自定义侦听器：iris.Listener(net.Listener)或
	// iris.TLS(cert,key)或iris.AutoTLS()，请参阅“custom-listener”示例。
}