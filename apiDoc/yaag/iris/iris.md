# `YAAG `生成`iris web`框架项目`API`文档
## API 生成步骤
- 下载`YAAG`中间件
> go get github.com/betacraft/yaag/...
- 导入依赖包
> import github.com/betacraft/yaag/yaag
> Import github.com/betacraft/yaag/irisyaag
- 初始化`yaag`
> yaag.Init(&yaag.Config(On: true, DocTile: "Iris", DocPath: "apidoc.html"))
- 注册`yaag`中间件
> app.Use(irisyaag.New())
> irisyaag记录响应主体并向apidoc提供所有必要的信息
## 目录结构
> 主目录`iris`

```html
    —— apidoc.html (执行命令后生成)
    —— apidoc.html.json (执行命令后生成)
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris"

	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
)
/*
	下载包 go get github.com/betacraft/yaag/...
*/
type myXML struct {
	Result string `xml:"result"`
}

func main() {
	app := iris.New()
	//初始化中间件
	yaag.Init(&yaag.Config{
		On:       true,                 //是否开启自动生成API文档功能
		DocTitle: "Iris",
		DocPath:  "apidoc.html",        //生成API文档名称存放路径
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})
	//注册中间件
	app.Use(irisyaag.New())
	app.Get("/json", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"result": "Hello World!"})
	})
	app.Get("/plain", func(ctx iris.Context) {
		ctx.Text("Hello World!")
	})
	app.Get("/xml", func(ctx iris.Context) {
		ctx.XML(myXML{Result: "Hello World!"})
	})
	app.Get("/complex", func(ctx iris.Context) {
		value := ctx.URLParam("key")
		ctx.JSON(iris.Map{"value": value})
	})
	//运行HTTP服务器。
	//每个传入的请求都会重新生成和更新“apidoc.html”文件。

	//编写调用这些处理程序的测试，保存生成的apidoc.html，apidoc.html.json。
	//在制作时关闭yaag中间件。
	app.Run(iris.Addr(":8080"))
}
```
### 提示
1. 运行上面的例子，并请求其中任意一个接口，会生成`apidoc.html`，`apidoc.html.json`两个文件
2. 如果不关闭`yaag`中间件，`apidoc.html`，`apidoc.html.json`文件会随着每一次请求而重新生成
3. 如果没有翻墙，可能看不到生成效果，因为`apidoc.html`文件的引入了很多谷歌的`cnd`,替换成功国内支持即可