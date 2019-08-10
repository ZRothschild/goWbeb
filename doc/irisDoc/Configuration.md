###配置
 在iris中 初始化应用程序 已经使用了默认的配置值,所以我们可以无需使用任何配置也可以启动我们的应
 用程序,像下面一样只需要简单的几行代码就可以运行我们web应用
 #### 通过程序内部配置

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
        //我们可以用这种方法单独定义我们的配置项
        app.Configure(iris.WithConfiguration(iris.Configuration{ DisableStartupLog:false}))
        //也可以使用app.run的第二个参数
        app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.Configuration{
            DisableInterruptHandler:           false,
            DisablePathCorrection:             false,
            EnablePathEscape:                  false,
            FireMethodNotAllowed:              false,
            DisableBodyConsumptionOnUnmarshal: false,
            DisableAutoFireStatusCode:         false,
            TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
            Charset:                           "UTF-8",
        }))
        //通过多参数配置 但是上面两种方式是我们最推荐的
        // 我们使用With+配置项名称 如WithCharset("UTF-8") 其中就是With+ Charset的组合
        //app.Run(iris.Addr(":8080"), iris.WithoutStartupLog, iris.WithCharset("UTF-8"))
        //当使用app.Configure(iris.WithoutStartupLog, iris.WithCharset("UTF-8"))设置配置项时
        //需要app.run()面前使用
    }
```   
####  通过TOML配置文件
我们在config 目录下新建main.tml
```go
    DisablePathCorrection = false
    EnablePathEscape = false
    FireMethodNotAllowed = true
    DisableBodyConsumptionOnUnmarshal = false
    TimeFormat = "Mon, 01 Jan 2006 15:04:05 GMT"
    Charset = "UTF-8"
    [Other]
        MyServerName = "iris"
```
在程序内做以下使用
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
        // 通过文件配置 我们可以更加方便的切换开发环境配置和生产环境.
        app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./configs/iris.tml")))
    }
```  
####通过YAML配置文件
```go
    DisablePathCorrection: false
    EnablePathEscape: false
    FireMethodNotAllowed: true
    DisableBodyConsumptionOnUnmarshal: true
    TimeFormat: Mon, 01 Jan 2006 15:04:05 GMT
    Charset: UTF-8
```
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
        app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
    }
```
####Built'n配置器
```go
// err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
// 当配置此项 如果web服务器 出现异常 我们将返回nil.
// 参考`Configuration的IgnoreServerErrors方法
// 地址: https://github.com/kataras/iris/tree/master/_examples/http-listening/listen-addr/omit-server-errors
func WithoutServerError(errors ...error) Configurator

// 当主服务器打开时，是否显示启动信息 如下
//Now listening on: http://localhost:8080
// Application started. Press CTRL+C to shut down.

var WithoutStartupLog

//当按下ctrl+C 时 禁止关闭当前程序(不会中止程序的运行)
var WithoutInterruptHandler

//路径重新定义(默认关闭)比如当访问/user/info 当该路径不存在的时候自动访问/user对应的handler
var WithoutPathCorrection

//如果此字段设置为true，则将创建一个新缓冲区以从请求主体读取。
var WithoutBodyConsumptionOnUnmarshal

//如果为true则关闭http错误状态代码处理程序自动执行
var WithoutAutoFireStatusCode

//转义路径
var WithPathEscape

//开启优化
var WithOptimizations

//不允许重新指向方法
var WithFireMethodNotAllowed

//设置时间格式
func WithTimeFormat(timeformat string) Configurator

//设值程序字符集
func WithCharset(charset string) Configurator

//启用或添加新的或现有的请求标头名称
func WithRemoteAddrHeader(headerName string) Configurator

//取消现有的请求标头名称
func WithoutRemoteAddrHeader(headerName string) Configurator

//自定义配置 key=>value
func WithOtherValue(key string, val interface{}) Configurator
```