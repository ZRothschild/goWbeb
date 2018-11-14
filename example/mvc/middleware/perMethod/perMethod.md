# `IRIS MVC`前置中间件使用
## 目录结构
> 主目录`perMethod`

```html
    —— main.go
```
## 代码示例
> `main.go`

```go
/*
如果要将其用作整个控制器的中间件
你可以使用它的路由，它只是一个子路由添加了中间件，就像你通常使用标准API一样：

我将向您展示将中间件添加到mvc应用程序的4种不同方法，
所有这4个做同样的事情，选择你喜欢的，
当我需要在某个地方注册中间件时，我更喜欢最后一个代码片段
否则我也会选择第一个：

// 1
mvc.Configure(app.Party("/user"), func(m *mvc.Application) {
     m.Router.Use(cache.Handler(10*time.Second))
})

// 2
// same:
userRouter := app.Party("/user")
userRouter.Use(cache.Handler(10*time.Second))
mvc.Configure(userRouter, ...)

// 3
userRouter := app.Party("/user", cache.Handler(10*time.Second))
mvc.Configure(userRouter, ...)

// 4
// same:
app.PartyFunc("/user", func(r iris.Party){
    r.Use(cache.Handler(10*time.Second))
    mvc.Configure(r, ...)
})

如果要将中间件用于单个路由，
对于已由引擎注册的单个控制器方法
而不是自定义`Handle`（你可以添加最后一个参数上的中间件）
并且它不依赖于`Next Handler`来完成它的工作
然后你只需在方法上调用它：

var myMiddleware := myMiddleware.New(...) //返回一个iris/context.Handler类型

type UserController struct{}
func (c *UserController) GetSomething(ctx iris.Context) {
	// ctx.Proceed检查myMiddleware是否调用`ctx.Next()`
     //在其中，如果是，则返回true，否则返回false。
    nextCalled := ctx.Proceed(myMiddleware)
    if !nextCalled {
        return
    }
   //其他工作，这在这里执行是允许的
}

最后，如果您想在特定方法上添加中间件
这取决于下一个和整个链，那么你必须这样做
像下面的例子一样使用`AfterActivation`：
*/
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
	//如果你在主函数中完成所有操作，则不必使用.Configure
	//mvc.Configure和mvc.New(...).Configure()只是拆分你的代码更好，
	//这里我们使用最简单的形式：
	m := mvc.New(app)
	m.Handle(&exampleController{})
	app.Run(iris.Addr(":8080"))
}

type exampleController struct{}

func (c *exampleController) AfterActivation(a mvc.AfterActivation) {
	//根据您想要的方法名称选择路由 修改
	index := a.GetRoute("Get")
	//只是将处理程序作为您想要使用的中间件预先添加。
	//或附加“done”处理程序。
	index.Handlers = append([]iris.Handler{cacheHandler}, index.Handlers...)
}

func (c *exampleController) Get() string {
	//每隔10秒刷新一次，你会看到不同的时间输出。
	now := time.Now().Format("Mon, Jan 02 2006 15:04:05")
	return "last time executed without cache: " + now
}
```