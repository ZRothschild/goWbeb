### Hosts

监听服务

您可以启动服务器监听任何类型的`net.Listener`甚至`http.Server`实例。
服务器的初始化方法应该在最后通过Run函数传递。

Go开发人员用于服务其服务器的最常用方法是传递“hostname：ip”形式的网络地址。
有了Iris，我们使用的iris.Addr是一种iris.Runner类型

//用网络地址监听tcp 0.0.0.0:8080

app.Run(iris.Addr(":8080"))

有时您在应用程序的其他位置创建了标准的net / http服务器，并希望使用它来为Iris Web应用程序提供服务

// 与之前相同，但使用自定义的http.Server，也可能在其他地方使用

app.Run(iris.Server(&http.Server{Addr:":8080"}))

最高级的用法是创建自定义或标准net.Listener并将其传递给app.Run

// 使用自定义的net.Listener

```go
    l, err := net.Listen("tcp4", ":8080")
    if err != nil {
        panic(err)
    }
    app.Run(iris.Listener(l))
```
一个更完整的示例，使用仅限unix的套接字文件功能
```go
    package main
    import (
        "os"
        "net"
        "github.com/kataras/iris"
    )
    func main() {
        app := iris.New()
        // UNIX socket
        if errOs := os.Remove(socketFile); errOs != nil && !os.IsNotExist(errOs) {
            app.Logger().Fatal(errOs)
        }
        l, err := net.Listen("unix", socketFile)
        if err != nil {
            app.Logger().Fatal(err)
        }
        if err = os.Chmod(socketFile, mode); err != nil {
            app.Logger().Fatal(err)
        }
        app.Run(iris.Listener(l))
    }

```
UNIX和BSD主机可以优先考虑重用端口功能

```go
    package main
    import (
        // Package tcplisten provides customizable TCP net.Listener with various
        // performance-related options:
        //
        //   - SO_REUSEPORT. This option allows linear scaling server performance
        //     on multi-CPU servers.
        //     See https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/ for details.
        //
        //   - TCP_DEFER_ACCEPT. This option expects the server reads from the accepted
        //     connection before writing to them.
        //
        //   - TCP_FASTOPEN. See https://lwn.net/Articles/508865/ for details.
        "github.com/valyala/tcplisten"
        "github.com/kataras/iris"
    )
   //go get github.com/valyala/tcplisten
   //go run main.go
    func main() {
        app := iris.New()
        app.Get("/", func(ctx iris.Context) {
            ctx.HTML("<h1>Hello World!</h1>")
        })
        listenerCfg := tcplisten.Config{
            ReusePort:   true,
            DeferAccept: true,
            FastOpen:    true,
        }
        l, err := listenerCfg.NewListener("tcp", ":8080")
        if err != nil {
            app.Logger().Fatal(err)
        }
        app.Run(iris.Listener(l))
    }
```
#### HTTP / 2和安全
如果您有已签名的文件密钥，则可以根据这些证书密钥使用该iris.TLS服务https
```go
    app.Run(iris.TLS("127.0.0.1:443", "mycert.cert", "mykey.key"))
```
该方法时，你的应用程序已经准备好，
你应该使用的生产是iris.AutoTLS其开始安全服务器所提供的自动认证
[https://letsencrypt.org](https://letsencrypt.org)为免费

// Automatic TLS
```go
    app.Run(iris.AutoTLS(":443", "example.com", "admin@example.com"))
```


任何 iris.Runner

有时你可能想要一些非常特别的东西来听，这不是一种类型的`net.Listener`。你能够做到这一点iris.Raw，但你负责这种方法

//使用任何func（）错误，
//启动听众的责任取决于你这个方式，
//为了简单起见，我们将使用
//`net / http`包的ListenAndServe函数
```go
    app.Run(iris.Raw(&http.Server{Addr:":8080"}).ListenAndServe)
```

主机配置器

所有上述形式的倾听都接受了最后的，可变的论证func(*iris.Supervisor)。这用于为通过这些函数传递的特定主机添加配置程序。

例如，假设我们要添加一个在服务器关闭时触发的回调

```go
    app.Run(iris.Addr(":8080", func(h *iris.Supervisor) {
        h.RegisterOnShutdown(func() {
            println("server terminated")
        })
    }))
```

您甚至可以在app.Run方法之前执行此操作，但区别在于这些主机配置程序将被执行到您可能用于为您的Web应用程序提供服务的所有主机
（通过app.NewHost我们将在一分钟内看到）

```go
    app := iris.New()
    app.ConfigureHost(func(h *iris.Supervisor) {
        h.RegisterOnShutdown(func() {
            println("server terminated")
        })
    })
    app.Run(iris.Addr(":8080"))
```
Application#Hosts在该Run方法之后，字段可以提供对为您的应用程序提供服务的所有主机的访问权限。

但最常见的情况是您可能需要在app.Run方法之前访问主机，有两种获取访问主机主管的方法，请参阅下文。

我们已经看到了如何通过app.Run或的第二个参数配置所有应用程序的主机app.ConfigureHost。还有一种方法更适合简单的场景，即使用app.NewHost创建新主机并使用其中一个Serve或多个Listen函数通过iris#RawRunner 启动应用程序。

请注意，这种方式需要额外导入net/http包。

```go
    h := app.NewHost(&http.Server{Addr:":8080"})
    h.RegisterOnShutdown(func(){
        println("server terminated")
    })
    app.Run(iris.Raw(h.ListenAndServe))
```
####多主机

您可以使用多个服务器为您的Iris Web应用程序提供服务，因此iris.Router与net/http/Handler功能兼容，您可以理解，它可以在任何net/http服务器上进行调整，但是有一种更简单的方法，
通过使用app.NewHost也可以复制所有主机配置程序，它关闭连接到特定Web应用程序的所有主机app.Shutdown。

```go
    app := iris.New()
    app.Get("/", indexHandler)

    //在不同的goroutine中运行，以便不阻止主要的“goroutine”。
    go app.Run(iris.Addr(":8080"))

    // 启动第二个服务器，它正在监听tcp 0.0.0.0:9090，
    //没有“go”关键字，因为我们想要在最后一次服务器运行时阻止。
    app.NewHost(&http.Server{Addr:":9090"}).ListenAndServe()
```
关机（优雅）
让我们继续学习如何捕获CONTROL + C / COMMAND + C或unix kill命令并优雅地关闭服务器。

> 正确关闭CONTROL + C / COMMAND + C或者当发送的kill命令是ENABLED BY-DEFAULT时。

为了手动管理应用程序中断时要执行的操作，我们必须使用该选项禁用默认行为WithoutInterruptHandler 并注册新的中断处理程序（全局，跨所有可能的主机）。

如下代码:
```go
    package main
    import (
        "context"
        "time"
        "github.com/kataras/iris"
    )
    func main() {
        app := iris.New()
        iris.RegisterOnInterrupt(func() {
            timeout := 5 * time.Second
            ctx, cancel := context.WithTimeout(context.Background(), timeout)
            defer cancel()
            // close all hosts
            app.Shutdown(ctx)
        })
        app.Get("/", func(ctx iris.Context) {
            ctx.HTML(" <h1>hi, I just exist in order to see if the server is closed</h1>")
        })
        app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
    }
```

