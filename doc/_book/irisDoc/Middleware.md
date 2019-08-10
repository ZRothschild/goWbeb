#####中间件

当我们在Iris中讨论中间件时，我们讨论的是在HTTP请求生命周期中在主处理程序代码之前和/或之后运行代码。例如，记录中间件可能会将传入的请求详细信息写入日志，然后在写入有关日志响应的详细信息之前调用处理程序代码。关于中间件的一个很酷的事情是这些单元非常灵活和可重用。

下面一个简单的示例

```go
    package main
    import "github.com/kataras/iris"
    func main() {
        app := iris.New()
        app.Get("/", before, mainHandler, after)
        app.Run(iris.Addr(":8080"))
    }
    func before(ctx iris.Context) {
        shareInformation := "this is a sharable information between handlers"

        requestPath := ctx.Path()
        println("Before the mainHandler: " + requestPath)

        ctx.Values().Set("info", shareInformation)
        ctx.Next() //继续执行下一个handler，在本例中是mainHandler。
    }
    func after(ctx iris.Context) {
        println("After the mainHandler")
    }
    func mainHandler(ctx iris.Context) {
        println("Inside mainHandler")
        // take the info from the "before" handler.
        info := ctx.Values().GetString("info")
        // write something to the client as a response.
        ctx.HTML("<h1>Response</h1>")
        ctx.HTML("<br/> Info: " + info)
        ctx.Next() // 继续下一个handler 这里是after
    }
```
试试看:
```go
    $ go run main.go # and navigate to the http://localhost:8080
    Now listening on: http://localhost:8080
    Application started. Press CTRL+C to shut down.
    Before the mainHandler: /
    Inside mainHandler
    After the mainHandler
```

全局使用中间件

```go
    package main
    import "github.com/kataras/iris"
    func main() {
        app := iris.New()
        //将“before”处理程序注册为将要执行的第一个处理程序
        //在所有域的路由上。
        //或使用`UseGlobal`注册一个将跨子域触发的中间件。
        app.Use(before)

        //将“after”处理程序注册为将要执行的最后一个处理程序
        //在所有域的路由'处理程序之后。
        app.Done(after)

        // register our routes.
        app.Get("/", indexHandler)
        app.Get("/contact", contactHandler)

        app.Run(iris.Addr(":8080"))
    }
    func before(ctx iris.Context) {
         // [...]
    }
    func after(ctx iris.Context) {
        // [...]
    }
    func indexHandler(ctx iris.Context) {
        ctx.HTML("<h1>Index</h1>")
        ctx.Next() // 执行通过`Done`注册的“after”处理程序。
    }
    func contactHandler(ctx iris.Context) {
        // write something to the client as a response.
        ctx.HTML("<h1>Contact</h1>")
        ctx.Next() // 执行通过`Done`注册的“after”处理程序。
    }
```

探索与发现
下面你可以看到一些有用的处理程序的源代码来学习

| Middleware名称          | 例子地址  
| :----:                 | :----: 
| basic authentication   | [basicauth](https://github.com/kataras/iris/tree/master/_examples/authentication/ )
| Google reCAPTCHA       | [recaptcha](https://github.com/kataras/iris/tree/master/_examples/miscellaneous/recaptcha)
| ocalization and internationalization | [i81n](https://github.com/kataras/iris/tree/master/_examples/miscellaneous/i81n)
| request logger    | [request-logger](https://github.com/kataras/iris/tree/master/_examples/http_request/request-logger) 
| article_id     |[profiling (pprof)](https://github.com/kataras/iris/tree/master/_examples/miscellaneous/pprof) 
| article_id     | [recovery](https://github.com/kataras/iris/tree/master/_examples/miscellaneous/recover) 


一些真正帮助您完成特定任务的中间件  

| Middleware名称     | 描述            | 例子地址  
| :----:            | :----:         | :----: 
| jwt               | 中间件在传入请求的Authorization标头上检查JWT并对其进行解码。|[iris-contrib/middleware/jwt/_example](https://github.com/iris-contrib/middleware/tree/master/jwt/_example )
| cors	            |HTTP访问控制。| [iris-contrib/middleware/cors/_example](https://github.com/iris-contrib/middleware/tree/master/cors/_example)
| secure            |实现一些快速安全性的中间件获胜。| [iris-contrib/middleware/secure/_example](https://github.com/iris-contrib/middleware/tree/master/secure/_example/main.go)
| tollbooth         |用于限制HTTP请求的通用中间件。| [iris-contrib/middleware/tollbooth/_examples/limit-handler](https://github.com/iris-contrib/middleware/tree/master/tollbooth/_examples/limit-handler) 
| cloudwatch        |AWS cloudwatch指标中间件。|[iris-contrib/middleware/cloudwatch/_example](https://github.com/iris-contrib/middleware/tree/master/cloudwatch/_example) 
| new relic         |官方New Relic Go Agent。| [iris-contrib/middleware/newrelic/_example](https://github.com/iris-contrib/middleware/tree/master/newrelic/_example) 
| prometheus        |轻松为prometheus检测工具创建指标端点| [	iris-contrib/middleware/prometheus/_example](https://github.com/iris-contrib/middleware/tree/master/prometheus/_example) 
| casbin            |一个授权库，支持ACL，RBAC，ABAC等访问控制模型| [iris-contrib/middleware/casbin/_examples](https://github.com/iris-contrib/middleware/tree/master/casbin/_examples)