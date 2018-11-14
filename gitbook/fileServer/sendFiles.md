# `Send Files`
## 目录结构
> 主目录`sendFiles`

```html
—— files
    —— first.zip
—— main.go
```
## 示例代码
> `main.go`

```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		file := "./files/first.zip"
		ctx.SendFile(file, "c.zip")
	})
	app.Run(iris.Addr(":8080"))
}
```