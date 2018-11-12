# 小项目示例
## 目录结构
> 主目录`overview`

```html
—— views
    —— user
        —— create_verification.html
        —— profile.html
—— main.go
```
## 示例代码
> `main.go`

```go
package main

import (
	"github.com/kataras/iris"
)
//User结构体用户绑定提交参数，与json返回
type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

func main() {
	app := iris.New()
	// app.Logger().SetLevel("disable")禁用错误日志记录
	//使用std html/template引擎定义模板
	//使用“.html”文件扩展名解析并加载“./views”文件夹中的所有文件。
	//在每个请求上重新加载模板（开发模式）。
	app.RegisterView(iris.HTML("./views", ".html").Reload(true))
	//为特定的http错误注册自定义处理程序。
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		// .Values用于在处理程序，中间件之间进行通信。
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) internal server error")
	})
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	// app.Done(func(ctx iris.Context) {]})
	// POST: scheme://mysubdomain.$domain.com/decode
	app.Subdomain("mysubdomain.").Post("/decode", func(ctx iris.Context) {})
	// 请求方法 POST: http://localhost:8080/decode
	app.Post("/decode", func(ctx iris.Context) {
		var user User
		ctx.ReadJSON(&user)
		ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
	})
	// 请求方法 GET: http://localhost:8080/encode
	app.Get("/encode", func(ctx iris.Context) {
		doe := User{
			Username:  "Johndoe",
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
			Age:       25,
		}
		ctx.JSON(doe)
	})
	// 请求方法 GET: http://localhost:8080/profile/anytypeofstring
	app.Get("/profile/{username:string}", profileByUsername)
	usersRoutes := app.Party("/users", logThisMiddleware)
	{
		// 请求方法 GET: http://localhost:8080/users/42
		usersRoutes.Get("/{id:int min(1)}", getUserByID)
		// 请求方法 POST: http://localhost:8080/users/create
		usersRoutes.Post("/create", createUser)
	}
	//在localhost端口8080上监听传入的HTTP/1.x和HTTP/2客户端
	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))
}

func logThisMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())
	//.Next继续运行下一个链上的处理程序（中间件）
	//如果没有，程序就会在这里停止了，不会向下面继续执行
	ctx.Next()
}

func profileByUsername(ctx iris.Context) {
	// .Params用于获取动态路径参数
	username := ctx.Params().Get("username")
	ctx.ViewData("Username", username)
	//渲染“./views/user/profile.html”
	//使用 {{ .Username }}将动态路径参数，渲染到html页面
	ctx.View("user/profile.html")
}

func getUserByID(ctx iris.Context) {
	userID := ctx.Params().Get("id") //或直接转换为：.Values().GetInt/GetInt64等...
	user := User{Username: "username" + userID}
	ctx.XML(user)
}

func createUser(ctx iris.Context) {
	var user User
	err := ctx.ReadForm(&user)
	if err != nil {
		ctx.Values().Set("error", "creating user, read and parse form failed. "+err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	//渲染“./views/user/create_verification.html”
	//{{ . }}相当于User结构体，例如{{ .Username }}表示User结构体的Username字段等等...
	ctx.ViewData("", user)
	ctx.View("user/create_verification.html")
}
```
> `/views/user/create_verification.html`

```html
<html>
    <head><title>Create verification</title></head>
    <body>
        <h1> Create Verification </h1>
        <table style="width:550px">
        <tr>
            <th>Username</th>
            <th>Firstname</th>
            <th>Lastname</th>
            <th>City</th>
            <th>Age</th>
        </tr>
        <tr>
            <td>{{ .Username }}</td>
            <td>{{ .Firstname }}</td>
            <td>{{ .Lastname }}</td>
            <td>{{ .City }}</td>
            <td>{{ .Age }}</td>
        </tr>
        </table> 
    </body>
</html>
```
> `/views/user/profile.html`

```html
<html>
    <head><title>Profile page</title></head>
    <body>
        <h1> Profile </h1>
        <b> {{ .Username }} </b>
    </body>
</html>
```