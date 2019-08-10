###安装iris

####iris安装要求golang版本至少为1.8,建议1.9(本文档按照1.9进行编写)
  
> `$ go get -u github.com/kataras/iris`

>注解:Go 1.9支持类型别名，Iris已经为Go 1.8.3做好了准备,它有一个文件,kataras/iris/context.go 
对Go 1.9的限制构建访问权限声明了Iris的所有类型别名，所有这些都在同一个地方。

在 `Go 1.9 `之前，你必须导入“github.com/kataras/iris/context”来创建一个Handler：
如一下程序
```go
    package main
    import (
        "github.com/kataras/iris"
        "github.com/kataras/iris/context"//需要单独引入
    )
    func main() {
        app := iris.New()
        app.Get("/", func(ctx context.Context){})
        app.Run(iris.Addr(":8080"))
    }
```
从 Go 1.9开始，在您不必导入之后，您可以选择性地执行此操作：
```go
    package main
    import "github.com/kataras/iris"
    func main() {
        app := iris.New()
        app.Get("/", func(ctx iris.Context){})
        app.Run(iris.Addr(":8080"))
    }
```
同样的 kataras/iris/core/router/APIBuilder#PartyFunc
```go
    package main
    import (
        "github.com/kataras/iris"
        //"github.com/kataras/iris/core/router" 1.9之后需要引入
    )
    func main()  {
    app := iris.New()
    app.Get("/", func(ctx iris.Context){})
    app.Run(iris.Addr(":8080"))
    app.PartyFunc("/cpanel", func(child iris.Party) { 
         child.Get("/", func(ctx iris.Context){})
    })
     // OR
     cpanel := app.Party("/cpanel")
     cpanel.Get("/", func(ctx iris.Context){})
   }
```