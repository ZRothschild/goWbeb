# 单页面应用(`embedded Single Page Application Other router`)
## 目录结构
> 主目录`embeddedSPAWithOtherRoutes`

```html
—— main.go
—— public
    —— css
        —— main.css
    —— index.html
    —— app.js
```
## 示例代码
> `main.go`

```go
package main

import "github.com/kataras/iris"

// $ go get -u github.com/shuLhan/go-bindata/...
// $ go-bindata ./public/...
// $ go build
// $ ./embedded-single-page-application-with-other-routes

func newApp() *iris.Application {
	app := iris.New()
	app.OnErrorCode(404, func(ctx iris.Context) {
		ctx.Writef("404 not found here")
	})
	app.StaticEmbedded("/", "./public", Asset, AssetNames)
	//注意：
	//如果您想要一个动态索引页面，请查看file-server/embedded-single-page-application
	//正在注册基于bindata的视图引擎和根路由。
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	app.Get("/.well-known", func(ctx iris.Context) {
		ctx.WriteString("well-known")
	})
	app.Get(".well-known/ready", func(ctx iris.Context) {
		ctx.WriteString("ready")
	})
	app.Get(".well-known/live", func(ctx iris.Context) {
		ctx.WriteString("live")
	})
	app.Get(".well-known/metrics", func(ctx iris.Context) {
		ctx.Writef("metrics")
	})
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080/index.html
	// http://localhost:8080/app.js
	// http://localhost:8080/css/main.css
	// http://localhost:8080/ping
	// http://localhost:8080/.well-known
	// http://localhost:8080/.well-known/ready
	// http://localhost:8080/.well-known/live
	// http://localhost:8080/.well-known/metrics
	//记住：我们可以使用根通配符`app.Get("/{param:path}")`并手动提供文件。
	app.Run(iris.Addr(":8080"))
}
```
> `/public/app.js`

```js
window.alert("app.js loaded from \"/");
```
> `/public/index.html`

```html
<html>
<head>
    <title>{{ .Page.Title }}</title>
</head>
<body>
    <h1> Hello from index.html </h1>
    <script src="/app.js">  </script>
</body>
</html>
```
> `/public/css/main.css`

```css
body {
    background-color: black;
}
```