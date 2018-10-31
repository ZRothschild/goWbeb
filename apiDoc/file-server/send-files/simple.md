# Send Files

### 示例代码

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