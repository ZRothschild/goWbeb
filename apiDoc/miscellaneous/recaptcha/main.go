package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recaptcha"
)

//密钥应通过https://www.google.com/recaptcha获取
const (
	recaptchaPublic = ""
	recaptchaSecret = ""
)

func showRecaptchaForm(ctx iris.Context, path string) {
	ctx.HTML(recaptcha.GetFormHTML(recaptchaPublic, path))
}

func main() {
	app := iris.New()
	// On both Get and Post on this example, so you can easly
	// use a single route to show a form and the main subject if recaptcha's validation result succeed.
	//在此示例的Get和Post上，您可以轻松
	//使用单个路径显示表单，并且重新验证结果的主要主题成功。
	app.HandleMany("GET POST", "/", func(ctx iris.Context) {
		if ctx.Method() == iris.MethodGet {
			showRecaptchaForm(ctx, "/")
			return
		}
		result := recaptcha.SiteFerify(ctx, recaptchaSecret)
		if !result.Success {
			/* 如果你想要或什么都不做，重定向到这里 */
			ctx.HTML("<b> failed please try again </b>")
			return
		}
		ctx.Writef("succeed.")
	})
	app.Run(iris.Addr(":8080"))
}