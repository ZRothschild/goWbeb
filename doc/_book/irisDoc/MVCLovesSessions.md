#### MVC 结合session

```go
    // +build go1.9
    package main
    import (
        "fmt"
        "time"
        "github.com/kataras/iris"
        "github.com/kataras/iris/mvc"
        "github.com/kataras/iris/sessions"
    )
    // VisitController处理根路由。.
    type VisitController struct {
        //当前请求 Session
        //它的初始化是由我们添加到`visitApp`的依赖函数发生的。
        Session *sessions.Session

        //从MVC绑定的time.time，
        //绑定字段的顺序无关紧要。
        StartTime time.Time
    }
    // Get handles
    // Method: GET
    // Path: http://localhost:8080
    func (c *VisitController) Get() string {
       //它将整数的“访问”值递增1，
        //如果带有“访问”键的条目不存在，它将为您创建。
        visits := c.Session.Increment("visits", 1)
        // write the current, updated visits.
        since := time.Now().Sub(c.StartTime).Seconds()
        return fmt.Sprintf("%d visit from my current session in %0.1f seconds of server's up-time",visits, since)
    }
    func newApp() *iris.Application {
        app := iris.New()
        sess := sessions.New(sessions.Config{Cookie: "mysession_cookie_name"})
        visitApp := mvc.New(app.Party("/"))
       //将当前* session.Session绑定到`VisitController.Session`
        //和time.Now（）到`VisitController.StartTime`。
        visitApp.Register(
                //如果依赖是一个接受的函数
                 //一个Context并返回一个值
                 //然后控制器解析此函数的结果类型
                 //并且在每个请求上它将使用Context调用该函数
                 //并将结果（此处为* sessions.Session）设置为控制器的字段。
                 //如果没有字段或函数的输入参数，则注册依赖项
                 //使用者然后在服务器运行之前忽略这些依赖项，
                 //这样你就可以绑定很多dependecies并在不同的控制器中使用它们。
            sess.Start,
            time.Now(),
        )
        visitApp.Handle(new(VisitController))
        return app
    }
    func main() {
        app := newApp()
        // 1.打开浏览器（不在私人模式下,如google的隐身模式）
        // 2.导航到http：// localhost：8080
        // 3.刷新页面一些次数
        // 4.关闭浏览器
        // 5.重新打开浏览器并重新播放2
        app.Run(iris.Addr(":8080"), iris.WithoutVersionChecker)
    }
```