# `IRIS session`使用
## 目录结构
> 主目录`sessions`

```html
    —— main.go
    —— routes
        —— index.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"time"
	"./routes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero" // <- 导入
	"github.com/kataras/iris/sessions"
)

func main() {
	app := iris.New()
	sessionManager := sessions.New(sessions.Config{
		Cookie:       "site_session_id",
		Expires:      60 * time.Minute,
		AllowReclaim: true,
	})
	//注册
	//动态依赖关系，比如* sessions.Session，来自`sessionManager.Start(ctx)*sessions.Session` < - 它接受一个Context并返回
	// something - >这称为动态请求 - 时间依赖，并且某些东西可以作为输入参数用于处理程序，
	//没有关于依赖项数量限制，每个处理程序将在服务器运行之前构建一次，并且它将仅使用它所需的依赖项。
	hero.Register(sessionManager.Start)
	//将任何函数转换为iris Handler，使用独特的Iris超快速依赖注入来解析它们的输入参数
	//用于服务或动态依赖，例如* sessions.Session，来自sessionManager.Start(ctx)* sessions.Session)< - 它接受一个Context并返回
	// 某些东西->这称为动态请求时依赖性。
	indexHandler := hero.Handler(routes.Index)
	// Method: GET
	// Path: http://localhost:8080
	app.Get("/", indexHandler)
	app.Run(
		iris.Addr(":8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}
```
> `/routes/index.go`

```go
package routes

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

// Index将根据此用户/session 所执行的访问来增加一个简单的int版本。
func Index(ctx iris.Context, session *sessions.Session) {
	//每一次访问自增一，如果不存在就先为你创建一个visits
	visits := session.Increment("visits", 1)
	//打印出当前的visits值
	ctx.Writef("%d visit(s) from my current session", visits)
}

/*
您还可以执行MVC功能可以执行的任何操作，即：
func Index(ctx iris.Context,session *sessions.Session) string {
	visits := session.Increment("visits", 1)
	return fmt.Spritnf("%d visit(s) from my current session", visits)
}
//你也可以省略iris.Context输入参数并使用LoginForm等依赖注入。< - 查看mvc示例。
*/
```