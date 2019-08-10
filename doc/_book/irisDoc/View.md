####视图

Iris支持开箱即用的5个模板引擎，开发人员仍然可以使用任何外部golang模板引擎，
如 context/context#ResponseWriter() is an io.Writer.

所有这五个模板引擎都具有通用API的共同特征，如布局，模板功能，特定于派对的布局，部分渲染等

标准的html,它的模板解析器就是 [golang.org/pkg/html/template/](https://golang.org/pkg/html/template/)
Django,它的模板解析器就是 [github.com/flosch/pongo2](https://github.com/flosch/pongo2)
Pug(Jade),它的模板解析器就是 [github.com/Joker/jade](https://github.com/Joker/jade)
Handlebars, 它的模板解析器 [github.com/aymerick/raymond](https://github.com/aymerick/raymond)
Amber, 它的模板解析器 [github.com/eknkc/amber](https://github.com/eknkc/amber)

概述
```go
    package main
    import "github.com/kataras/iris"
    func main() {
        app := iris.New()
        // 从 "./views" 文件夹加载所以的模板
        // 其中扩展名为“.html”并解析它们
        // 使用标准的`html / template`包。
        app.RegisterView(iris.HTML("./views", ".html"))

        // 方法:    GET
        // 资源:  http://localhost:8080
        app.Get("/", func(ctx iris.Context) {
            //绑定数据
            ctx.ViewData("message", "Hello world!")
            // 渲染视图文件: ./views/hello.html
            ctx.View("hello.html")
        })
         // 方法:    GET
        //资源:  http://localhost:8080/user/42
        app.Get("/user/{id:long}", func(ctx iris.Context) {
            userID, _ := ctx.Params().GetInt64("id")
            ctx.Writef("User ID: %d", userID)
        })
        //启动服务
        app.Run(iris.Addr(":8080"))
    }
```