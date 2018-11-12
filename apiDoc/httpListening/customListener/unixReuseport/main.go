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