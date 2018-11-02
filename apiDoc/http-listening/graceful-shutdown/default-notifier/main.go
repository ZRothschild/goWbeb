package main

import (
	stdContext "context"
	"time"
	"github.com/kataras/iris"
)
//继续之前：
//正常关闭control+C/command+C或当发送的kill命令是ENABLED BY-DEFAULT。
//为了手动管理应用程序中断时要执行的操作，
//我们必须使用选项`WithoutInterruptHandler`禁用默认行为
//并注册一个新的中断处理程序(全局，所有主机)。
func main() {
	app := iris.New()
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		//关闭所有主机
		app.Shutdown(ctx)
	})
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(" <h1>hi, I just exist in order to see if the server is closed</h1>")
	})
	// http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}