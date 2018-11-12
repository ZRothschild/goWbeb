# `Single Page Application`基本运用
## 目录结构
> 主目录`basic`
```html
—— main.go
—— main_test.go
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

import (
	"github.com/kataras/iris"
)
//与嵌入式单页面应用程序相同但没有go-bindata，文件是"原始"存储在
//当前系统目录
var page = struct {
	Title string
}{"Welcome"}

func newApp() *iris.Application {
	app := iris.New()
	app.RegisterView(iris.HTML("./public", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Page", page)
		ctx.View("index.html")
	})
	//或者只是按原样提供index.html：
	// app.Get("/{f:path}", func(ctx iris.Context) {
	// 	ctx.ServeFile("index.html", false)
	// })
	assetHandler := app.StaticHandler("./public", false, false)
	//作为SPA的替代方案，您可以查看/routing /dynamic-path/root-wildcard
	app.SPA(assetHandler)
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	// http://localhost:8080/index.html
	// http://localhost:8080/app.js
	// http://localhost:8080/css/main.css
	app.Run(iris.Addr(":8080"))
}
```
> `main_test.go`
```go
package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"github.com/kataras/iris/httptest"
)

type resource string

func (r resource) String() string {
	return string(r)
}

func (r resource) strip(strip string) string {
	s := r.String()
	return strings.TrimPrefix(s, strip)
}

func (r resource) loadFromBase(dir string) string {
	filename := r.String()

	if filename == "/" {
		filename = "/index.html"
	}

	fullpath := filepath.Join(dir, filename)

	b, err := ioutil.ReadFile(fullpath)
	if err != nil {
		panic(fullpath + " failed with error: " + err.Error())
	}

	result := string(b)

	return result
}

var urls = []resource{
	"/",
	"/index.html",
	"/app.js",
	"/css/main.css",
}

func TestSPA(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.Debug(false))

	for _, u := range urls {
		url := u.String()
		contents := u.loadFromBase("./public")
		contents = strings.Replace(contents, "{{ .Page.Title }}", page.Title, 1)

		e.GET(url).Expect().
			Status(httptest.StatusOK).
			Body().Equal(contents)
	}
}
```
> `public/app.js`
```js
window.alert("app.js loaded from \"/");
```
> `public/index.html`
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
> `public/css/main.css`
```css
body {
    background-color: black;
}
```