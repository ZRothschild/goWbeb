# `IRIS`自定义`http`服务
## 简单服务
### 目录结构
> 主目录`easyWay`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"net/http"
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
	//这里有任何自定义字段 Handler和ErrorLog自动设置到服务器
	srv := &http.Server{Addr: ":8080"}
	// http://localhost:8080/
	// http://localhost:8080/mypath
	app.Run(iris.Server(srv)) // 等同于 app.Run(iris.Addr(":8080"))
	// 更多：
	//如果您需要在同一个应用程序中使用多个服务器，请参阅“multi”。
	//用于自定义listener：iris.Listener（net.Listener）或
	// iris.TLS(cert，key)或iris.AutoTLS()，请参阅“custom-listener”示例。
}
```
## 多服务
### 目录结构
> 主目录`multi`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"net/http"
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
	//注意：如果第一个动作是“go app.Run”，则不需要它。
	if err := app.Build(); err != nil {
		panic(err)
	}
	//启动侦听localhost：9090的辅助服务器。
	//如果您需要在同一个应用程序中使用多个服务器，请对Listen函数使用“go”关键字。
	// http://localhost:9090/
	// http://localhost:9090/mypath
	srv1 := &http.Server{Addr: ":9090", Handler: app}
	go srv1.ListenAndServe()
	println("Start a server listening on http://localhost:9090")
	//启动一个“second-secondary”服务器，监听localhost：5050。
	// http://localhost:5050/
	// http://localhost:5050/mypath
	srv2 := &http.Server{Addr: ":5050", Handler: app}
	go srv2.ListenAndServe()
	println("Start a server listening on http://localhost:5050")
	//注意：app.Run完全是可选的，我们已经使用app.Build构建了应用程序，
	//你可以改为创建一个新的http.Server。
	// http://localhost:8080/
	// http://localhost:8080/mypath
	app.Run(iris.Addr(":8080")) //在这里以后就不可以监听了
}
```
## `iris`服务与原始服务共存
### 目录结构
> 主目录`stdWay`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"net/http"
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
	//在自定义http.Server上使用'app'作为http.Handler之前调用.Build
	app.Build()
	//创建我们的自定义服务器并分配Handler/Router
	srv := &http.Server{Handler: app, Addr: ":8080"} //你必须设置Handler:app和Addr,请参阅“iris-way”,它会自动执行此操作
	// http://localhost:8080/
	// http://localhost:8080/mypath
	println("Start a server listening on http://localhost:8080")
	srv.ListenAndServe() // 等同于 app.Run(iris.Addr(":8080"))
	//注意：
	//根本不显示所有。 即使应用程序的配置允许，中断处理程序也是如此。
	//`.Run`是唯一一个关心这三者的函数。
	// 更多：
	//如果您需要在同一个应用程序中使用多个服务器，请参阅“multi”。
	//用于自定义侦听器：iris.Listener(net.Listener)或
	// iris.TLS(cert,key)或iris.AutoTLS()，请参阅“custom-listener”示例。
}
```