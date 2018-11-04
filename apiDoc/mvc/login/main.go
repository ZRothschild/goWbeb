//文件: main.go
package main

import (
	"time"
	"./datasource"
	"./repositories"
	"./services"
	"./web/controllers"
	"./web/middleware"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

func main() {
	app := iris.New()
	//你有完整的调试消息，在使用MVC时你很有用
	//确保您的代码与Iris的MVC架构保持一致。
	app.Logger().SetLevel("debug")
	//加载模板文件
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.StaticWeb("/public", "./web/public")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})
	// ----为我们的控制器服务----
	//准备我们的存储库和服务。
	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}
	repo := repositories.NewUserRepository(db)
	userService := services.NewUserService(repo)
	//“/user”基于mvc的应用程序。
	users := mvc.New(app.Party("/users"))
	//添加基本身份验证（admin：password）中间件
	//用于基于/ users的请求。
	users.Router.Use(middleware.BasicAuth)
	//将“userService”绑定到UserController的Service（interface）字段。
	users.Register(userService)
	users.Handle(new(controllers.UsersController))
	//“/user”基于mvc的应用程序。
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))
	// http://localhost:8080/noexist
	// and all controller's methods like
	// http://localhost:8080/users/1
	// http://localhost:8080/user/register
	// http://localhost:8080/user/login
	// http://localhost:8080/user/me
	// http://localhost:8080/user/logout
	//基本身份验证：“admin”，“password”，请参阅“./middleware/basicauth.go”源文件。
	app.Run(
		//在localhost：8080启动Web服务器
		iris.Addr("localhost:8080"),
		//按下CTRL/CMD+C时跳过错误的服务器：
		iris.WithoutServerError(iris.ErrServerClosed),
		//启用更快的json序列化和优化：
		iris.WithOptimizations,
	)
}