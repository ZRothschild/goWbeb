# 数据返回类型(`write rest`)
## 目录结构
> 主目录`writeRest`
```html
    —— main.go
```
## 代码示例
> `main.go`
```go
package main

import (
	"encoding/xml"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)
//用户绑定结构
type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

// ExampleXML只是一个要查看的测试结构，表示xml内容类型
type ExampleXML struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func main() {
	app := iris.New()
	// 读取
	app.Post("/decode", func(ctx iris.Context) {
		// 参考 /http_request/read-json/main.go
		var user User
		ctx.ReadJSON(&user)
		ctx.Writef("%s %s is %d years old and comes from %s!", user.Firstname, user.Lastname, user.Age, user.City)
	})
	// Write
	app.Get("/encode", func(ctx iris.Context) {
		peter := User{
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
			Age:       25,
		}
		//手动设置内容类型: ctx.ContentType("application/javascript")
		ctx.JSON(peter)
	})
	//其他内容类型
	app.Get("/binary", func(ctx iris.Context) {
		//当您想要强制下载原始字节内容时有用下载文件
		ctx.Binary([]byte("Some binary data here."))
	})
	app.Get("/text", func(ctx iris.Context) {
		ctx.Text("Plain text here")
	})
	app.Get("/json", func(ctx iris.Context) {
		ctx.JSON(map[string]string{"hello": "json"}) // or myjsonStruct{hello:"json}
	})
	app.Get("/jsonp", func(ctx iris.Context) {
		ctx.JSONP(map[string]string{"hello": "jsonp"}, context.JSONP{Callback: "callbackName"})
	})
	app.Get("/xml", func(ctx iris.Context) {
		ctx.XML(ExampleXML{One: "hello", Two: "xml"}) // or iris.Map{"One":"hello"...}
	})
	app.Get("/markdown", func(ctx iris.Context) {
		ctx.Markdown([]byte("# Hello Dynamic Markdown -- iris"))
	})
	// http://localhost:8080/decode
	// http://localhost:8080/encode
	// http://localhost:8080/binary
	// http://localhost:8080/text
	// http://localhost:8080/json
	// http://localhost:8080/jsonp
	// http://localhost:8080/xml
	// http://localhost:8080/markdown

	//`iris.WithOptimizations`是一个可选的配置器，
	//如果传递给`Run`那么它将确保应用程序，尽快响应客户端。

	// `iris.WithoutServerError` 是一个可选的配置器,
	//如果传递给`Run`那么它不会将传递的错误，实际的服务器错误打印。
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
```