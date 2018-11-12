package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.ConfigureHost(func(host *iris.Supervisor) { // <- 重要
		//您可以使用某些主机的方法控制流或延迟某些内容：
		// host.RegisterOnError
		// host.RegisterOnServe
		host.RegisterOnShutdown(func() {
			app.Logger().Infof("Application shutdown on signal")
		})
	})
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1>Hello</h1>\n")
	})
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	/*
	对于默认信号中断事件，使用`iris.RegisterOnInterrupt`可以更简单地通知全局关闭。
	您甚至可以通过查看：“gracefulShutdown”示例进一步了解它。
	*/
}