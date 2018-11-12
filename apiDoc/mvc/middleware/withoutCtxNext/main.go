/*
Package main是处理程序执行流程的行为更改的简单示例，
通常我们需要`ctx.Next()`来调用,路由处理程序链中的下一个处理程序，
但是使用新的`ExecutionRules`，我们可以更改此默认行为。
请继续阅读以下内容。

`Party＃SetExecutionRules`改变处理程序本身之外的路由处理程序的执行流程

For example, if for some reason the desired result is the (done or all) handlers to be executed no matter what
even if no `ctx.Next()` is called in the previous handlers, including the begin(`Use`),
the main(`Handle`) and the done(`Done`) handlers themselves, then:

例如，如果由于某种原因，期望的结果是（完成或所有）处理程序无论如何都要执行
即使在前面的处理程序中没有调用`ctx.Next()`，包括begin(`Use`)，
main（`Handle`）和done（`Done`）处理程序本身，然后：

Party#SetExecutionRules(iris.ExecutionRules {
  Begin: iris.ExecutionOptions{Force: true},
  Main:  iris.ExecutionOptions{Force: true},
  Done:  iris.ExecutionOptions{Force: true},
})

注意，如果`true`那么“break”处理程序链的唯一剩余方法是`ctx.StopExecution()`现在`ctx.Next（）`无关紧要。

这些规则是按先后进行的，因此如果`Party`创建了一个子规则，那么同样的规则也将应用于该规则。
可以使用`Party #SetExecutionRules（iris.ExecutionRules {}）`来重置这些规则（在`Party＃Handle`之前）。

最常见的使用方案可以在Iris MVC应用程序中找到;
当我们想要特定mvc应用程序的`Party`的`Done`处理程序时
要执行，但我们不想在`exampleController #EndRequest`上添加`ctx.Next()` */

package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) { ctx.Redirect("/example") })
	// example := app.Party("/example")
	// example.SetExecutionRules && mvc.New(example) 或...
	m := mvc.New(app.Party("/example"))
	//重要
	//所有选项都可以用Force填充：true，所有的都会很好的兼容
	m.Router.SetExecutionRules(iris.ExecutionRules{
		//Begin: <- from `Use[all]` 到`Handle[last]` 程序执行顺序，执行all，即使缺少`ctx.Next()`也执行all。
		// Main: <- all `Handle` 执行顺序，执行所有>> >>。
		Done: iris.ExecutionOptions {Force:true},// < - 从`Handle [last]`到`Done [all]`程序执行顺序，执行全部>> >>。
		})
	m.Router.Done(doneHandler)
	// m.Router.Done(...)
	// ...
	m.Handle(&exampleController{})
	app.Run(iris.Addr(":8080"))
}

func doneHandler(ctx iris.Context) {
	ctx.WriteString("\nFrom Done Handler")
}

type exampleController struct{}

func (c *exampleController) Get() string {
	return "From Main Handler"
	//注意，这里我们不绑定`Context`，我们不调用它的`Next()`
	//函数以调用`doneHandler`，
	//这是我们自动完成的，因为我们用`SetExecutionRules`改变了执行规则。
	//因此最终输出是：
	//来自Main Handler
	//来自Done Handler
}