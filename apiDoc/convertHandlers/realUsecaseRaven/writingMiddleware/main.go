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
