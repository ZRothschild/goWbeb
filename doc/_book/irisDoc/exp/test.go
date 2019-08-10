package main

import (
	"github.com/kataras/iris"
	"fmt"
)

func main() {
	app := iris.New()
	app.Configure(iris.WithConfiguration(iris.Configuration{}))
	app.Get("/name", Dof)
	// [...]

	app.Configure(iris.WithFireMethodNotAllowed)
	// Good when you want to modify the whole configuration.
	app.Run(iris.Addr(":80"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	}))

}
func Dof(ctx iris.Context) {
	fmt.Println(ctx.Path())
	fmt.Println(ctx.HandlerName())
	ctx.HTML("<b>Hello!</b>")

}