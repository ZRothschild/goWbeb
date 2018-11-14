# `CORS`跨域资源共享
## `CORS`简介
`CORS`是一个`W3C`标准，全称是"跨域资源共享"（`Cross-origin resource sharing`）。它允许浏览器向跨源(协议 + 域名 + 端口)服务器，
发出`XMLHttpRequest`请求，从而克服了`AJAX`只能同源使用的限制。

`CORS`需要浏览器和服务器同时支持。它的通信过程，都是浏览器自动完成，不需要用户参与。对于开发者来说，`CORS`通信与同源的`AJAX`通信没
有差别，代码完全一样。浏览器一旦发现`AJAX`请求跨源，就会自动添加一些附加的头信息，有时还会多出一次附加的请求，但用户不会有感觉。因此，
实现`CORS`通信的关键是服务器。只要服务器实现了`CORS`接口，就可以跨源通信
## 目录结构
> 主目录`simple`
```html
    —— main.go
```
## 代码示例 
> `main.go`

```go
package main

// go get -u github.com/iris-contrib/middleware/...

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
)

//跨与请求 下面代码表示 http://foo.com 站点下的 ajax 可以跨域请求 localhost:8080 接口
//当http://foo.com 为 * 表示所有域名都可以请求

//AllowedOrigins 该字段是必须的。
// 它的值要么是请求时Origin字段的值，要么是一个*，表示接受任意域名的请求。

//AllowCredentials  该字段可选。它的值是一个布尔值，表示是否允许发送Cookie。
//默认情况下，Cookie不包括在CORS请求之中。
//设为true，即表示服务器明确许可，Cookie可以包含在请求中，一起发给服务器。
//这个值也只能设为true，如果服务器不要浏览器发送Cookie，删除该字段即可。

func main() {
	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://foo.com"},   //允许通过的主机名称
		AllowCredentials: true,
	})
	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions) // <- 对于预检很重要。
	{
		v1.Get("/home", func(ctx iris.Context) {
			ctx.WriteString("Hello from /home")
		})
		v1.Get("/about", func(ctx iris.Context) {
			ctx.WriteString("Hello from /about")
		})
		v1.Post("/send", func(ctx iris.Context) {
			ctx.WriteString("sent")
		})
		v1.Put("/send", func(ctx iris.Context) {
			ctx.WriteString("updated")
		})
		v1.Delete("/send", func(ctx iris.Context) {
			ctx.WriteString("deleted")
		})
	}
	app.Run(iris.Addr("localhost:8080"))
}
````