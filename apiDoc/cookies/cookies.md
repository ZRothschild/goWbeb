# cookies使用

### 基础cookies操作

```go
package main

import "github.com/kataras/iris"

func newApp() *iris.Application {
	app := iris.New()
	app.Get("/cookies/{name}/{value}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		value := ctx.Params().Get("value")
		ctx.SetCookieKV(name, value) // <-- 设置一个Cookie
		// 另外也可以用: ctx.SetCookie(&http.Cookie{...})
		// 如果要设置自定义存放路径：
		// ctx.SetCookieKV(name, value, iris.CookiePath("/custom/path/cookie/will/be/stored"))

		ctx.Request().Cookie(name)
		//如果您希望仅对当前请求路径可见：
		//（请注意，如果服务器发送空cookie的路径，所有浏览器都兼容，将会使用客户端自定义路径）
		// ctx.SetCookieKV(name, value, iris.CookieCleanPath /* or iris.CookiePath("") */)
		// 学习更多:
		//                              iris.CookieExpires(time.Duration)
		//                              iris.CookieHTTPOnly(false)
		ctx.Writef("cookie added: %s = %s", name, value)
	})
	app.Get("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		value := ctx.GetCookie(name) // <-- 检索，获取Cookie
		//判断命名cookie不存在，再获取值
		// cookie, err := ctx.Request().Cookie(name)
		// if err != nil {
		//  handle error.
		// }
		ctx.WriteString(value)
	})
	app.Delete("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.RemoveCookie(name) // <-- 删除Cookie
		//如果要设置自定义路径：
		// ctx.SetCookieKV(name, value, iris.CookiePath("/custom/path/cookie/will/be/stored"))
		ctx.Writef("cookie %s removed", name)
	})
	return app
}

func main() {
	app := newApp()
	// GET:    http://localhost:8080/cookies/my_name/my_value
	// GET:    http://localhost:8080/cookies/my_name
	// DELETE: http://localhost:8080/cookies/my_name
	app.Run(iris.Addr(":8080"))
}
```

### 如何加密cookies

```go
package main
//开发人员可以使用任何库添加自定义cookie编码器/解码器。
//在这个例子中，我们使用gorilla的securecookie包：
// $ go get github.com/gorilla/securecookie
// $ go run main.go

import (
	"github.com/kataras/iris"
	"github.com/gorilla/securecookie"
)

var (
	// AES仅支持16,24或32字节的密钥大小。
	//您需要准确提供该密钥字节大小，或者从您键入的内容中获取密钥。
	hashKey  = []byte("the-big-and-secret-fash-key-here")
	blockKey = []byte("lot-secret-of-characters-big-too")
	sc       = securecookie.New(hashKey, blockKey)
)

func newApp() *iris.Application {
	app := iris.New()
	app.Get("/cookies/{name}/{value}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		value := ctx.Params().Get("value")
		ctx.SetCookieKV(name, value, iris.CookieEncode(sc.Encode)) // <--设置一个Cookie
		ctx.Writef("cookie added: %s = %s", name, value)
	})
	app.Get("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		value := ctx.GetCookie(name, iris.CookieDecode(sc.Decode)) // <--检索Cookie
		ctx.WriteString(value)
	})
	app.Delete("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.RemoveCookie(name) // <-- 删除Cookie
		ctx.Writef("cookie %s removed", name)
	})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}
```

[Go Web Iris中文网](https://www.studyiris.com/)