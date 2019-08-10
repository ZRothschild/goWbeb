####使用路由

#####基本介绍

Iris 支持所有HTTP方法，开发人员还可以为不同方法注册相同路径的处理程序。

第一个参数是HTTP方法，第二个参数是路径的请求路径，
第三个可变参数应该包含一个或多个iris.Handler，
当用户从服务器请求该特定的资源路径时，由注册的顺序执行。

示例代码:

```go
    package main
    import (
        "github.com/kataras/iris"
    )
    func main(){
        app := iris.New()
        app.Handle("GET", "/contact", func(ctx iris.Context) {
            ctx.HTML("<h1> Hello from /contact </h1>")
        })
    }
```
为了使最终开发人员更容易，iris为所有HTTP方法提供了功能。第一个参数是路由的请求路径，
第二个可变参数应该包含一个或多个iris.Handler，当用户从服务器请求该特定的资源路径时，由注册顺序执行。

示例代码:

```go
    package main
    import (
        "github.com/kataras/iris"
    )
    func main(){
        app := iris.New()
        //GET 方法
        app.Get("/", handler)
        // POST 方法
        app.Post("/", handler)
        // PUT 方法
        app.Put("/", handler)
        // DELETE 方法
        app.Delete("/", handler)
        //OPTIONS 方法
        app.Options("/", handler)
        //TRACE 方法
        app.Trace("/", handler)
        //CONNECT 方法
        app.Connect("/", handler)
        //HEAD 方法
        app.Head("/", handler)
        // PATCH 方法
        app.Patch("/", handler)
        //任意的http请求方法如option等
        app.Any("/", handler)

    }
    func handler(ctx iris.Context){
        ctx.Writef("Hello from method: %s and path: %s", ctx.Method(), ctx.Path())
    }
```    
分组路由
由路径前缀分组的一组路由可以（可选）共享相同的中间件处理程序和模板布局。一个组也可以有一个嵌套组。

.Party 正在用于分组路由，开发人员可以声明无限数量的（嵌套）组。

```go
    package main
    import (
        "github.com/kataras/iris"
    )
    func main(){
        app := iris.New()
        //请在参数化路径部分
        users := app.Party("/users", myAuthMiddlewareHandler)
        // http://localhost:8080/users/42/profile
        users.Get("/{id:int}/profile", userProfileHandler)
        // http://localhost:8080/users/inbox/1
        users.Get("/inbox/{id:int}", userMessageHandler)
    }
    func myAuthMiddlewareHandler(ctx iris.Context){
        ctx.WriteString("Authentication failed")
    }
    func userProfileHandler(ctx iris.Context) {//
        id:=ctx.Params().Get("id")
        ctx.WriteString(id)
    }
    func userMessageHandler(ctx iris.Context){
        id:=ctx.Params().Get("id")
        ctx.WriteString(id)
    }
```
也可以使用接受子路由器（Party）的功能编写相同的内容。
```go
    package main
    import (
        "github.com/kataras/iris"
    )
    func main(){
        app := iris.New()
        app.PartyFunc("/users", func(users iris.Party) {
            users.Use(myAuthMiddlewareHandler)
            // http://localhost:8080/users/42/profile
            users.Get("/{id:int}/profile", userProfileHandler)
            // http://localhost:8080/users/messages/1
            users.Get("/inbox/{id:int}", userMessageHandler)
        })
    }
    func myAuthMiddlewareHandler(ctx iris.Context){
        ctx.WriteString("Authentication failed")
        ctx.Next()//继续执行后续的handler
    }
    func userProfileHandler(ctx iris.Context) {//
        id:=ctx.Params().Get("id")
        ctx.WriteString(id)
    }
    func userMessageHandler(ctx iris.Context){
        id:=ctx.Params().Get("id")
        ctx.WriteString(id)
    }
```

