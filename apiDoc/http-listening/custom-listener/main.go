package main

import (
	"net"
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
	//创建任何自定义tcp侦听器，unix sock文件或tls tcp侦听器。
	l, err := net.Listen("tcp4", ":8080")
	if err != nil {
		panic(err)
	}
	//使用自定义侦听器
	app.Run(iris.Listener(l))
}