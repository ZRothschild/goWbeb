// file: main.go
package main

import (
	"./datasource"
	"./repositories"
	"./services"
	"./web/middleware"
	"./web/routes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	//加载模板文件
	app.RegisterView(iris.HTML("./web/views", ".html"))
	//获取仓库所电影资源
	repo := repositories.NewMovieRepository(datasource.Movies)
	//创建我们的电影服务，我们将它绑定到电影应用程序的依赖项。
	movieService := services.NewMovieService(repo)
	hero.Register(movieService)
	//使用hero处理程序对应我们的路由
	app.PartyFunc("/hello", func(r iris.Party) {
		r.Get("/", hero.Handler(routes.Hello))
		r.Get("/{name}", hero.Handler(routes.HelloName))
	})
	app.PartyFunc("/movies", func(r iris.Party) {
		//添加基本basic authentication（admin：password）中间件用户/movies 请求
		r.Use(middleware.BasicAuth)
		r.Get("/", hero.Handler(routes.Movies))
		r.Get("/{id:long}", hero.Handler(routes.MovieByID))
		r.Put("/{id:long}", hero.Handler(routes.UpdateMovieByID))
		r.Delete("/{id:long}", hero.Handler(routes.DeleteMovieByID))
	})
	// http://localhost:8080/hello
	// http://localhost:8080/hello/iris
	// http://localhost:8080/movies
	// http://localhost:8080/movies/1
	app.Run(
		//在localhost：8080启动Web服务器
		iris.Addr("localhost:8080"),
		//按下CTRL/CMD+C时跳过错误的服务器：
		iris.WithoutServerError(iris.ErrServerClosed),
		//启用更快的json序列化和优化：
		iris.WithOptimizations,
	)
}