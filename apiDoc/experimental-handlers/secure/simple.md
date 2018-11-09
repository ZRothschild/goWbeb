# 域名安全加密

## 示例代码 main.go

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/secure"
)

func main() {
	s := secure.New(secure.Options{
		// AllowedHosts是允许的完全限定域名列表。默认为空列表，允许任何和所有主机名。
		AllowedHosts:            []string{"ssl.example.com"},
		//如果SSLRedirect设置为true，则仅允许HTTPS请求。默认值为false。
		SSLRedirect:             true,
		//如果SSLTemporaryRedirect为true，则在重定向时将使用a 302。默认值为false（301）。
		SSLTemporaryRedirect:    false,
		// SSLHost是用于将HTTP请求重定向到HTTPS的主机名。默认值为“”，表示使用相同的主机。
		SSLHost:                 "ssl.example.com",
		// SSLProxyHeaders是一组标题键，其关联值表示有效的HTTPS请求。在使用Nginx时很有用：`map[string]string{"X-Forwarded-Proto”:“https"}`。默认为空白map。
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"},
		// STSSeconds是Strict-Transport-Security标头的max-age。默认值为0，不包括header。
		STSSeconds:              315360000,
		//如果STSIncludeSubdomains设置为true，则`includeSubdomains`将附加到Strict-Transport-Security标头。默认值为false。
		STSIncludeSubdomains:    true,
		//如果STSPreload设置为true，则`preload`标志将附加到Strict-Transport-Security标头。默认值为false。
		STSPreload:              true,
		//仅当连接是HTTPS时才包含STS标头。如果要强制始终添加，请设置为true."IsDevelopment"仍然覆盖了这一点。默认值为false。
		ForceSTSHeader:          false,
		//如果FrameDeny设置为true，则添加值为"DENY"的X-Frame-Options标头。默认值为false。
		FrameDeny:               true,
		// CustomFrameOptionsValue允许使用自定义值设置X-Frame-Options标头值。这会覆盖FrameDeny选项。
		CustomFrameOptionsValue: "SAMEORIGIN",
		//如果ContentTypeNosniff为true，则使用值nosniff添加X-Content-Type-Options标头。默认值为false。
		ContentTypeNosniff:      true,
		//如果BrowserXssFilter为true，则添加值为1的X-XSS-Protection标头;模式= block`。默认值为false。
		BrowserXSSFilter:        true,
		// ContentSecurityPolicy允许使用自定义值设置Content-Security-Policy标头值。默认为""。
		ContentSecurityPolicy:   "default-src 'self'",
		// PublicKey实现HPKP以防止伪造证书的MITM攻击。默认为""。
		PublicKey:               `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`,
		//这将导致在开发期间忽略AllowedHosts，SSLRedirect和STSSeconds/STSIncludeSubdomains选项。 部署到生产时，请务必将其设置为false。
		IsDevelopment: true,
	})
	app := iris.New()
	app.Use(s.Serve)
	app.Get("/home", func(ctx iris.Context) {
		ctx.Writef("Hello from /home")
	})
	app.Run(iris.Addr(":8080"))
}
```