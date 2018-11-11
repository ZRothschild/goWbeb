package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recaptcha"
)
//密钥应通过https://www.google.com/recaptcha获取
const (
	recaptchaPublic = "6Lf3WywUAAAAAKNfAm5DP2J5ahqedtZdHTYaKkJ6"
	recaptchaSecret = "6Lf3WywUAAAAAJpArb8nW_LCL_PuPuokmEABFfgw"
)

func main() {
	app := iris.New()
	r := recaptcha.New(recaptchaSecret)
	app.Get("/comment", showRecaptchaForm)
	//在主处理程序之前传递中间件或使用`recaptcha.SiteVerify`。
	app.Post("/comment", r, postComment)
	app.Run(iris.Addr(":8080"))
}

var htmlForm = `<form action="/comment" method="POST">
	    <script src="https://www.google.com/recaptcha/api.js"></script>
		<div class="g-recaptcha" data-sitekey="%s"></div>
    	<input type="submit" name="button" value="Verify">
</form>`

func showRecaptchaForm(ctx iris.Context) {
	contents := fmt.Sprintf(htmlForm, recaptchaPublic)
	ctx.HTML(contents)
}

func postComment(ctx iris.Context) {
	// [...]
	ctx.JSON(iris.Map{"success": true})
}