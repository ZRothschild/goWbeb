//main包显示了如何将中间件添加到mvc应用程序中
//使用它的`Router`，它是主iris app的子路由器（iris.Party）。
package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/cache"
	"github.com/kataras/iris/mvc"
)

var cacheHandler = cache.Handler(10 * time.Second)

func main() {
	app := iris.New()
	mvc.Configure(app, configure)
	// http://localhost:8080
	// http://localhost:8080/other
	// //每10秒刷新一次，你会看到不同的时间输出。
	app.Run(iris.Addr(":8080"))
}

func configure(m *mvc.Application) {
	m.Router.Use(cacheHandler)
	m.Handle(&exampleController{
		timeFormat: "Mon, Jan 02 2006 15:04:05",
	})
}

type exampleController struct {
	timeFormat string
}

func (c *exampleController) Get() string {
	now := time.Now().Format(c.timeFormat)
	return "last time executed without cache: " + now
}

func (c *exampleController) GetOther() string {
	now := time.Now().Format(c.timeFormat)
	return "/other: " + now
}