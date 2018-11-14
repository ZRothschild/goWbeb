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