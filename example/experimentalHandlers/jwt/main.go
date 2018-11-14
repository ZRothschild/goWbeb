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
	//如果解密成功，将会进入这里,获取解密了的token
	token := ctx.Values().Get("jwt").(*jwt.Token)
	//或者这样
	//userMsg :=ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	// userMsg["id"].(float64) == 1
	// userMsg["nick_name"].(string) == iris

	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")
	//可以了解一下token的数据结构
	ctx.Writef("%s", token.Signature)
}

func main() {
	app := iris.New()
	//jwt中间件
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		//这个方法将验证jwt的token
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//自己加密的秘钥或者说盐值
			return []byte("My Secret"), nil
		},
		//设置后，中间件会验证令牌是否使用特定的签名算法进行签名
		//如果签名方法不是常量，则可以使用ValidationKeyGetter回调来实现其他检查
		//重要的是要避免此处的安全问题：https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		//加密的方式
		SigningMethod: jwt.SigningMethodHS256,
		//验证未通过错误处理方式
		//ErrorHandler: func(context.Context, string)

		//debug 模式
		//Debug: bool
	})
	app.Use(jwtHandler.Serve)
	//解释：
	//jwtmiddleware.New是配置中间件的错误返回，是否为调试模式，机密秘钥，加密模式等
	//app.Use(jwtHandler.Serve) 是把中间件注册到处理程序中
	//注册次中间件的路由，中间件每一次都会去获取header头Authorization字段用户判断

	// 生成加密串过程
	//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//		"nick_name": "iris",
	//		"email":"go-iris@qq.com",
	//		"id":"1",
	//		"iss":"Iris",
	//		"iat":time.Now().Unix(),
	//		"jti":"9527",
	//		"exp":time.Now().Add(10*time.Hour * time.Duration(1)).Unix(),
	//	})
	//  把token已约定的加密方式和加密秘钥加密，当然也可以使用不对称加密
	//	tokenString, _ := token.SignedString([]byte("My Secret"))
	//  登录时候，把tokenString返回给客户端，然后需要登录的页面就在header上面附此字符串
	//  eg: header["Authorization"] = "bears "+tokenString

	app.Get("/ping", myHandler)
	app.Run(iris.Addr("localhost:3001"))
} //不要忘记查看../jwt_test.go.了解如何设置自己的自定义声明