# `iris`获取引用者(`extract referer`)
## 代码示例
### 主目录`extractReferer`
> 文件名称`main.go`
```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx context.Context) /*或iris.Context，Go 1.9+也是如此*/ {
		// GetReferrer提取并返回指定的Referer标头中的信息
		//在https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy中或通过URL查询参数是referer。
		r := ctx.GetReferrer()
		switch r.Type {
		case context.ReferrerSearch:
			ctx.Writef("Search %s: %s\n", r.Label, r.Query)
			ctx.Writef("Google: %s\n", r.GoogleType)
		case context.ReferrerSocial:
			ctx.Writef("Social %s\n", r.Label)
		case context.ReferrerIndirect:
			ctx.Writef("Indirect: %s\n", r.URL)
		}
	})
	//URL查询参数是referer
	// http://localhost:8080?referer=https://twitter.com/Xinterio/status/1023566830974251008
	// http://localhost:8080?referer=https://www.google.com/search?q=Top+6+golang+web+frameworks&oq=Top+6+golang+web+frameworks
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```
## 目录结构
> 主目录`extractReferer`
```html
    —— main.go
```
## 提示
1. 查询网站页面跳转referrer信息