# 以`gzip`格式写入数据(`write gzip`)
## 目录结构
> 主目录`writeGzip`
```html
    —— main.go
```
## 代码示例
> `main.go`
```go
package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()
	//返回一个gz压缩包
	app.Get("/", func(ctx iris.Context) {
		ctx.WriteGzip([]byte("Hello World!"))
		ctx.Header("X-Custom",
			"Headers can be set here after WriteGzip as well, because the data are kept before sent to the client when using the context's GzipResponseWriter and ResponseRecorder.")
	})
	//以gzip格式返回
	app.Get("/2", func(ctx iris.Context) {
		//与`WriteGzip`相同。
		//然而，GzipResponseWriter为您提供了更多选项，例如
		//重置数据，禁用等等，查看其方法。
		ctx.GzipResponseWriter().WriteString("Hello World!")
	})
	app.Run(iris.Addr(":8080"))
}
```