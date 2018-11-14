package main

import (
	"time"
	"./routes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero" // <- 导入
	"github.com/kataras/iris/sessions"
)

func main() {
	app := iris.New()
	sessionManager := sessions.New(sessions.Config{
		Cookie:       "site_session_id",
		Expires:      60 * time.Minute,
		AllowReclaim: true,
	})
	//注册
	//动态依赖关系，比如* sessions.Session，来自`sessionManager.Start(ctx)*sessions.Session` < - 它接受一个Context并返回
	// something - >这称为动态请求 - 时间依赖，并且某些东西可以作为输入参数用于处理程序，
	//没有关于依赖项数量限制，每个处理程序将在服务器运行之前构建一次，并且它将仅使用它所需的依赖项。
	hero.Register(sessionManager.Start)
	//将任何函数转换为iris Handler，使用独特的Iris超快速依赖注入来解析它们的输入参数
	//用于服务或动态依赖，例如* sessions.Session，来自sessionManager.Start(ctx)* sessions.Session)< - 它接受一个Context并返回
	// 某些东西->这称为动态请求时依赖性。
	indexHandler := hero.Handler(routes.Index)
	// Method: GET
	// Path: http://localhost:8080
	app.Get("/", indexHandler)
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}