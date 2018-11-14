# `cookies`使用
## 1.基础`cookies`操作
### 目录结构
> 主目录`basic`

```html
    —— main.go
    —— main_test.go
```
### 代码示例
> `main.go`

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
> `main_test.go`

```go
package main

import (
	"fmt"
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestCookiesBasic(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://example.com"))

	cookieName, cookieValue := "my_cookie_name", "my_cookie_value"

	// Test Set A Cookie.
	t1 := e.GET(fmt.Sprintf("/cookies/%s/%s", cookieName, cookieValue)).Expect().Status(httptest.StatusOK)
	t1.Cookie(cookieName).Value().Equal(cookieValue) // validate cookie's existence, it should be there now.
	t1.Body().Contains(cookieValue)

	// Test Retrieve A Cookie.
	t2 := e.GET(fmt.Sprintf("/cookies/%s", cookieName)).Expect().Status(httptest.StatusOK)
	t2.Body().Equal(cookieValue)

	// Test Remove A Cookie.
	t3 := e.DELETE(fmt.Sprintf("/cookies/%s", cookieName)).Expect().Status(httptest.StatusOK)
	t3.Body().Contains(cookieName)

	t4 := e.GET(fmt.Sprintf("/cookies/%s", cookieName)).Expect().Status(httptest.StatusOK)
	t4.Cookies().Empty()
	t4.Body().Empty()
}
```
## 2.`cookies`加密
### 目录结构
> 主目录`securecookie`

```html
    —— main.go
    —— main_test.go
```
### 代码示例
> `main.go`

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
		//加密值
		ctx.SetCookieKV(name, value, iris.CookieEncode(sc.Encode)) // <--设置一个Cookie
		ctx.Writef("cookie added: %s = %s", name, value)
	})
	app.Get("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		//解密值
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
> `main_test.go`

```go
package main

import (
	"fmt"
	"testing"
	"github.com/kataras/iris/httptest"
)

func TestCookiesBasic(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://example.com"))

	cookieName, cookieValue := "my_cookie_name", "my_cookie_value"

	// Test Set A Cookie.
	t1 := e.GET(fmt.Sprintf("/cookies/%s/%s", cookieName, cookieValue)).Expect().Status(httptest.StatusOK)
	// note that this will not work because it doesn't always returns the same value:
	// cookieValueEncoded, _ := sc.Encode(cookieName, cookieValue)
	t1.Cookie(cookieName).Value().NotEqual(cookieValue) // validate cookie's existence and value is not on its raw form.
	t1.Body().Contains(cookieValue)

	// Test Retrieve A Cookie.
	t2 := e.GET(fmt.Sprintf("/cookies/%s", cookieName)).Expect().Status(httptest.StatusOK)
	t2.Body().Equal(cookieValue)

	// Test Remove A Cookie.
	t3 := e.DELETE(fmt.Sprintf("/cookies/%s", cookieName)).Expect().Status(httptest.StatusOK)
	t3.Body().Contains(cookieName)

	t4 := e.GET(fmt.Sprintf("/cookies/%s", cookieName)).Expect().Status(httptest.StatusOK)
	t4.Cookies().Empty()
	t4.Body().Empty()
}
```