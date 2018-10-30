package main

import (
	"github.com/kataras/iris"
	"github.com/casbin/casbin"
	cm "github.com/iris-contrib/middleware/casbin"
)

// $ go get github.com/casbin/casbin
// $ go run main.go
// Enforcer映射模型和casbin服务的策略，我们也在main_test上使用此变量。
var Enforcer = casbin.NewEnforcer("casbinmodel.conf", "casbinpolicy.csv")

func newApp() *iris.Application {
	casbinMiddleware := cm.New(Enforcer)
	app := iris.New()
	app.Use(casbinMiddleware.ServeHTTP)
	app.Get("/", hi)
	app.Get("/dataset1/{p:path}", hi) // p, alice, /dataset1/*, GET
	app.Post("/dataset1/resource1", hi)
	app.Get("/dataset2/resource2", hi)
	app.Post("/dataset2/folder1/{p:path}", hi)
	app.Any("/dataset2/resource1", hi)
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

func hi(ctx iris.Context) {
	ctx.Writef("Hello %s", cm.Username(ctx.Request()))
}
