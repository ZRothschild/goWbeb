# `embedded Single Page Application`(单页面应用)
## 目录结构
> 主目录`embeddedSinglePageApplication`
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

// $ go get -u github.com/shuLhan/go-bindata/...
// $ go-bindata ./public/...
// $ go build
// $ ./embedded-single-page-application

var page = struct {
	Title string
}{"Welcome"}

func newApp() *iris.Application {
	app := iris.New()
	app.RegisterView(iris.HTML("./public", ".html").Binary(Asset, AssetNames))

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Page", page)
		ctx.View("index.html")
	})

	assetHandler := iris.StaticEmbeddedHandler("./public", Asset, AssetNames, false) //如果使用`go-bindata`工具，请保持false。
	//作为SPA的替代方案，您可以查看/routing/dynamic-path/root-wildcard
	//也是例子
	// 要么
	// app.StaticEmbedded如果您不想在index.html上重定向并且简单地为您的SPA应用程序提供服务（推荐）。

	// public / index.html是一个动态视图，它由root手工绘制，
	//我们不希望作为原始数据显示，所以我们会
	//'app.SPA`的返回值来修改`IndexNames`;
	app.SPA(assetHandler).AddIndexName("index.html")
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
//请注意，将执行app.Use/UseGlobal/Done
//仅限于注册的路由，如我们的app.Get("/"，..）。
//文件服务器将不受限制，但你仍然可以通过修饰它的assetHandler来添加中间件。

//使用此方法，与静态Web("/"，'./ public')不同，后者不再按设计工作，
//所有自定义http错误和所有路由都可以正常使用已注册的文件服务器
//到服务器的根路径
```
> `main_test.go`
```go
package main

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"github.com/kataras/iris/httptest"
)

type resource string

func (r resource) contentType() string {
	switch filepath.Ext(r.String()) {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	default:
		return "text/html"
	}
}

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
	if runtime.GOOS != "windows" {
		// result = strings.Replace(result, "\n", "\r\n", -1)
	}
	return result
}

var urls = []resource{
	"/",
	"/index.html",
	"/app.js",
	"/css/main.css",
}

func TestSPAEmbedded(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	for _, u := range urls {
		url := u.String()
		contents := u.loadFromBase("./public")
		contents = strings.Replace(contents, "{{ .Page.Title }}", page.Title, 1)

		e.GET(url).Expect().
			Status(httptest.StatusOK).
			ContentType(u.contentType(), app.ConfigurationReadOnly().GetCharset()).
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