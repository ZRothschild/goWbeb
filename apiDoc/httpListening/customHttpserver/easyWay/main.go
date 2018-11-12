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
	//这里有任何自定义字段 Handler和ErrorLog自动设置到服务器
	srv := &http.Server{Addr: ":8080"}
	// http://localhost:8080/
	// http://localhost:8080/mypath
	app.Run(iris.Server(srv)) // 等同于 app.Run(iris.Addr(":8080"))
	// 更多：
	//如果您需要在同一个应用程序中使用多个服务器，请参阅“multi”。
	//用于自定义listener：iris.Listener（net.Listener）或
	// iris.TLS(cert，key)或iris.AutoTLS()，请参阅“custom-listener”示例。
}