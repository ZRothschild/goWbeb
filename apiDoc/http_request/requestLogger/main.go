package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func main() {
	app := iris.New()
	customLogger := logger.New(logger.Config{
		//状态显示状态代码
		Status: true,
		// IP显示请求的远程地址
		IP: true,
		//方法显示http方法
		Method: true,
		// Path显示请求路径
		Path: true,
		// Query将url查询附加到Path。
		Query: true,
		//Columns：true，
		// 如果不为空然后它的内容来自`ctx.Values(),Get("logger_message")
		//将添加到日志中。
		MessageContextKeys: []string{"logger_message"},
		//如果不为空然后它的内容来自`ctx.GetHeader（“User-Agent”）
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(customLogger)
	h := func(ctx iris.Context) {
		ctx.Writef("Hello from %s", ctx.Path())
	}
	app.Get("/", h)
	app.Get("/1", h)
	app.Get("/2", h)
	//因此，http错误有自己的处理程序
	//注册中间人应该手动完成。
	/*
		 app.OnErrorCode(404 ,customLogger, func(ctx iris.Context) {
			ctx.Writef("My Custom 404 error page ")
		 })
	*/
	//或捕获所有http错误:
	app.OnAnyErrorCode(customLogger, func(ctx iris.Context) {
		//这应该被添加到日志中，因为`logger.Config＃MessageContextKey`
		ctx.Values().Set("logger_message",
			"a dynamic message passed to the logs")
		ctx.Writef("My Custom error page")
	})
	// http://localhost:8080
	// http://localhost:8080/1
	// http://localhost:8080/2
	// http://lcoalhost:8080/notfoundhere
	//查看控制台上的输出
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}