package controllers

import (
	"../templates"
	"github.com/kataras/iris"
)
// Hello使用已编译的../templates/hello.qtpl.go文件渲染我们的../templates/hello.qtpl文件。
func Hello(ctx iris.Context) {
	// vars := make(map[string]interface{})
	// vars["message"] = "Hello World!"
	// vars["name"] = ctx.Params().Get("name")
	// [...]
	// &templates.Hello{ Vars: vars }
	// [...]
	//但是，作为替代方案，我们建议您使用`ctx.ViewData(key，value)`
	//为了能够从中间件（其他处理程序）修改`templates.Hello #Vars`
	ctx.ViewData("message", "Hello World!")
	ctx.ViewData("name", ctx.Params().Get("name"))
	// set view data to the `Vars` template's field
	tmpl := &templates.Hello{
		Vars: ctx.GetViewData(),
	}
	//渲染模板
	ExecuteTemplate(ctx, tmpl)
}