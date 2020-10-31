package main

//在这个包中，我将向您展示如何覆盖现有Context的函数和方法。
//您可以轻松了看懂custom-context示例，以便于了解如何添加新功能到您自己的上下文context（需要自定义处理程序）。
//
//这种方法更容易理解，并且当您要覆盖现有方法时，它会更快：

// In this package I'll show you how to override the existing Context's functions and methods.
// You can easly navigate to the custom-context example to see how you can add new functions
// to your own context (need a custom handler).
//
// This way is far easier to understand and it's faster when you want to override existing methods:
import (
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

//创建您自己的自定义上下文，放入您需要的任何字段。

// Create your own custom Context, put any fields you wanna need.
type MyContext struct {
	//可选的第1部分：嵌入（可选，但如果您不想覆盖所有上下文方法，则为必需）

	// Optional Part 1: embed (optional but required if you don't want to override all context's methods)
	iris.Context
}

//（可选）：如果MyContext实现context.Context，则在编译时验证。
var _ iris.Context = &MyContext{} // optionally: validate on compile-time if MyContext implements context.Context.

//重要的是您将覆盖上下文
//具有嵌入式上下文。 为了使"*MyContext"通过路由处理函数，下面设置是必须的

// The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*MyContext".
func (ctx *MyContext) Do(handlers context.Handlers) {
	context.Do(ctx, handlers)
}

//如果您要覆盖上下文，则第二个重要的
//具有嵌入式上下文。 为了使"*MyContext"通过路由处理函数，下面设置是必须的

// The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*MyContext".
func (ctx *MyContext) Next() {
	context.Next(ctx)
}

//覆盖您想要的任何上下文方法...
//下面重写了 HTML方法

// Override any context's method you want...
// [...]

func (ctx *MyContext) HTML(format string, args ...interface{}) (int, error) {
	ctx.Application().Logger().Infof("Executing .HTML function from MyContext")

	ctx.ContentType("text/html")
	return ctx.Writef(format, args...)
}

func main() {
	app := iris.New()
	// app.Logger().SetLevel("debug")

	//唯一需要的一个：
	//这是您定义自己的上下文的方式
	//从Iris的通用上下文池创建和获取。

	// The only one Required:
	// here is how you define how your own context will
	// be created and acquired from the iris' generic context pool.
	app.ContextPool.Attach(func() iris.Context {
		return &MyContext{
			//可选的第3部分：

			// Optional Part 3:
			Context: context.NewContext(app),
		}
	})
	//在./view/**目录内的.html文件上注册视图引擎。

	// Register a view engine on .html files inside the ./view/** directory.
	app.RegisterView(iris.HTML("./view", ".html"))

	//像往常一样注册你的路由
	// register your route, as you normally do
	app.Handle("GET", "/", recordWhichContextJustForProofOfConcept, func(ctx iris.Context) {
		//下面使用了自定义上下文的的HTML方法。

		// use the context's overridden HTML method.
		ctx.HTML("<h1> Hello from my custom context's HTML! </h1>")
	})

	//这将由MyContext.Context执行
	//如果MyContext没有直接定义View函数则会执行Iris自己的View函数

	// this will be executed by the MyContext.Context
	// if MyContext is not directly define the View function by itself.
	app.Handle("GET", "/hi/{firstname:alphabetical}", recordWhichContextJustForProofOfConcept, func(ctx iris.Context) {
		firstname := ctx.Params().Get("firstname")
		ctx.ViewData("firstname", firstname)
		ctx.Gzip(true)

		ctx.View("hi.html")
	})
	app.Run(iris.Addr(":8080"))
}

//应该始终打印“($PATH)路由函数正在从'MyContext'执行” 也就是ctx iris.Context是 'MyContext'

// should always print "($PATH) Handler is executing from 'MyContext'"
func recordWhichContextJustForProofOfConcept(ctx iris.Context) {
	ctx.Application().Logger().Infof("(%s) Handler is executing from: '%s'", ctx.Path(), reflect.TypeOf(ctx).Elem().Name())
	ctx.Next()
}

//查看"new-implementation"，以了解如何创建具有新功能的全新Context。

// Look "new-implementation" to see how you can create an entirely new Context with new functions.
