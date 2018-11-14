# `IRIS`优雅关闭服务
## 自定义通知关闭服务
### 目录结构
> 主目录`customNotifier`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	stdContext "context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1>hi, I just exist in order to see if the server is closed</h1>")
	})
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			println("shutdown...")
			timeout := 5 * time.Second
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
			defer cancel()
			app.Shutdown(ctx)
		}
	}()
	//启动服务器并禁用默认的中断处理程序
	//我们自己处理它清晰简单，没有任何问题。
	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}
```
## 默认通知关闭服务
### 目录结构
> 主目录`defaultNotifier`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	stdContext "context"
	"time"
	"github.com/kataras/iris"
)
//继续之前：
//正常关闭control+C/command+C或当发送的kill命令是ENABLED BY-DEFAULT。
//为了手动管理应用程序中断时要执行的操作，
//我们必须使用选项`WithoutInterruptHandler`禁用默认行为
//并注册一个新的中断处理程序(全局，所有主机)。
func main() {
	app := iris.New()
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		//关闭所有主机
		app.Shutdown(ctx)
	})
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(" <h1>hi, I just exist in order to see if the server is closed</h1>")
	})
	// http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}
```