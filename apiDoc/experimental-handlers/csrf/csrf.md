# `csrf`防御

##  `csrf`介绍

`CSRF`(`Cross-site request forgery`)跨站请求伪造，也被称为`One Click Attack`或者`Session Riding`，通常缩写为`CSRF`或者`XSRF`，
是一种对网站的恶意利用。尽管听起来像跨站脚本（`XSS`），但它与`XSS`非常不同，`XSS`利用站点内的信任用户，而`CSRF`则通过伪装来自受信
任用户的请求来利用受信任的网站。与`XSS`攻击相比，`CSRF`攻击往往不大流行（因此对其进行防范的资源也相当稀少）和难以防范，所以被认为
比`XSS`更具危险性

## 攻击案例

攻击通过在授权用户访问的页面中包含链接或者脚本的方式工作。例如：一个网站用户`Bob`可能正在浏览聊天论坛，而同时另一个用户`Alice`也在
此论坛中，并且后者刚刚发布了一个具有`Bob`银行链接的图片消息。设想一下，`Alice`编写了一个在`Bob`的银行站点上进行取款的`form`提交的
链接，并将此链接作为图片src。如果`Bob`的银行在`cookie`中保存他的授权信息，并且此`cookie`没有过期，那么当`Bob`的浏览器尝试装载图
片时将提交这个取款`form`和他的`cookie`，这样在没经`Bob`同意的情况下便授权了这次事务。`CSRF`是一种依赖`web`浏览器的、被混淆过的代
理人攻击（`deputy attack`）。在上面银行示例中的代理人是`Bob`的`web`浏览器，它被混淆后误将`Bob`的授权直接交给了`Alice`使用

## `CSRF`防御

- 通过`referer`,`token` 或者验证码来检测用户提交
- 尽量不要在页面的链接中暴露用户隐私信息
- 对于用户修改删除等操作最好都使用`post`操作
- 避免全站通用的`cookie`，严格设置`cookie`的域

## 示例代码

> go代码 `main.go`
```go
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
```
> `html`代码 ./views/user/signup.html

```html
<form method="POST" action="/user/signup">
    {{ .csrfField }}
<button type="submit">Proceed</button>
</form>
```