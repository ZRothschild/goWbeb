####路由和反向查找

正如Handlers章节中所提到的，Iris提供了几种处理程序注册方法，每种方法都返回一个Route实例。

路由命名
 路由命名很容易,我们只调用返回的*Route，并使用Name字段来定义一个名称:
反向路径
 反向路径,也就是从路径名生成url
 ```go
    package main
    import (
        "github.com/kataras/iris"
    )
    func main() {
        app := iris.New()
        // define a function
        h := func(ctx iris.Context) {
            ctx.HTML("<b>Hi</b1>")
        }
        // handler registration and naming
        home := app.Get("/", h)
        home.Name = "home"
        // or
        app.Get("/about", h).Name = "about"
        app.Get("/page/{id}", h).Name = "page"

        app.Run(iris.Addr(":8080"))
    }
```
 当我们为特定路径注册处理程序时，我们就能够基于传递给Iris的结构化数据创建url。在上面的示例中，我们命名了三个路由器，其中一个甚至使用了参数。如果我们使用默认的html/模板视图引擎，我们可以使用一个简单的操作来逆转路由(以及生成实际的url):

 ```go
    Home: {{ urlpath "home" }}
    About: {{ urlpath "about" }}
    Page 17: {{ urlpath "page" "17" }}
```
上面的代码将生成以下输出：
> Home: http://localhost:8080/
> About: http://localhost:8080/about
> Page 17: http://localhost:8080/page/17

我们可以使用以下方法/函数来处理命名路由（及其参数）：

GetRoutes函数用于获取所有已注册的路由
GetRoute（routeName string）方法按名称检索路由
URL（routeName string，paramValues ... interface {}）方法，用于根据提供的参数生成url字符串
Path（routeName string，paramValues ... interface {}方法，根据提供的值生成URL的路径（没有主机和协议）部分