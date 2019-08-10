# 请求事务(`TransactionScope`)
## 目录结构
> 主目录`transactions`

```html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func main() {
	app := iris.New()
	//子域可以与所有可用的路由器一起使用，就像其他功能一样。
	app.Get("/", func(ctx context.Context) {
		ctx.BeginTransaction(func(t *context.Transaction) {
			// 选择步骤：如果为true，那么下一个转录将不会被执行，如果为fails则向反
			// t.SetScope（context.RequestTransactionScope）
			//可选步骤：
			//在此处创建新的自定义错误类型，以跟踪状态代码和错误消息
			err := context.NewTransactionErrResult()
			//如果我们想要回滚这个函数clojure中的任何错误，我们应该使用t.Context。
			t.Context().Text("Blablabla this should not be sent to the client because we will fill the err with a message and status")
			//在这里虚拟化一个虚假的错误，为了测试这个例子
			fail := true
			if fail {
				err.StatusCode = iris.StatusInternalServerError
				//注意：如果为空原因则会触发默认或自定义http错误（如ctx.FireStatusCode）
				err.Reason = "Error: Virtual failure!!"
                //选择步骤：
                //但是如果我们想要在事务失败时将错误消息回发给客户端，这很有用
                //如果原因为空，那么交易成功完成，
                //否则我们回滚整个返回的响应体，
                //header和Cookie，状态代码以及此事务中的所有内容
                t.Complete(err)
			}
		})
		ctx.BeginTransaction(func(t *context.Transaction) {
			t.Context().HTML("<h1>This will sent at all cases because it lives on different transaction and it doesn't fails</h1>")
			// *如果我们没有任何'throw error'逻辑，则不需要scope.Complete（）
		})
		// OPTIONALLY，取决于用法：
		//无论如何，在上下文的事务中发生的情景都会发送给客户端
		ctx.HTML("<h1>Let's add a second html message to the response, " +
			"if the transaction was failed and it was request scoped then this message would " +
			"not been shown. But it has a transient scope(default) so, it is visible as expected!</h1>")
	})
	app.Run(iris.Addr(":8080"))
}
```
## 提示
1. 尝试修改一下`fail`变量值，当与到错误，则会不显示`BeginTransaction`里面的内容，也就是说回滚了