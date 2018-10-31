// iris提供了一些基本的中间件，大部分用于学习曲线。
//您可以将任何net/http请求的中间件与iris.FromStd（把net/http转化成iris context形式）包装器一起使用
//适用于Golang新手的JWT net/http视频教程：https://www.youtube.com/watchv=dgJFeqeXVKw

//这个中间件是唯一一个从外部源克隆的中间件：https://github.com/auth0/go-jwt-middleware
//（因为它使用“context”来定义用户，但我们不需要这样，因此简单的iris.FromStd将无法按预期工作。）
package main

// $ go get -u github.com/dgrijalva/jwt-go
// $ go run main.go

import (
	"github.com/kataras/iris"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

func myHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")
	ctx.Writef("%s", user.Signature)
}

func main() {
	app := iris.New()
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		//设置后，中间件会验证令牌是否使用特定的签名算法进行签名
		//如果签名方法不是常量，则可以使用ValidationKeyGetter回调来实现其他检查
		//重要的是要避免此处的安全问题：https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
	app.Use(jwtHandler.Serve)
	app.Get("/ping", myHandler)
	app.Run(iris.Addr("localhost:3001"))
} //不要忘记查看../jwt_test.go.了解如何设置自己的自定义声明
