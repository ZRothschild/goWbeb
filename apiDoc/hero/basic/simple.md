# Hero Basic

### 示例 main.go

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func main() {
	app := iris.New()
	// 1.直接把hello函数转化成iris请求处理函数
	helloHandler := hero.Handler(hello)
	app.Get("/{to:string}", helloHandler)
	// 2.把结构体实例注入hero,在把在结构体方法转化成iris请求处理函数
	hero.Register(&myTestService{
		prefix: "Service: Hello",
	})
	helloServiceHandler := hero.Handler(helloService)
	app.Get("/service/{to:string}", helloServiceHandler)
	// 3.注册一个iris请求处理函数，是以from表单格式x-www-form-urlencoded数据类型,以LoginForm类型映射
	// 然后把login方法转化成iris请求处理函数
	hero.Register(func(ctx iris.Context) (form LoginForm) {
		//绑定from方式提交以x-www-form-urlencoded数据格式传输的from数据，并返回相应结构体
		ctx.ReadForm(&form)
		return
	})
	loginHandler := hero.Handler(login)
	app.Post("/login", loginHandler)
	// http://localhost:8080/your_name
	// http://localhost:8080/service/your_name
	app.Run(iris.Addr(":8080"))
}

func hello(to string) string {
	return "Hello " + to
}

type Service interface {
	SayHello(to string) string
}

type myTestService struct {
	prefix string
}

func (s *myTestService) SayHello(to string) string {
	return s.prefix + " " + to
}

func helloService(to string, service Service) string {
	return service.SayHello(to)
}

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func login(form LoginForm) string {
	return "Hello " + form.Username
}
``` 

### 提示

Handler接口有一个handler函方法，它可以接受任何匹配的输入参数
使用Hero的"Dependencies"和任何输出结果; 像string，int（string，int），
自定义结构，结果（查看|响应）以及您可以想象的任何内容。
它返回一个标准的`iris/context.Handler`，它可以在Iris应用程序的任何地方使用，
作为中间件或简单路由处理程序或子域的处理程序