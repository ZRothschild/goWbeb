package main

import (
	"math/rand"
	"time"
	"github.com/kataras/iris"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	app := iris.New()
	m := prometheusMiddleware.New("serviceName", 300, 1200, 5000)
	app.Use(m.ServeHTTP)
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		//错误代码处理程序不与其他路由共享相同的中间件，所以单独执行错误
		m.ServeHTTP(ctx)
		ctx.Writef("Not Found")
	})
	app.Get("/", func(ctx iris.Context) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		ctx.Writef("Slept for %d milliseconds", sleep)
	})
	app.Get("/metrics", iris.FromStd(prometheus.Handler()))
	// http://localhost:8080/
	// http://localhost:8080/anotfound
	// http://localhost:8080/metrics
	app.Run(iris.Addr(":8080"))
}
