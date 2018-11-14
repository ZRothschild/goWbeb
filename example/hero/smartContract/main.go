package main

import (
	"github.com/kataras/iris"
	"fmt"
	"github.com/jmespath/go-jmespath"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/hero"
	"strings"
)
/*
$ go get github.com/jmespath/go-jmespath
*/
func newApp() *iris.Application {
	app := iris.New()
	// PartyFunc 等同于 usersRouter := app.Party("/users")
	//但它为我们提供了一种简单的方法来调用路由组的注册路由方法，
	//即来自另一个可以处理这组路由的包的函数。
	app.PartyFunc("/users", registerUsersRoutes)
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080/users?query=[?Name == 'John Doe'].Age
	// < - 客户端将收到用户的年龄，他的名字是John Doe
	//您还可以测试query = [0] .Name以检索第一个用户的名称。
	//甚至query= [0:3].Age打印前三个年龄段。
	//了解有关jmespath以及如何过滤的更多信息：
	// http://jmespath.readthedocs.io/en/latest/ and
	// https://github.com/jmespath/go-jmespath/tree/master/fuzz/testdata
	// http://localhost:8080/users
	// http://localhost:8080/users/William%20Woe
	// http://localhost:8080/users/William%20Woe/age
	app.Run(iris.Addr(":8080"))
}

/*
开始使用路由
*/
func registerUsersRoutes(usersRouter iris.Party) {
	// GET: /users
	usersRouter.Get("/", getAllUsersHandler)
	usersRouter.PartyFunc("/{name:string}", registerUserRoutes)
}

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var usersSample = []*user{
	{"William Woe", 25},
	{"Mary Moe", 15},
	{"John Doe", 17},
}

func getAllUsersHandler(ctx iris.Context) {
	err := sendJSON(ctx, usersSample)
	if err != nil {
		fail(ctx, iris.StatusInternalServerError, "unable to send a list of all users: %v", err)
		return
	}
}

//开始使用  USERS路由组的子路由组
func registerUserRoutes(userRouter iris.Party) {
	//为此子路由器创建一个新的依赖注入管理器
	userDeps := hero.New()
	//你也可以使用global/package-level的hero.Register(userDependency)，正如我们在其他例子中已经学到的那样。
	userDeps.Register(userDependency)
	// GET: /users/{name:string}
	userRouter.Get("/", userDeps.Handler(getUserHandler))
	// GET: /users/{name:string}/age
	userRouter.Get("/age", userDeps.Handler(getUserAgeHandler))
}

var userDependency = func(ctx iris.Context) *user {
	name := strings.Title(ctx.Params().Get("name"))
	for _, u := range usersSample {
		if u.Name == name {
			return u
		}
	}
	// you may want or no to handle the error here, either way the main route handler
	// is going to be executed, always. A dynamic dependency(per-request) is not a middleware, so things like `ctx.Next()` or `ctx.StopExecution()`
	// do not apply here, look the `getUserHandler`'s first lines; we stop/exit the handler manually
	// if the received user is nil but depending on your app's needs, it is possible to do other things too.
	// A dynamic dependency like this can return more output values, i.e (*user, bool).
	//你可能想要或不想在这里处理错误，无论是主路由处理程序
	//将永远执行。 动态依赖（每个请求）不是中间件，所以像`ctx.Next()`或`ctx.StopExecution()`
	//不要在这里申请，看看`getUserHandler`的第一行; 我们手动停止/退出处理程序
	//如果收到的用户是零，但根据您的应用程序的需要，也可以做其他事情。
	//像这样的动态依赖可以返回更多的输出值，即(*user，bool)。
	fail(ctx, iris.StatusNotFound, "user with name '%s' not found", name)
	return nil
}

func getUserHandler(ctx iris.Context, u *user) {
	if u == nil {
		return
	}
	sendJSON(ctx, u)
}

func getUserAgeHandler(ctx iris.Context, u *user) {
	if u == nil {
		return
	}
	ctx.Writef("%d", u.Age)
}

/*请记住，使用'hero'可以获得类似mvc的函数，所以这也可以工作：
func getUserAgeHandler(u *user) string {
	if u == nil {
		return ""
	}
	return fmt.Sprintf("%d", u.Age)
}
*/

/* ENDERS USERS.USER SUB ROUTER */
/* 用户路由器结束 */
//可选择手动HTTP错误的常见JSON响应。
type httpError struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func (h httpError) Error() string {
	return fmt.Sprintf("Status Code: %d\nReason: %s", h.Code, h.Reason)
}

func fail(ctx context.Context, statusCode int, format string, a ...interface{}) {
	err := httpError{
		Code:   statusCode,
		Reason: fmt.Sprintf(format, a...),
	}
	//记录所有> = 500内部错误。
	if statusCode >= 500 {
		ctx.Application().Logger().Error(err)
	}
	ctx.StatusCode(statusCode)
	ctx.JSON(err)
	//没有下一个处理程序将运行。
	ctx.StopExecution()
}

// JSON辅助函数，为最终用户提供插入字符或过滤响应的能力，您可以选择这样做。
//如果您想在Iris的Context中看到该函数，则会引发[Feature Request]问题并链接此示例。
func sendJSON(ctx iris.Context, resp interface{}) (err error) {
	indent := ctx.URLParamDefault("indent", "  ")
	// i.e [?Name == 'John Doe'].Age # to output the [age] of a user which his name is "John Doe".
	if query := ctx.URLParam("query"); query != "" && query != "[]" {
		resp, err = jmespath.Search(query, resp)
		if err != nil {
			return
		}
	}
	_, err = ctx.JSON(resp, context.JSON{Indent: indent, UnescapeHTML: true})
	return err
}