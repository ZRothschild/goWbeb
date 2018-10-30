//此中间件提供跨站点请求伪造保护
//它安全地生成一个掩码（每个请求唯一）令牌
//可以嵌入HTTP响应中（例如表单字段或HTTP标头）。
//原始令牌存储在会话中，该会话无法访问
// 攻击者（如果您使用的是HTTPS后续请求是 期望包含此令牌，该令牌与会话令牌进行比较。
//匹配令牌失败包HTTP 403 Forbidden 错误响应。
package main
// $ go get -u github.com/iris-contrib/middleware/...
import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/csrf"
)
func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	//请注意，提供的身份验证密钥应为32个字节应用程序重新启动时保持不变
	protect := csrf.Protect([]byte("9AB0F421E53A477C084477AEA06096F5"),
		csrf.Secure(false)) //默认为true，但在没有https（devmode）的情况下传递`false`。
	users := app.Party("/user", protect)
	{
		users.Get("/signup", getSignupForm)
		//没有有效令牌的POST请求将返回HTTP 403 Forbidden。
		users.Post("/signup", postSignupForm)
	}
	// GET: http://localhost:8080/user/signup
	// POST: http://localhost:8080/user/signup
	app.Run(iris.Addr(":8080"))
}

func getSignupForm(ctx iris.Context) {
	// views/user/signup.html只需要一个{{.csrfField}}模板标记
	// csrf.TemplateField将CSRF令牌注入即可！
	ctx.ViewData(csrf.TemplateTag, csrf.TemplateField(ctx))
	ctx.View("user/signup.html")
	//我们也可以直接从csrf.Token（ctx）中获取检索令牌
	//在请求标头中设置它 - ctx.GetHeader("X-CSRF-Token"，token)
	//如果您要向客户端或前端JavaScript发送JSON，这将非常有用
	//框架
}

func postSignupForm(ctx iris.Context) {
	ctx.Writef("You're welcome mate!")
}
