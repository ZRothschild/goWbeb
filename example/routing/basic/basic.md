# `route`基础用法
## 目录结构
> 主目录`basic`
```html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")
	//注册一个自定处理处理 hhttp 状态为404未找到路径错误的路由处理函数
	//触发条件是路由没有被找到，或通过手动调用 ctx.StatusCode(iris.StatusNotFound)
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
	// GET -> HTTP 的请求方法
	// / -> /路由名称
	// func(ctx iris.Context) -> 路由处理函数
	//
	//第三个可变参数应该包含一个或多个路由处理函数，他们将被顺序执行
	//例子如下:
	app.Handle("GET", "/", func(ctx iris.Context) {
		// 可以参考 https://github.com/kataras/iris/wiki/Routing-context-methods
		//详细介绍了所有 context 的所有可用方法不仅仅是 ctx.Path()
		ctx.HTML("Hello from " + ctx.Path()) // Hello from /
	})
	app.Get("/home", func(ctx iris.Context) {
		ctx.Writef(`Same as app.Handle("GET", "/", [...])`)
	})
	//同一路径中的不同路径参数类型。
	app.Get("/u/{p:path}", func(ctx iris.Context) {
		ctx.Writef(":string, :int, :uint, :alphabetical and :path in the same path pattern.")
	})
	app.Get("/u/{username:string}", func(ctx iris.Context) {
		ctx.Writef("before username (string), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("username (string): %s", ctx.Params().Get("username"))
	})
	app.Get("/u/{id:int}", func(ctx iris.Context) {
		ctx.Writef("before id (int), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("id (int): %d", ctx.Params().GetIntDefault("id", 0))
	})
	app.Get("/u/{uid:uint}", func(ctx iris.Context) {
		ctx.Writef("before uid (uint), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("uid (uint): %d", ctx.Params().GetUintDefault("uid", 0))
	})
	app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
		ctx.Writef("before firstname (alphabetical), current route name: %s\n", ctx.RouteName())
		ctx.Next()
	}, func(ctx iris.Context) {
		ctx.Writef("firstname (alphabetical): %s", ctx.Params().Get("firstname"))
	})
	/*
		/u/some/path/here 对应 :path
		/u/abcd 对应 :alphabetical (如果不写 :alphabetical 默认是 :string)
		/u/42 对应 :uint (如果不写 :uint 默认是 :int)
		/u/-1 对应 :int (如果不写 :int 默认是 :string)
		/u/abcd123 对应 :string
	*/
	// Pssst, don't forget dynamic-path example for more "magic"!
	// Pssst，别忘了使用动态路径示例获得更意外惊喜
	app.Get("/api/users/{userid:uint64 min(1)}", func(ctx iris.Context) {
		userID, err := ctx.Params().GetUint64("userid")
		if err != nil {
			ctx.Writef("error while trying to parse userid parameter," +
				"this will never happen if :uint64 is being used because if it's not a valid uint64 it will fire Not Found automatically.")
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
		ctx.JSON(map[string]interface{}{
			//当然，您可以传递任何自定义的结构化go值。
			"user_id": userID,
		})
	})
	// app.Post("/", func(ctx iris.Context){}) -> for POST http method.
	// app.Put("/", func(ctx iris.Context){})-> for "PUT" http method.
	// app.Delete("/", func(ctx iris.Context){})-> for "DELETE" http method.
	// app.Options("/", func(ctx iris.Context){})-> for "OPTIONS" http method.
	// app.Trace("/", func(ctx iris.Context){})-> for "TRACE" http method.
	// app.Head("/", func(ctx iris.Context){})-> for "HEAD" http method.
	// app.Connect("/", func(ctx iris.Context){})-> for "CONNECT" http method.
	// app.Patch("/", func(ctx iris.Context){})-> for "PATCH" http method.
	// app.Any("/", func(ctx iris.Context){}) for all http methods.

	//相同的路由可以对应多个不同的http 请求方法
	//您可以使用以下命令捕获任何路由创建错误：
	//路线， err := app.Get(...)
	//为路由设置名称：route.Name =“ myroute”

	//您还可以按路径前缀对路由进行分组，共享中间件和完成的需要处理的动作。

	adminRoutes := app.Party("/admin", adminMiddleware)

	adminRoutes.Done(func(ctx iris.Context) { // Done 在ctx.Next()后面才会调用，所以下面 / 路由会调用一下ctx.Next()
		ctx.HTML("<h1>Hello from fdkjin/</h1>")
		ctx.Application().Logger().Infof("response sent to " + ctx.Path())
	})
	// adminRoutes.Layout("/views/layouts/admin.html")
	// 为这些路由设置视图布局，请参view图示例。
	// GET: http://localhost:8080/admin
	adminRoutes.Get("/", func(ctx iris.Context) {
		// [...]
		ctx.StatusCode(iris.StatusOK) // 默认是 200 == iris.StatusOK
		ctx.HTML("<h1>Hello from admin/</h1>")
		ctx.Next() // 为了执行路由组的 Done" Handler(s)
})
	// GET: http://localhost:8080/admin/login
	adminRoutes.Get("/login", func(ctx iris.Context) {
		// [...]
	})
	// POST: http://localhost:8080/admin/login
	adminRoutes.Post("/login", func(ctx iris.Context) {
		// [...]
	})
	// 子域名比上面更容易, 执行要在host localhost或127.0.0.1
	// unix 路径在 /etc/hosts  windows 路径在 C:/windows/system32/drivers/etc/hosts
	v1 := app.Party("v1.")
	{ //花括号是可选的，它只是样式的一种，以可视方式对路线进行分组。
		// http://v1.localhost:8080
		//注意：对于版本特定的功能，请改用_examples /versioning。
		v1.Get("/", func(ctx iris.Context) {
			ctx.HTML(`Version 1 API. go to <a href="/api/users">/api/users</a>`)
		})
		usersAPI := v1.Party("/api/users")
		{
			// http://v1.localhost:8080/api/users
			usersAPI.Get("/", func(ctx iris.Context) {
				ctx.Writef("All users")
			})
			// http://v1.localhost:8080/api/users/42
			usersAPI.Get("/{userid:int}", func(ctx iris.Context) {
				ctx.Writef("user with id: %s", ctx.Params().Get("userid"))
			})
		}
	}
	// 通配符匹配子域名
	wildcardSubdomain := app.Party("*.")
	{
		wildcardSubdomain.Get("/", func(ctx iris.Context) {
			ctx.Writef("Subdomain can be anything, now you're here from: %s", ctx.Subdomain())
		})
	}
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	// http://localhost:8080/home
	// http://localhost:8080/api/users/42
	// http://localhost:8080/admin
	// http://localhost:8080/admin/login
	//
	// http://localhost:8080/api/users/0
	// http://localhost:8080/api/users/blabla
	// http://localhost:8080/wontfound
	//
	// http://localhost:8080/u/abcd
	// http://localhost:8080/u/42
	// http://localhost:8080/u/-1
	// http://localhost:8080/u/abcd123
	// http://localhost:8080/u/some/path/here
	//
	//  修改host情况下:
	//  http://v1.localhost:8080
	//  http://v1.localhost:8080/api/users
	//  http://v1.localhost:8080/api/users/42
	//  http://anything.localhost:8080
	app.Run(iris.Addr(":8080"))
}

func adminMiddleware(ctx iris.Context) {
	// [...]
	ctx.HTML("Hello from data" + ctx.Path())
	ctx.Next()  //移至下一个处理程序，如果有任何身份验证逻辑，则不要调用这方法。
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>.")
}

//注意：
//路由参数名称仅包含字母，符号则 _ ，数字将不被允许
//如果无法注册路由，则应用会在没有任何警告的情况下崩溃

//请参阅“file-server/single-page-application”以了解另一个功能“ WrapRouter”的工作方式。

```