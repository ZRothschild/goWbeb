# `JWT`
## `JWT`介绍
`JSON Web Token(JWT)`是一个开放标准`(RFC 7519)`，它定义了一种紧凑且独立的方式，可以在各方之间作为`JSON`对象安全地传输信息。
此信息可以通过数字签名进行验证和信任。`JWT`可以使用秘密（使用`HMAC`算法）或使用`RSA`或`ECDSA`的公钥/私钥对进行签名。
## 什么场景适合用`JWT`
1. **授权**：这是使用`JWT`的最常见方案。一旦用户登录，每个后续请求将包括`JWT`，允许用户访问该令牌允许的路由，服务和资源。
`Single Sign On(单点登录)`是一种现在广泛使用`JWT`的功能，因为它的开销很小，并且能够在不同的域中轻松使用。
2. **信息交换**：`JWT`是在各方之间安全传输信息的好方法。 因为`JWT`可以签名 - 例如，使用公钥/私钥对 - 您可以确定发件人是他们所说的人。
此外，由于使用`Header`和`payload`,`Signature`，您还可以验证内容是否未被篡改。
##  `JWT`结构
> 在紧凑的形式中，`JWT`由三个部分组成，用点（.）分隔，
1. Header
2. Payload
3. Signature
> 因此，JWT通常如下所(token)
```html
xxxxx.yyyyy.zzzzz (Header.Payload.Signature)
````
### `Header`
`Header`通常由两部分组成：令牌的类型，即`JWT`，以及正在使用的散列算法，例如`HMAC SHA256`或`RSA`。示例如下：
```html
{
  "alg": "HS256",
  "typ": "JWT"
}
````
> 然后，这个JSON被编码为Base64Url，形成JWT的第一部分。
### `Payload`
`token`的第二部分是`payload`，其中包含`claims`。`claims`是关于实体（通常是用户）和其他数据的声明。`claims`有三种类型：
`registered`, `public`, `private claims`
1. `Registered claims`：这些是一组预定义声明，不是强制性的，但建议使用，以提供一组有用的，可互操作的声明。 其中一些是：
`iss`(发行人)，`exp`(过期时间)，`sub`(主题)，`aud`(观众)等。
2. `Public claims`：这些可以由使用`JWT`的人随意定义。但是为避免冲突，应在`IANA JSON Web`令牌注册表中定义它们，或者将其定义为
包含防冲突命名空间的`URI`。
3. `Private claims`：这些声明是为了在同意使用它们的各方之间共享信息而创建的，并且既不是注册声明也不是公开声明。
> `Payload`示例
```html
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
````
> 然后，`Payload`经过`Base64Url`编码，形成`JSON Web Token`的第二部分，数据虽然是不可串改，但是确实透明的
### `Signature`
要创建签名部分，您必须采用`base64Url`编码`header`，`base64Url`编码的`payload`，`secret`，`header`中指定的算法，并对其进行签名

例如，如果要使用`HMAC SHA256`算法，将按以下方式创建签名：
```html
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
````
签名用于验证消息在此过程中未被更改，并且，在使用私钥签名的令牌的情况下，它还可以验证`JWT`的发件人是否是它所声称的人。
##  合并`jwt`三部分
输出是三个由点分隔的`Base64-URL`字符串，可以在`HTML`和`HTTP`环境中轻松传递，与`SAML`等基于`XML`的标准相比更加紧凑。

下面显示了一个`JWT`，它具有`Header`和`Payload`，并使`Signature`
```html
//encoded
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6
IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

//decoded

//Header 部分
{
  "alg": "HS256",
  "typ": "JWT"
}

//Payload 部分

{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}

//Signature 部分
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  your-256-bit-secret
)
```
## `JWT`工作原理
在身份验证中，当用户使用其凭据成功登录时，将返回`JSON Web Token`。由于`Token`是凭证，因此必须非常小心以防止出现安全问题。
一般情况下，您不应该将令牌保留的时间超过要求。

每当用户想要访问受保护的路由或资源时，用户代理应该使用承载模式发送`JWT`，通常在`Authorization Header`中。`Header`的内容应如下所示：
```html
Authorization: Bearer <token>
```
在某些情况下，这可以是无状态授权机制。服务器的受保护路由将在`Authorization Header`中检查有效的`JWT`，如果存在，则允许用户访问
受保护的资源。如果`JWT`包含必要的数据，则可以减少查询数据库以进行某些操作的需要，尽管可能并非总是如此。

如果在`Authorization Header`中发送`Token`，则跨域资源共享(`CORS`)将不会成为问题，因为它不使用`cookie`。

下图显示了如何获取`JWT`并用于访问`API`或资源：
![JWT工作流程图](https://cdn2.auth0.com/docs/media/articles/api-auth/client-credentials-grant.png)
1. 应用程序或客户端向授权服务器请求授权。这是通过其中一个不同的授权流程执行的。例如，典型的`OpenID Connect`兼容`Web`应用程序
将使用授权代码流通过`/oauth/authorize`端点。
2. 授予授权后，授权服务器会向应用程序返回访问`Token`。
3. 应用程序使用访问`Token`来访问受保护资源（如`API`）。
> 请注意，使用签名`Token`，`Token`中包含的所有信息都会向用户或其他方公开，即使他们无法更改。这意味着您不应该在`Token`中放置秘密信息。
## 目录结构
> 主目录`jwt`
```html
    —— main.go
```
## 代码示例 
> `main.go`
```go
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
}
```