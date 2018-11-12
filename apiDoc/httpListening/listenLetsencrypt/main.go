//包main提供与letsencrypt.org的单行集成
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello from SECURE SERVER!")
	})
	app.Get("/test2", func(ctx iris.Context) {
		ctx.Writef("Welcome to secure server from /test2!")
	})
	app.Get("/redirect", func(ctx iris.Context) {
		ctx.Redirect("/test2")
	})
	//注意：这不适用于这样的域名，
	//使用真正的白名单域（或由空格分割的域）
	//而不是非公开的电子邮件。
	app.Run(iris.AutoTLS(":443", "example.com", "mail@example.com"))
}
