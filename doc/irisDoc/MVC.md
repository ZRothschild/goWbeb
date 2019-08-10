####MVC介绍
![run2](images\web_mvc_diagram.png)

Iris对MVC(模型视图控制器)模式有一流的支持，您在其他任何地方都找不到这些东西

Iris web框架支持请求数据、模型、持久数据和以最快的速度执行的绑定。

#####特性

支持所有HTTP方法，例如，如果想要提供GET，那么控制器应该有一个名为Get（）的函数，您可以定义多个方法函数在同一个Controller中提供。

通过BeforeActivation自定义事件回调，每个控制器，将自定义控制器的struct的方法作为具有自定义路径（即使使用正则表达式参数化路径）的处理程序提供。

示例代码:
```go
    import (
        "github.com/kataras/iris"
        "github.com/kataras/iris/mvc"
    )
    func main() {
        app := iris.New()
        mvc.Configure(app.Party("/root"), myMVC)
        app.Run(iris.Addr(":8080"))
    }
    func myMVC(app *mvc.Application) {
        // app.Register(...)
        // app.Router.Use/UseGlobal/Done(...)
        app.Handle(new(MyController))
    }
    type MyController struct {}
    func (m *MyController) BeforeActivation(b mvc.BeforeActivation) {
        // b.Dependencies().Add/Remove
        // b.Router().Use/UseGlobal/Done // 以及您已经知道的任何标准API调用

        // 1-> Method
        // 2-> Path
        // 3-> 控制器的函数名称将被解析为处理程序
        // 4-> 应该在MyCustomHandler之前运行的任何处理程序
        b.Handle("GET", "/something/{id:long}", "MyCustomHandler", anyMiddleware...)
    }
    // GET: http://localhost:8080/root
    func (m *MyController) Get() string { return "Hey" }
    // GET: http://localhost:8080/root/something/{id:long}
    func (m *MyController) MyCustomHandler(id int64) string { return "MyCustomHandler says Hey" }
```

Controller结构中的持久性数据（在请求之间共享数据），通过定义对依赖项的服务或具有Singleton控制器作用域。

共享控制器之间的依赖关系或在父MVC应用程序上注册它们，并能够在Controller内的BeforeActivation可选事件回调中修改每个控制器的依赖关系，即
func（c * MyController）BeforeActivation（b mvc.BeforeActivation）{b.Dependencies （）add/remove（...）}
。

访问Context作为控制器的字段（没有手动绑定是neede），即Ctx iris.Context或通过方法的输入参数，即
func（ctx iris.Context，otherArguments ...）。

Controller结构中的模型（在Method函数中设置并由View呈现）。
您可以从控制器的方法返回模型，或者在请求生命周期中设置字段，并在同一请求生命周期中将该字段返回到另一个方法。
像以前一样，mvc应用程序有自己的路由器，这是一种iris/route.Party，标准的 iris api。控制器可以注册到任何一方，包括子域名，Party的开始和完成处理程序按预期工作。


可选的BeginRequest（ctx）函数在方法执行之前执行任何初始化，对调用中间件或许多方法使用相同的数据集合很有用。
可选的EndRequest（ctx）函数，用于在执行任何方法后执行任何终结。

继承，递归，参见我们的mvc.SessionController，它将Session * sessions.Session和Manager * sessions.Sessions作为嵌入字段，由其BeginRequest填充，在这里。这只是一个示例，您可以使用从管理器的Start作为动态依赖关系返回到MVC应用程序的sessions.Session，即mvcApp.Register（sessions.New（sessions.Config {Cookie：“iris_session_id”}）。 ）。

通过控制器方法的输入参数访问动态路径参数，不需要绑定。当您使用Iris的默认语法来解析来自控制器的处理程序时，您需要使用By字来为方法添加后缀，大写是一个新的子路径。例：
 如这种形式 mvc.New(app.Party("/user")).Handle(new(user.Controller))
 则:
 * func(*Controller) Get() - GET:/user.
 * func(*Controller) Post() - POST:/user.
 * func(*Controller) GetLogin() - GET:/user/login
 * func(*Controller) PostLogin() - POST:/user/login
 * func(*Controller) GetProfileFollowers() - GET:/user/profile/followers
 * func(*Controller) PostProfileFollowers() - POST:/user/profile/followers
 * func(*Controller) GetBy(id int64) - GET:/user/{param:long}
 * func(*Controller) PostBy(id int64) - POST:/user/{param:long}

这样也是哦 mvc.New(app.Party("/profile")).Handle(new(profile.Controller))

* func(*Controller) GetBy(username string) - GET:/profile/{param:string}


mvc.New(app.Party("/assets")).Handle(new(file.Controller))

* func(*Controller) GetByWildard(path string) - GET:/assets/{param:path}

方法函数接收器支持的类型：int，int64，bool和string。

通过输出参数响应，可选，即
```go
func(c *ExampleController) Get() string |
                                (string, string) |
                                (string, int) |
                                int |
                                (int, string) |
                                (string, error) |
                                error |
                                (int, error) |
                                (any, bool) |
                                (customStruct, error) |
                                customStruct |
                                (customStruct, int) |
                                (customStruct, string) |
                                mvc.Result or (mvc.Result, error)
```

当mvc.Result是一个interface时,只包含该函数：：Dispatch（ctx iris.Context） 

#####使用Iris MVC进行代码重用

通过创建彼此独立的组件，开发人员能够在其他应用程序中快速轻松地重用组件。对于具有不同数据的另一个应用程序，可以为一个应用程序重构相同（或类似）的视图，因为视图只是处理数据如何显示给用户。

如果您不熟悉后端Web开发，请首先阅读有关MVC架构模式的内容，一个好的开始就是。
[维基百科文章](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller)

##### 快速MVC教程第1部分

此示例等同于[https://github.com/kataras/iris/blob/master/_examples/hello-world/main.go](https://github.com/kataras/iris/blob/master/_examples/hello-world/main.go)

似乎你必须编写额外的代码并不值得，但请记住，
这个例子没有使用iris mvc功能，比如Model，Persistence或View引擎都没有Session，
它对于学习目的来说非常简单，可能你在你的应用程序的任何地方都不会使用简单的控制器。

在我的个人笔记本电脑上，在每个20MB吞吐量的“/ hello”路径上使用MVC的这个例子的成本是每个20MB吞吐量大约2MB，
大多数应用程序都可以容忍，但你可以选择最适合你的Iris，低级处理程序：
性能或高级控制器：在大型应用程序上更易于维护和更小的代码库。

```go
    package main
    import (
        "github.com/kataras/iris"
        "github.com/kataras/iris/mvc"
        "github.com/kataras/iris/middleware/logger"
        "github.com/kataras/iris/middleware/recover"
    )
    func main() {
        app := iris.New()
        //（可选）添加两个内置处理程序
        // 可以从任何http相关的恐慌中恢复
        // 并将请求记录到终端。
        app.Use(recover.New())
        app.Use(logger.New())

        // 基于根路由器服务控制器， "/".
        mvc.New(app).Handle(new(ExampleController))

        // http://localhost:8080
        // http://localhost:8080/ping
        // http://localhost:8080/hello
        // http://localhost:8080/custom_path
        app.Run(iris.Addr(":8080"))
    }
    //ExampleController服务于 "/", "/ping" and "/hello".
    type ExampleController struct{}
    // Get serves
    // Method:   GET
    // Resource: http://localhost:8080
    func (c *ExampleController) Get() mvc.Result {
        return mvc.Response{
            ContentType: "text/html",
            Text:        "<h1>Welcome</h1>",
        }
    }
    // GetPing serves
    // Method:   GET
    // Resource: http://localhost:8080/ping
    func (c *ExampleController) GetPing() string {
        return "pong"
    }
    // GetHello serves
    // Method:   GET
    // Resource: http://localhost:8080/hello
    func (c *ExampleController) GetHello() interface{} {
        return map[string]string{"message": "Hello Iris!"}
    }
    //在控制器适应主应用程序之前调用一次BeforeActivation
    //当然在服务器运行之前。
    //在版本9之后，您还可以为特定控制器的方法添加自定义路由。
    //在这里您可以注册自定义方法的处理程序
    //使用带有`ca.Router`的标准路由器做一些你可以做的事情，没有mvc，
    //并添加将绑定到控制器的字段或方法函数的输入参数的依赖项。
    func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
        anyMiddlewareHere := func(ctx iris.Context) {
            ctx.Application().Logger().Warnf("Inside /custom_path")
            ctx.Next()
        }
        b.Handle("GET", "/custom_path", "CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)

        //甚至添加基于此控制器路由器的全局中间件，
        //在这个例子中是根“/”：
        // b.Router（）。使用（myMiddleware）
    }
    // CustomHandlerWithoutFollowingTheNamingGuide serves
    // Method:   GET
    // Resource: http://localhost:8080/custom_path
    func (c *ExampleController) CustomHandlerWithoutFollowingTheNamingGuide() string {
        return "hello from the custom handler without following the naming guide"
    }
    // GetUserBy serves
    // Method:   GET
    // Resource: http://localhost:8080/user/{username:string}
    //是一个保留的“关键字”来告诉框架你要去的
    //在函数的输入参数中绑定路径参数，它也是
    //有助于在同一控制器中使用“Get”和“GetBy”。
    //
    // func (c *ExampleController) GetUserBy(username string) mvc.Result {
    //     return mvc.View{
    //         Name: "user/username.html",
    //         Data: username,
    //     }
    // }
    /*
    可以使用多个，工厂会确定
    为每条路线注册了正确的http方法
    对于此控制器，如果需要，请取消注释：
    */
    func (c *ExampleController) Post() {}
    func (c *ExampleController) Put() {}
    func (c *ExampleController) Delete() {}
    func (c *ExampleController) Connect() {}
    func (c *ExampleController) Head() {}
    func (c *ExampleController) Patch() {}
    func (c *ExampleController) Options() {}
    func (c *ExampleController) Trace() {}
    */
    /*
    func (c *ExampleController) All() {}
    //        OR
    func (c *ExampleController) Any() {}

    func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
        // 1 -> the HTTP Method
        // 2 -> the route's path
        // 3 -> this controller's method name that should be handler for that route.
        b.Handle("GET", "/mypath/{param}", "DoIt", optionalMiddlewareHere...)
    }
    // After activation, all dependencies are set-ed - so read only access on them
    // but still possible to add custom controller or simple standard handlers.
    func (c *ExampleController) AfterActivation(a mvc.AfterActivation) {}
    */
```
>_examples / mvc和mvc / controller_test.go
文件用简单的范例解释每个功能，
他们展示了如何利用Iris MVC Binder，Iris MVC模型以及更多......

在控制器中以HTTP方法（Get，Post，Put，Delete ...）为前缀的每个导出的func都可以作为HTTP端点调用。在上面的示例中，所有func都将一个字符串写入响应。请注意每种方法之前的注释。

HTTP端点是Web应用程序中的可定位URL，例如http：// localhost：8080 / helloworld，并结合使用的协议：HTTP，Web服务器的网络位置（包括TCP端口）：localhost：8080和目标URI / helloworld。

第一条评论声明这是一个HTTP GET方法，它通过将“/ helloworld”附加到基本URL来调用。第三个注释指定通过将“/ helloworld / welcome”附加到URL来调用的HTTP GET方法。

Controller知道如何处理GetBy上的“name”或GetWelcomeBy中的“name”和“numTimes”，因为By关键字，并且构建没有样板的动态路由;第三个注释指定HTTP GET动态方法，
该方法由任何以“/ helloworld / welcome”开头并后跟两个路径部分的URL调用，第一个可以接受任何值，第二个只能接受数字，
i，e ：“http：// localhost：8080 / helloworld / welcome / golang / 32719”，
否则将向客户端发送404 Not Found HTTP Error。

单击[movieMVC](https://docs.iris-go.com/mvc_3.html)到“movieMVC应用程序”子部分。