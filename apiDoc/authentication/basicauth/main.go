package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/basicauth"
)

func newApp() *iris.Application {
	app := iris.New()
	authConfig := basicauth.Config{
		Users:   map[string]string{"myusername": "mypassword", "mySecondusername": "mySecondpassword"},
		Realm:   "Authorization Required", // 默认表示域 "Authorization Required"
		Expires: time.Duration(30) * time.Minute,
	}
	authentication := basicauth.New(authConfig)
	//作用范围 全局 app.Use(authentication) 或者 (app.UseGlobal 在Run之前)
	//作用范围 单个路由 app.Get("/mysecret", authentication, h)
	app.Get("/", func(ctx iris.Context) { ctx.Redirect("/admin") })
	//作用范围  Party
	needAuth := app.Party("/admin", authentication)
	{
		//http://localhost:8080/admin
		needAuth.Get("/", h)
		// http://localhost:8080/admin/profile
		needAuth.Get("/profile", h)
		// http://localhost:8080/admin/settings
		needAuth.Get("/settings", h)
	}
	return app
}

func main() {
	app := newApp()
	// open http://localhost:8080/admin
	app.Run(iris.Addr(":8080"))
}

func h(ctx iris.Context) {
	username, password, _ := ctx.Request().BasicAuth()
	//第三个参数因为中间件所以不需要判断其值，否则不会执行此处理程序
	ctx.Writef("%s %s:%s", ctx.Path(), username, password)
}
