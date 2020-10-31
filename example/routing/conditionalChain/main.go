package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	v1 := app.Party("/api/v1")
	myFilter := func(ctx iris.Context) bool {
		//请勿在生产环境中执行此操作，请使用会话或/和数据库调用等。
		ok, _ := ctx.URLParamBool("admin")
		return ok
	}
	onlyWhenFilter1 := func(ctx iris.Context) {
		ctx.Application().Logger().Infof("admin: %v", ctx.Params())
		ctx.HTML("<h1>Hello in</h1><br>")
		ctx.Next()
	}
	onlyWhenFilter2 := func(ctx iris.Context) {
		//您可以一直使用上一个求情处理方法存储的数据
		//执行类似ofc的操作。
		//
		//当前路由处理方法设置 ：ctx.Values().Set("is_admin", true)
		//下一个路由处理方法可以获取到上一个设置的值 ：isAdmin := ctx.Values().GetBoolDefault("is_admin", false)
		//
		//，但让我们简化一下：
		ctx.HTML("<h1>Hello Admin</h1><br>")
		ctx.Next()
	}
	// HERE:
	// It can be registered anywhere, as a middleware.
	// It will fire the `onlyWhenFilter1` and `onlyWhenFilter2` as middlewares (with ctx.Next())
	// if myFilter pass otherwise it will just continue the handler chain with ctx.Next() by ignoring
	// the `onlyWhenFilter1` and `onlyWhenFilter2`.
	// 这里：
	//它可以在任何地方注册为中间件。
	//它将触发`onlyWhenFilter1`和`onlyWhenFilter2`作为中间件（使用ctx.Next（））
	//如果myFilter通过，否则它将通过忽略ctx.Next（）继续处理程序链
	//`onlyWhenFilter1`和`onlyWhenFilter2`。

	//onlyWhenFilter1 function and onlyWhenFilter2 function only onlyWhenFilter2 functions are executed
	myMiddleware := iris.NewConditionalHandler(myFilter, onlyWhenFilter1, onlyWhenFilter2)
	v1UsersRouter := v1.Party("/users", myMiddleware)
	v1UsersRouter.Get("/", func(ctx iris.Context) {
		ctx.HTML("requested: <b>/api/v1/users</b>")
	})
	// http://localhost:8080/api/v1/users
	// http://localhost:8080/api/v1/users?admin=true
	app.Run(iris.Addr(":8000"))
}
