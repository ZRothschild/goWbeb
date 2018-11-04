// 文件: middleware/basicauth.go

package middleware

import "github.com/kataras/iris/middleware/basicauth"

// BasicAuth中间件示例。
var BasicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{
		"admin": "password",
	},
})