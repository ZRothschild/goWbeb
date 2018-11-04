package controllers

import (
	"../templates"
	"github.com/kataras/iris"
)
// ExecuteTemplate将“tmpl”部分模板渲染给`context＃ResponseWriter`。
func ExecuteTemplate(ctx iris.Context, tmpl templates.Partial) {
	ctx.Gzip(true)
	ctx.ContentType("text/html")
	templates.WriteTemplate(ctx.ResponseWriter(), tmpl)
}