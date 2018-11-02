package main

import (
	"net/url"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello from the SECURE server")
	})
	app.Get("/mypath", func(ctx iris.Context) {
		ctx.Writef("Hello from the SECURE server on path /mypath")
	})
	//启动一个新的服务器,监听:80并重定向,到安全地址，然后:
	target, _ := url.Parse("https://127.0.1:443")
	go host.NewProxy("127.0.0.1:80", target).ListenAndServe()
	//在端口443上启动服务器(HTTPS)这是一个阻塞功能
	app.Run(iris.TLS("127.0.0.1:443", "mycert.cert", "mykey.key"))
}