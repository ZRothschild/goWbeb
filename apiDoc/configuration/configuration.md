# 配置组
所有配置的值都有默认值，所有你可以直接使用`iris.New()`

在`listen`函数之前配置是没用的，所以应该在`Application＃Run / 2`（第二个参数）上传递它

`Iris`有一个名为`Configurator`的类型，它是一个`func(* iris.Application)`，任何函数
完成此操作可以在`Application＃Configure`或`Application＃Run/2`中传递。

`Application#ConfigurationReadOnly()`返回配置值.
## 加载配置类型介绍
### 通过结构体加载配置
#### 目录结构
> 主目录`fromConfigurationStructure`
```html
    —— main.go
```
#### 代码示例
> `main.go`
```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})
	// [...]
	//当您想要修改整个配置时非常容易。
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.Configuration{ //默认配置:
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	}))
	// 或在run之前:
	// app.Configure(iris.WithConfiguration(iris.Configuration{...}))
}
```

### 通过toml加载配置 
#### 目录结构
> 主目录`fromTomlFile`
```html
    —— main.go
    —— configs
        —— iris.tml
```
#### 代码示例
> `main.go`
```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})
	// [...]
	//当你有两个配置时很好，一个用于开发，另一个用于生产用途。
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./configs/iris.tml")))
	// 会run之前加载:
	// app.Configure(iris.WithConfiguration(iris.TOML("./configs/iris.tml")))
	// app.Run(iris.Addr(":8080"))
}
```
> `/configs/iris.tml`
```tml
DisablePathCorrection = false
EnablePathEscape = false
FireMethodNotAllowed = true
DisableBodyConsumptionOnUnmarshal = false
TimeFormat = "Mon, 01 Jan 2006 15:04:05 GMT"
Charset = "UTF-8"

[Other]
	MyServerName = "iris"
```
### 通过`YAML`加载配置
#### 目录结构
> 主目录`fromYamlFile`
```html
    —— main.go
    —— configs
        —— iris.yml
```
#### 代码示例
> `main.go`
```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})
	// [...]
	//当你有两个配置时很好，一个用于开发，另一个用于生产用途。
	//如果iris.YAML的输入字符串参数为“〜”，则它从主目录加载配置
	//并且可以在许多iris实例之间共享。
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
	// 在run之前加载:
	// app.Configure(iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
	// app.Run(iris.Addr(":8080"))
}
```
> `/configs/iris.yml`

```yml
DisablePathCorrection: false
EnablePathEscape: false
FireMethodNotAllowed: true
DisableBodyConsumptionOnUnmarshal: true
TimeFormat: Mon, 01 Jan 2006 15:04:05 GMT
Charset: UTF-8
```
### 直接通过单个配置选项配置
#### 目录结构
> 主目录`functional`
```html
    —— main.go
```
#### 代码示例
> `main.go`
```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})
	// [...]
	//当您想要更改某些配置字段时，这很好。
	//前缀：With，代码编辑器将帮助您浏览所有内容
	//配置选项，甚至没有参考文档的类型。
	app.Run(iris.Addr(":8080"), iris.WithoutStartupLog, iris.WithCharset("UTF-8"))
	// 在run之前加载:
	// app.Configure(iris.WithoutStartupLog, iris.WithCharset("UTF-8"))
	// app.Run(iris.Addr(":8080"))
}
```
----------------------------------------------------------------------------------------------------------------------------------
## Built'n配置器
```go
// WithoutServerError将忽略错误，来自主应用程序的`Run`函数。
// 用法：
// err：= app.Run(iris.Addr（":8080"),iris.WithoutServerError(iris.ErrServerClosed))
// 如果服务器的错误是`http/iris＃ErrServerClosed`，//将返回`nil`。
// 也参见`Configuration＃IgnoreServerErrors [] string`。

// 示例: https://github.com/kataras/iris/tree/master/_examples/http-listening/listen-addr/omit-server-errors
func WithoutServerError(errors ...error) Configurator

//WithoutStartupLog 将关闭第一次运行终端就发送信息
var WithoutStartupLog

// WithoutInterruptHandler禁用自动正常服务器关闭，当按下control / cmd + C时
var WithoutInterruptHandler

// WithoutPathCorrection禁用PathCorrection设置
//参见`配置`
var WithoutPathCorrection

// WithoutBodyConsumptionOnUnmarshal禁用BodyConsumptionOnUnmarshal设置
//参见`配置`
var WithoutBodyConsumptionOnUnmarshal

// WithoutAutoFireStatusCode禁用AutoFireStatusCode设置
//参见`配置`
var WithoutAutoFireStatusCode

// WithPathEscape使PathEscape设置变为enanbles
//参见`配置`
var WithPathEscape

// WithOptimizations可以强制应用程序优化以获得尽可能最佳的性能
//参见`配置`
var WithOptimizations

// WithFireMethodNotAllowed使FireMethodNotAllowed设置变为可用
//参见`配置`
var WithFireMethodNotAllowed

//使用时间格式设置时间格式设置
//参见`配置`
func WithTimeFormat(timeformat string) Configurator

// WithCharset设置Charset设置
//参见`配置`
func WithCharset(charset string) Configurator

// WithRemoteAddrHeader启用或添加新的或现有的请求标头名称
// 可用于验证客户端的真实IP
//现有值为：
//"X-Real-Ip"：false，
//"X-Forwarded-For"：false，
//"CF-Connecting-IP"：false
//查看`context.RemoteAddr()`了解更多信息
func WithRemoteAddrHeader(headerName string) Configurator

// WithoutRemoteAddrHeader禁用现有的请求标头名称
//可用于验证客户端的真实IP。
//现有值为：
//"X-Real-Ip":false，
//"X-Forwarded-For":false，
//"CF-Connecting-IP":false
//查看`context.RemoteAddr()`了解更多信息
func WithoutRemoteAddrHeader(headerName string) Configurator

// WithOtherValue根据其他设置的键添加值
//参见`配置`
func WithOtherValue(key string, val interface{}) Configurator
```
## 自定义配置器
使用`Configurator`，开发人员可以轻松地模块化他们的应用程序，示例:
```go
// 文件名称 counter/counter.go
package counter

import (
    "time"
    "github.com/kataras/iris"
    "github.com/kataras/iris/core/host"
)

func Configurator(app *iris.Application) {
    counterValue := 0
    go func() {
        ticker := time.NewTicker(time.Second)
        for range ticker.C {
            counterValue++
        }
        app.ConfigureHost(func(h *host.Supervisor) { // < - 这里：重要的
            h.RegisterOnShutdown(func() {
                ticker.Stop()
            })
        })
    }()
    app.Get("/counter", func(ctx iris.Context) {
        ctx.Writef("Counter value = %d", counterValue)
    })
}
```

```go
// 文件名称: main.go
package main

import (
    "counter"
    "github.com/kataras/iris"
)

func main() {
    app := iris.New()
    app.Configure(counter.Configurator)
    app.Run(iris.Addr(":8080"))
}
```