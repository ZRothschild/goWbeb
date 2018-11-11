# `quicktemplate`

## quicktemplate介绍

`Go`的快速，强大且易于使用的模板引擎。 针对热路径中的速度，零内存分配进行了优化。比`html/template`快20倍，
灵感来自[Mako templates](http://www.makotemplates.org/`)

1. 非常快。 模板转换为Go代码然后编译
2. `Quicktemplate`语法非常接近`Go` - 在开始使用`quicktemplate`之前无需学习另一种模板语言
3. 在模板编译期间几乎所有错误都被捕获，因此生产受模板相关错误的影响较小
4. 使用方便。有关详细信息，请参阅快速入门和示例
5. 强大。任意`Go`代码可以嵌入到模板中并与模板混合。小心这个功能 - 不要从模板中查询数据库`and/or`外部资源，除非你错过`Go`中的`PHP`方
式`:)`这种功能主要用于任意数据转换
6. 易于使用的模板继承由`Go`接口提供支持。 请参阅此示例以获取详细信
7. 模板被编译为单个二进制文件，因此无需将模板文件复制到服务器
> 模板无法在服务器上动态更新，因为它们被编译为单个二进制文件。如果您需要快速模板引擎来简单动态更新模板，请查看
[fasttemplate](https://github.com/valyala/fasttemplate)

## quicktemplate使用

首先，安装[quicktemplate](https://github.com/valyala/quicktemplate)包和
[quicktemplate compiler](https://github.com/valyala/quicktemplate/tree/master/qtc)

```sh
go get -u github.com/valyala/quicktemplate
go get -u github.com/valyala/quicktemplate/qtc
```
可以在https://github.com/valyala/quicktemplate找到完整的文档。

将模板文件保存到扩展名`*.qtpl`下的`templates`文件夹中，打开终端并在此文件夹中运行`qtc`。

如果一切顺利，`*.qtpl.go`文件必须出现在`templates`文件夹中。 这些文件包含所有`* .qtpl`文件的Go代码。

> 请记住，每次更改`/templates/*.qtpl`文件时，都必须运行`qtc`命令并重新构建应用程序。
## `iris`代码示例
### `controllers`文件夹
> 文件名 `execute_template.go`
```go
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
```
> 文件名 `hello.go`
```go
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
```
> 文件名 `index.go`
```go
package controllers

import (
	"../templates"
	"github.com/kataras/iris"
)
//索引使用已编译的../templates/index.qtpl.go文件渲染我们的../templates/index.qtpl文件。
func Index(ctx iris.Context) {
	tmpl := &templates.Index{}
	//渲染模板
	ExecuteTemplate(ctx, tmpl)
}
```
### `templates`文件夹
> 文件名 `base.qtpl`
```sh
这是我们模板的基础实现
{% interface
Partial {
	Body()
}
%}
模板编写实现Partial接口的模板
{% func Template(p Partial) %}
<html>
	<head>
		<title>Quicktemplate integration with Iris</title>
	</head>
	<body>
		<div>
			Header contents here...
		</div>
		<div style="margin:10px;">
			{%= p.Body() %}
		</div>
	</body>
	<footer>
		Footer contents here...
	</footer>
</html>
{% endfunc %}
基本模板实现。 如果需要，其他页面可以继承
仅覆盖某些部分方法。
{% code type Base struct {} %}
{% func (b *Base) Body() %}This is the base body{% endfunc %}
```
> 文件名 `hello.qtpl`
```sh
Hello模板，实现了Partial的方法
{% code
type Hello struct {
  Vars map[string]interface{}
}
%}
{% func (h *Hello) Body() %}
	<h1>{%v h.Vars["message"] %}</h1>
	<div>
		Hello <b>{%v h.Vars["name"] %}!</b>
	</div>
{% endfunc %}
```
> 文件名 `index.qtpl`
```sh
Hello模板，实现了Partial的方法。
{% code
type Index struct {}
%}
{% func (i *Index) Body() %}
	<h1>Index Page</h1>
	<div>
		This is our index page's body.
	</div>
{% endfunc %}
```
### 主目录文件夹`quicktemplate`
> 文件名称 `main.go`
```go
package main

import (
	"./controllers"
	"github.com/kataras/iris"
)

func newApp() *iris.Application {
	app := iris.New()
	app.Get("/", controllers.Index)
	app.Get("/{name}", controllers.Hello)
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	// http://localhost:8080/yourname
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```
> 文件名称 `main_test.go`
```go
package main

import (
	"fmt"
	"testing"
	"github.com/kataras/iris/httptest"
)

func TestResponseWriterQuicktemplate(t *testing.T) {
	baseRawBody := `
<html>
	<head>
		<title>Quicktemplate integration with Iris</title>
	</head>
	<body>
		<div>
			Header contents here...
		</div>
		<div style="margin:10px;">
	<h1>%s</h1>
	<div>
		%s
	</div>
		</div>
	</body>
	<footer>
		Footer contents here...
	</footer>
</html>
`
	expectedIndexRawBody := fmt.Sprintf(baseRawBody, "Index Page", "This is our index page's body.")
	name := "yourname"
	expectedHelloRawBody := fmt.Sprintf(baseRawBody, "Hello World!", "Hello <b>"+name+"!</b>")
	app := newApp()
	e := httptest.New(t, app)
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal(expectedIndexRawBody)
	e.GET("/" + name).Expect().Status(httptest.StatusOK).Body().Equal(expectedHelloRawBody)
}
```
## 目录结构
> 主目录`herotemplate`
```html
    —— controllers
        —— execute_template.go
        —— hello.go
        —— index.go
    —— template
        —— base.qtpl
        —— hello.qtpl
        —— index.qtpl
    —— main.go
    —— main_test.go
```
