package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/netutil"
)

func main() {
	app := iris.New()
	l, err := netutil.UNIX("/tmpl/srv.sock", 0666) //查看其代码以了解如何手动创建新的文件侦听器，这很容易
	if err != nil {
		panic(err)
	}
	app.Run(iris.Listener(l))
}
//更多参阅 "customListener/unixReuseport"