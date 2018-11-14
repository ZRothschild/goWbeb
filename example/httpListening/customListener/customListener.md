# `IRIS`自定义监听器
## 简单监听器
### 目录结构
> 主目录`customListener`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"net"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello from the server")
	})
	app.Get("/mypath", func(ctx iris.Context) {
		ctx.Writef("Hello from %s", ctx.Path())
	})
	//创建任何自定义tcp侦听器，unix sock文件或tls tcp侦听器。
	l, err := net.Listen("tcp4", ":8080")
	if err != nil {
		panic(err)
	}
	//使用自定义侦听器
	app.Run(iris.Listener(l))
}
```
## `linux`监听器
### 目录结构
> 主目录`unixReuseport`

```html
    —— main.go
    —— main_windows.go
```
### 代码示例
> `main.go`

```go
// +build linux darwin dragonfly freebsd netbsd openbsd rumprun
package main

import (
	//包tcplisten提供各种可自定义的TCP net.Listener
	//与性能相关的选项：
	// - SO_REUSEPORT。 此选项允许线性扩展服务器性能 在多CPU服务器上。
	//有关详细信息，请参阅https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/。
	// - TCP_DEFER_ACCEPT。 此选项期望服务器从接受的读取写入之前的连接。
	// - TCP_FASTOPEN。 有关详细信息，请参阅https://lwn.net/Articles/508865/。
	"github.com/valyala/tcplisten"
	"github.com/kataras/iris"
)
//注意只支持linux系统
// $ go get github.com/valyala/tcplisten
// $ go run main.go

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello World!</b>")
	})
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
> `main_windows.go`

```go
// +build windows
package main

func main() {
	panic("windows operating system does not support this feature")
}
```