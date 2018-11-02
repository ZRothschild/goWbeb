package main

import (
	"context"
	"time"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1>Hello, try to refresh the page after ~10 secs</h1>")
	})
	app.Logger().Info("Wait 10 seconds and check your terminal again")
	//在这里模拟一个关机动作......
	go func() {
		<-time.After(10 * time.Second)
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		//关闭所有主机，这将通知我们已注册的回调
		//在`configureHost` func中。
		app.Shutdown(ctx)
	}()
	// app.ConfigureHost(configureHost) - >或将“configureHost”作为`app.Addr`参数传递，结果相同。
	//像往常一样启动服务器，唯一的区别就是
	//我们正在添加第二个（可选）功能
	//配置刚刚创建的主机管理。
	// http：// localhost：8080
	//等待10秒钟并检查您的终端
	app.Run(iris.Addr(":8080", configureHost), iris.WithoutServerError(iris.ErrServerClosed))
	/*
	或者对于简单的情况，您可以使用：
	iris.RegisterOnInterrupt用于CTRL/CMD+C和OS事件的全局捕获。
	查看“graceful-shutdown”示例了解更多信息。
	*/
}

func configureHost(su *iris.Supervisor) {
	//这里我们可以完全访问将要创建的主机
	//在`app.Run`函数或`NewHost`。
	//我们在这里注册一个关闭“事件”回调：
	su.RegisterOnShutdown(func() {
		println("server is closed")
	})
	// su.RegisterOnError
	// su.RegisterOnServe
}