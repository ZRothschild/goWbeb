# `MVC` `hello world`简单用法
## 代码示例
> 文件名称 `main.go`
```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

//这个例子相当于
// /hello-world/main.go

//似乎是你的附加代码,必须写不值得但请记住，这个例子
//没有使用像iris mvc这样的功能
//Model，Persistence或View引擎,也不是Session，
//这对于学习目的来说非常简单,可能你永远不会用到这样的,作为应用程序中任何位置的简单控制器

//我们在这个示例中使用MVC
//在提供JSON的“/ hello”路径上
//在我的个人笔记本电脑上每20MB吞吐量大约2MB，
//它可以接受大多数应用程序，但你可以选择
//最适合你的是Iris，低级处理程序：性能
//或高级控制器：在大型应用程序上更易于维护和更小的代码库。

//当然你可以将所有这些都放到主函数中，它只是一个单独的函数
//用于main_test.go。
func newApp() *iris.Application {
	app := iris.New()
	//（可选）添加两个内置处理程序
	//可以从任何与http相关的panics中恢复
	//并将请求记录到终端。
	app.Use(recover.New())
	app.Use(logger.New())
	//控制器根路由路径"/"
	mvc.New(app).Handle(new(ExampleController))
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	// http://localhost:8080/custom_path
	app.Run(iris.Addr(":8080"))
}

// ExampleController提供 ”/”，“/ping”和 “/hello”路由选项
type ExampleController struct{}

// Get 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080
func (c *ExampleController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

// GetPing 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/ping
func (c *ExampleController) GetPing() string {
	return "pong"
}

// GetHello 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/hello
func (c *ExampleController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}

//在main函数调用controller之前调用一次BeforeActivation
//在版本9之后，您还可以为特定控制器的方法添加自定义路由
//在这里您可以注册自定义方法的处理程序
//使用带有`ca.Router`的标准路由器做一些你可以做的事情即使不是mvc
//并添加将绑定到控制器的字段或方法函数的输入参数的依赖项
func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("Inside /custom_path")
		ctx.Next()
	}
	b.Handle("GET", "/custom_path", "CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)
	//甚至添加基于此控制器路由的全局中间件，
	//在这个例子中是根“/”：
	// b.Router().Use(myMiddleware)
}

// CustomHandlerWithoutFollowingTheNamingGuide 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/custom_path
func (c *ExampleController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}

// GetUserBy 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/user/{username:string}
//是一个保留的关键字来告诉框架你要在函数的输入参数中绑定路径参数，
//在同一控制器中使用“Get”和“GetBy”可以实现
//
//func (c *ExampleController) GetUserBy(username string) mvc.Result {
// 	return mvc.View{
// 		Name: "user/username.html",
// 		Data: username,
// 	}
// }

/*
func (c *ExampleController) Post() {}
func (c *ExampleController) Put() {}
func (c *ExampleController) Delete() {}
func (c *ExampleController) Connect() {}
func (c *ExampleController) Head() {}
func (c *ExampleController) Patch() {}
func (c *ExampleController) Options() {}
func (c *ExampleController) Trace() {}
*/

/*
func (c *ExampleController) All() {}
//  或者
func (c *ExampleController) Any() {}

func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
	// 1 -> http 请求方法
	// 2 -> 请求路由
	// 3 -> 此控制器的方法名称应该是该路由的处理程序
	b.Handle("GET", "/mypath/{param}", "DoIt", optionalMiddlewareHere...)
}

//AfterActivation，所有依赖项都被设置,因此访问它们是只读
，但仍可以添加自定义控制器或简单的标准处理程序。
func (c *ExampleController) AfterActivation(a mvc.AfterActivation) {}
*/
```