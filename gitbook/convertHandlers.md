# 原始请求转换成`context.Handler`
## 1.面向中间件 `negroni-like`
### 目录结构
> 主目录`negroniLike`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"net/http"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	//FromStd将原生http.Handler,http.HandlerFunc和func(w，r，next)转换为context.Handler
	irisMiddleware := iris.FromStd(negronilikeTestMiddleware)
	app.Use(irisMiddleware)
	// 请求 GET: http://localhost:8080/
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1> Home </h1>")
		//这会打印错误
		//这个路由的处理程序永远不会被执行，因为中间件的标准没有通过
	})
	// 请求 GET: http://localhost:8080/ok
	app.Get("/ok", func(ctx iris.Context) {
		ctx.Writef("Hello world!")
		// 这将打印"OK. Hello world!"
	})
	// http://localhost:8080
	// http://localhost:8080/ok
	app.Run(iris.Addr(":8080"))
}

func negronilikeTestMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == "/ok" && r.Method == "GET" {
		w.Write([]byte("OK. "))
		next(w, r) // 中间件
		return
	}
	//否则打印错误
	w.WriteHeader(iris.StatusBadRequest)
	w.Write([]byte("Bad request"))
}
//如果要使用自定义上下文转换自定义处理程序，请查看routing/custom-context
//一个context.Handler
```
> `negroni`是`go`实现中间件的一种使用很广泛的模式

## 2.面向`net/Http`
### 目录结构
> 主目录`nethttp`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"net/http"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	//FromStd将原生http.Handler,http.HandlerFunc和func(w，r，next)转换为context.Handler
	irisMiddleware := iris.FromStd(nativeTestMiddleware)
	app.Use(irisMiddleware)
	// 请求 GET: http://localhost:8080/
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("Home")
	})
	// 请求 GET: http://localhost:8080/ok
	app.Get("/ok", func(ctx iris.Context) {
		ctx.HTML("<b>Hello world!</b>")
	})
	// http://localhost:8080
	// http://localhost:8080/ok
	app.Run(iris.Addr(":8080"))
}

func nativeTestMiddleware(w http.ResponseWriter, r *http.Request) {
	println("Request path: " + r.URL.Path)
}
//如果要使用自定义上下文转换自定义处理程序，请查看routing/custom-context
//一个context.Handler。
```
> 直接把原生的`net/http`满足请求接口的函数转化成中间件

## 3.错误中间件(raven)客户端 `https://sentry.io/welcome/`
### 1.修饰路由类型
#### 目录结构
> 主目录`realUsecaseRaven`

```html
    —— main.go
```
#### 代码示例
> `main.go`

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"github.com/kataras/iris"
	"github.com/getsentry/raven-go"
)

// https://docs.sentry.io/clients/go/integrations/http/
func init() {
	raven.SetDSN("https://<key>:<secret>@sentry.io/<project>")
}

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hi")
	})
	// WrapRouter的示例已在此处:
	// https://github.com/kataras/iris/blob/master/_examples/routing/custom-wrapper/main.go#L53
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, irisRouter http.HandlerFunc) {
		//完全相同的源代码:
		// https://github.com/getsentry/raven-go/blob/379f8d0a68ca237cf8893a1cdfd4f574125e2c51/http.go#L70
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				rvalStr := fmt.Sprint(rval)
				packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)), raven.NewHttp(r))
				raven.Capture(packet, nil)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		irisRouter(w, r)
	})
	app.Run(iris.Addr(":8080"))
}
```
### 2.直接中间件类型 
#### 目录结构
> 主目录`writingMiddleware`

```html
    —— main.go
```
#### 代码示例
> `main.go`

```go
package main

import (
	"errors"
	"fmt"
	"runtime/debug"
	"github.com/kataras/iris"
	"github.com/getsentry/raven-go"
)
//raven实现了Sentry错误日志记录服务的客户端。
//在此示例中，您将看到如何转换任何net/http中间件
//具有`(HandlerFunc)HandlerFunc`的形式。
//如果`raven.RecoveryHandler`的形式是
//`(http.HandlerFunc)`或`(http.HandlerFunc,下一个http.HandlerFunc)`
//你可以使用`irisMiddleware：= iris.FromStd(nativeHandler)`
//但它没有，但是你已经知道Iris可以直接使用net/http
//因为`ctx.ResponseWriter()`和`ctx.Request()`是原来的
// http.ResponseWriter和* http.Request。
//(这个是一个很大的优势，因此你可以永远使用Iris :))。
//本机中间件的源代码根本不会改变。
// https://github.com/getsentry/raven-go/blob/379f8d0a68ca237cf8893a1cdfd4f574125e2c51/http.go#L70
//唯一的补充是第18行和第39行(而不是handler(w，r))
//你有一个新的Iris中间件准备使用！
func irisRavenMiddleware(ctx iris.Context) {
	w, r := ctx.ResponseWriter(), ctx.Request()
	defer func() {
		if rval := recover(); rval != nil {
			debug.PrintStack()
			rvalStr := fmt.Sprint(rval)
			packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)), raven.NewHttp(r))
			raven.Capture(packet, nil)
			w.WriteHeader(iris.StatusInternalServerError)
		}
	}()
	ctx.Next()
}

// https://docs.sentry.io/clients/go/integrations/http/
func init() {
	raven.SetDSN("https://<key>:<secret>@sentry.io/<project>")
}

func main() {
	app := iris.New()
	app.Use(irisRavenMiddleware)
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hi")
	})
	app.Run(iris.Addr(":8080"))
}
```
## 提示
1. 主要介绍如何将原始请求转换成`context.Handler`
2. `Negroni`是面向`web`中间件的一种惯用方法.它是微小的,非侵入性的,并且鼓励使用`net/http`处理器
3. 错误中间件(raven)客户端为我们提供了一个非常清晰的报表,展示我们的响应并提高质量.这是完整的堆栈跟踪.这是一个非常强大的东西