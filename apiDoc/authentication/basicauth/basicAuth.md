# `HTTP BASIC`认证

### `BASIC`认证概述

在`HTTP`协议进行通信的过程中，`HTTP`协议定义了基本认证过程以允许`HTTP`服务器对`WEB`浏览器进行用户身份证的方法，当一个客户端向`HTTP`服务 器进行数据请求时，
如果客户端未被认证，则`HTTP`服务器将通过基本认证过程对客户端的用户名及密码进行验证，以决定用户是否合法。客户端在接收到`HTTP`服务器的身份认证要求后，
会提示用户输入用户名及密码，然后将用户名及密码以`BASE64`加密，加密后的密文将附加于请求信息中， 如当用户名为:`iris`，密码为:`123456`时，
客户端将用户名和密码用`:`合并，并将合并后的字符串用`BASE64`加密为密文，并于每次请求数据 时，将密文附加于请求头（`Request Header`）中。
`HTTP`服务器在每次收到请求包后，根据协议取得客户端附加的用户信息（`BASE64`加密的用户名和密码），解开请求包，对用户名及密码进行验证，
如果用 户名及密码正确，则根据客户端请求，返回客户端所需要的数据;否则，返回错误代码或重新要求客户端提供用户名及密码。

### `BASIC`认证的过程

- 客户端向服务器请求数据，请求的内容可能是一个网页或者是一个其它的`MIME`类型，此时，假设客户端尚未被验证，则客户端提供如下请求至服务器:

```go
    Get /index.html HTTP/1.0
    
    Host:www.studyiris.com
```
- 服务器向客户端发送验证请求代码401,服务器返回的数据:

```go
    HTTP/1.0 401 Unauthorised
    Server: nginx/1.0
    WWW-Authenticate: Basic realm="studyiris.com"
    Content-Type: text/html
    Content-Length: xxx
```

- 当符合`http1.0`或`1.1`规范的客户端浏览器收到`401`返回值时，将自动弹出一个登录窗口，要求用户输入用户名和密码
- 用户输入用户名和密码后，将用户名及密码以`BASE64`加密方式加密，并将密文放入前一条请求信息中，则客户端发送的第一条请求信息则变成如下内容:

```go
    Get /index.html HTTP/1.0
    Host:www.studyiris.com
    Authorization: Basic xxxxxxxxxxxxxxxxxxxxxxxxxxxx //加密串
```

- 服务器收到上述请求信息后，将`Authorization`字段后的用户信息取出、解密，将解密后的用户名及密码与用户数据库进行比较验证，
如用户名及密码正确，服务器则根据请求，将所请求资源发送给客户端

### `BASIC`认证缺点

`HTTP`基本认证的目标是提供简单的用户验证功能，其认证过程简单明了，适合于对安全性要求不高的系统或设备中，如大家所用路由器的配置页面的认证，
几乎 都采取了这种方式。其缺点是没有灵活可靠的认证策略，如无法提供域（`domain`或`realm`）认证功能，另外，`BASE64`的加密强度非常低。
当然，`HTTP`基本认证系统也可以其他加密技术一起，实现安全性能较高（相对）的认证系统

### `BASIC`iris 示例

```go
package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/basicauth"
)

func newApp() *iris.Application {
	app := iris.New()
	authConfig := basicauth.Config{
		Users:   map[string]string{"myusername": "mypassword", "mySecondusername": "mySecondpassword"},
		Realm:   "Authorization Required", // 默认表示域 "Authorization Required"
		Expires: time.Duration(30) * time.Minute,
	}
	authentication := basicauth.New(authConfig)
	//作用范围 全局 app.Use(authentication) 或者 (app.UseGlobal 在Run之前)
	//作用范围 单个路由 app.Get("/mysecret", authentication, h)
	app.Get("/", func(ctx iris.Context) { ctx.Redirect("/admin") })
	//作用范围  Party
	needAuth := app.Party("/admin", authentication)
	{
		//http://localhost:8080/admin
		needAuth.Get("/", h)
		// http://localhost:8080/admin/profile
		needAuth.Get("/profile", h)
		// http://localhost:8080/admin/settings
		needAuth.Get("/settings", h)
	}
	return app
}

func main() {
	app := newApp()
	// open http://localhost:8080/admin
	app.Run(iris.Addr(":8080"))
}

func h(ctx iris.Context) {
	username, password, _ := ctx.Request().BasicAuth()
	//第三个参数因为中间件所以不需要判断其值，否则不会执行此处理程序
	ctx.Writef("%s %s:%s", ctx.Path(), username, password)
}
```

### 提示

1. 运行上面的代码，访问`http://localhost:8080/admin`
2. 未验证时候会弹出一个验证框，让你输入用户名与密码，请认真看弹出框上面的内容

[Go Web Iris中文网](https://www.studyiris.com/)