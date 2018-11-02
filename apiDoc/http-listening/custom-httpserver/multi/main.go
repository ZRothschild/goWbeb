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
	//注意：如果第一个动作是“go app.Run”，则不需要它。
	if err := app.Build(); err != nil {
		panic(err)
	}
	//启动侦听localhost：9090的辅助服务器。
	//如果您需要在同一个应用程序中使用多个服务器，请对Listen函数使用“go”关键字。
	// http://localhost:9090/
	// http://localhost:9090/mypath
	srv1 := &http.Server{Addr: ":9090", Handler: app}
	go srv1.ListenAndServe()
	println("Start a server listening on http://localhost:9090")
	//启动一个“second-secondary”服务器，监听localhost：5050。
	// http://localhost:5050/
	// http://localhost:5050/mypath
	srv2 := &http.Server{Addr: ":5050", Handler: app}
	go srv2.ListenAndServe()
	println("Start a server listening on http://localhost:5050")
	//注意：app.Run完全是可选的，我们已经使用app.Build构建了应用程序，
	//你可以改为创建一个新的http.Server。
	// http://localhost:8080/
	// http://localhost:8080/mypath
	app.Run(iris.Addr(":8080")) //在这里以后就不可以监听了
}