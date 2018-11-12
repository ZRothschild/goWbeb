# Hosts

## Listen and Serve

您可以启动服务器监听任何类型的`net.Listener`甚至是`http.Server`实例.应该通过`Run`函数在最后传递服务器的初始化方法。

Go开发人员用于服务其服务器的最常用方法是通过传递“hostname：ip”形式的网络地址.有了iris,我们使用的是`iris.Runner`类型的`iris.Addr`

```go
//监听tcp网络地址0.0.0.0:8080
app.Run(iris.Addr(":8080"))
```

有时您在应用程序的其他位置创建了标准的net/http服务器，并希望使用它来为Iris Web应用程序提供服务
```go
//与之前相同，但使用自定义的http.Server，也可能在其他地方使用
app.Run(iris.Server(&http.Server{Addr:":8080"}))
```

最高级的用法是创建一个自定义或标准的`net.Listener`并将其传递给`app.Run`

```go
//使用自定义net.Listener
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
UNIX和BSD主机可以利用重用端口功能
```go
package main

import (
	//tcplisten包 提供各种可自定义的TCP net.Listener 与 性能相关的选项 Linux 特性
	//第一 SO_REUSEPORT。 此选项允许线性扩展服务器性能 在多CPU服务器上。
	//关详细信息，请参阅 https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/
	//第二 TCP_DEFER_ACCEPT。 此选项期望服务器从接受的读取写入之前的连接
	//第三 TCP_FASTOPEN 关详细信息，请参阅https://lwn.net/Articles/508865/。
	"github.com/valyala/tcplisten"
	"github.com/kataras/iris"
)
// 安装 tcplisten $ go get github.com/valyala/tcplisten
// $ go run main.go
func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello World!</b>")
	})
	//对用上面的三个选项
	listenerCfg := tcplisten.Config{
		ReusePort:   true,
		DeferAccept: true,
		FastOpen:    true,
	}
	l, err := listenerCfg.NewListener("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	app.Run(iris.Listener(l))
}
```

### HTTP/2 and Secure

如果您有签名文件密钥，您可以使用`iris.TLS`根据这些证书密钥提供https

```go
// TLS using files
app.Run(iris.TLS("127.0.0.1:443", "mycert.cert", "mykey.key"))
```
你的应用程序准备好**生产**时应该使用的方法是`iris.AutoTLS`，它启动一个安全的服务器，自动认证由https://letsencrypt.org提供**免费**

```go
//自动TLS
app.Run(iris.AutoTLS(":443", "example.com", "admin@example.com"))
```

### Any `iris.Runner`

有时你可能想要一些非常特别的东西来听，这不是一种`net.Listener`。 你可以通过`iris.Raw`来做到这一点，但是你要设置这个方法

```go
//使用任何func() error，
//启动listener责任取决于你使用的方式，
//为了简单起见，我们将使用
//`net/http`包的ListenAndServe函数。
app.Run(iris.Raw(&http.Server{Addr:":8080"}).ListenAndServe)
```

## Host configurators

所有上述形式的listening都接受了`func(* iris.Supervisor)`的最后一个变量参数。 这用于为通过这些函数传递的特定主机添加配置程序。

例如，假设我们要添加一个在何时触发的回调
服务器已关闭

```go
app.Run(iris.Addr(":8080", func(h *iris.Supervisor) {
    h.RegisterOnShutdown(func() {
        println("server terminated")
    })
}))
```

你甚至可以在`app.Run`方法之前做到这一点，但区别在于,这些主机配置程序将被执行到您可能用于为您的
Web应用程序提供服务的所有主机（通过`app.NewHost`我们将在一分钟内看到）

```go
app := iris.New()
app.ConfigureHost(func(h *iris.Supervisor) {
    h.RegisterOnShutdown(func() {
        println("server terminated")
    })
})
app.Run(iris.Addr(":8080"))
```

可以访问为您的应用程序提供服务的所有主机，在`Run`方法之后的`Application＃Hosts`字段。

但最常见的情况是您可能需要在`app.Run`方法之前访问主机，有两种方法可以访问主机主管，如下所示。

我们已经看到了如何通过`app.Run`或`app.ConfigureHost`的第二个参数配置所有应用程序的主机。 还有一种更适合简单场景的方法，那就是使用`app.NewHost`创建一个新主机
并使用其中一个“Serve”或“Listen”函数通过`iris＃Raw` Runner启动应用程序。

请注意，这种方式需要额外导入`net / http`包，示例代码：

```go
h := app.NewHost(&http.Server{Addr:":8080"})
h.RegisterOnShutdown(func(){
    println("server terminated")
})
app.Run(iris.Raw(h.ListenAndServe))
```

## Multi hosts

您可以使用多个服务为您的`Iris Web`应用程序提供服务,`iris.Router`与`net/http/Handler`功能兼容,因此您可以理解,它可以用于任何
`net/http`服务器,但有一种更简单的方法,使用`app.NewHost`,它也复制所有主机配置器,并关闭连接到`app.Shutdown`上特定`Web`应用程序的所有主机

```go
app := iris.New()
app.Get("/", indexHandler)
//在不同的goroutine中运行，防止阻塞goroutine
go app.Run(iris.Addr(":8080"))
//启动第二个正在侦听tcp 0.0.0.0:9090的服务器，
//没有“go”关键字，因为我们想要在最后一次服务器运行时阻止。
app.NewHost(&http.Server{Addr:":9090"}).ListenAndServe()
```
## Shutdown (Gracefully)

让我们继续学习如何捕获`CONTROL+C/COMMAND+C`或`unix kill`命令并优雅地关闭服务器

> 正确关闭`CONTROL+C/COMMAND+C`或者当发送的`kill`命令是`ENABLED BY-DEFAULT`时

为了手动管理应用程序中断时要执行的操作，我们必须使用选项`WithoutInterruptHandler`禁用默认行为
并注册一个新的中断处理程序（全局，跨所有的主机）

示例代码:
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