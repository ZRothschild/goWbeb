package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	mvc.Configure(app.Party("/basic"), basicMVC)
	app.Run(iris.Addr(":8080"))
}

func basicMVC(app *mvc.Application) {
	//当然，你可以在MVC应用程序中使用普通的中间件。
	app.Router.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Path: %s", ctx.Path())
		ctx.Next()
	})
	//把依赖注入，controller(s)绑定
	//可以是一个接受iris.Context并返回单个值的函数（动态绑定）
	//或静态结构值（service）。
	app.Register(
		sessions.New(sessions.Config{}).Start,
		&prefixedLogger{prefix: "DEV"},
	)
	// GET: http://localhost:8080/basic
	// GET: http://localhost:8080/basic/custom
	app.Handle(new(basicController))
	//所有依赖项被绑定在父 *mvc.Application
	//被克隆到这个新子身上，父的也可以访问同一个会话。
	// GET: http://localhost:8080/basic/sub
	app.Party("/sub").Handle(new(basicSubController))
}

// If controller's fields (or even its functions) expecting an interface
// but a struct value is binded then it will check
// if that struct value implements
// the interface and if true then it will add this to the
// available bindings, as expected, before the server ran of course,
// remember? Iris always uses the best possible way to reduce load
// on serving web resources.
//如果控制器结构体的字段（甚至其方法）需要接口
//但结构值是绑定的，然后它会检查
//如果该结构值实现
//接口，如果为true，则将其添加到在服务器运行之前，正如预期的那样可用绑定，
//记得吗？ Iris总是使用最好的方法来减少负载关于提供网络资源。

type LoggerService interface {
	Log(string)
}

type prefixedLogger struct {
	prefix string
}

func (s *prefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.prefix, msg)
}

type basicController struct {
	Logger LoggerService
	Session *sessions.Session
}

func (c *basicController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/custom", "Custom")
}

func (c *basicController) AfterActivation(a mvc.AfterActivation) {
	if a.Singleton() {
		panic("basicController should be stateless,a request-scoped,we have a 'Session' which depends on the context.")
	}
}

func (c *basicController) Get() string {
	count := c.Session.Increment("count", 1)
	body := fmt.Sprintf("Hello from basicController\nTotal visits from you: %d", count)
	c.Logger.Log(body)
	return body
}

func (c *basicController) Custom() string {
	return "custom"
}

type basicSubController struct {
	Session *sessions.Session
}

func (c *basicSubController) Get() string {
	count := c.Session.GetIntDefault("count", 1)
	return fmt.Sprintf("Hello from basicSubController.\nRead-only visits count: %d", count)
}