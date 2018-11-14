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