# `IRIS`服务监听启用安全传输协议
## 目录结构
> 主目录`listenTls`

```html
    —— main.go
    —— mycert.cert
    —— mykey.key
```
## 代码示例
> `main.go`

```go
package main

import (
	"net/url"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Hello from the SECURE server")
	})
	app.Get("/mypath", func(ctx iris.Context) {
		ctx.Writef("Hello from the SECURE server on path /mypath")
	})
	//启动一个新的服务器,监听:80并重定向,到安全地址，然后:
	target, _ := url.Parse("https://127.0.1:443")
	go host.NewProxy("127.0.0.1:80", target).ListenAndServe()
	//在端口443上启动服务器(HTTPS)这是一个阻塞功能
	app.Run(iris.TLS("127.0.0.1:443", "mycert.cert", "mykey.key"))
}
```
> `mycert.cert`

```smartyconfig
-----BEGIN CERTIFICATE-----
MIIDAzCCAeugAwIBAgIJAPDsxtKV4v3uMA0GCSqGSIb3DQEBBQUAMBgxFjAUBgNV
BAMMDTEyNy4wLjAuMTo0NDMwHhcNMTYwNjI5MTMxMjU4WhcNMjYwNjI3MTMxMjU4
WjAYMRYwFAYDVQQDDA0xMjcuMC4wLjE6NDQzMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEA0KtAOHKrcbLwWJXgRX7XSFyu4HHHpSty4bliv8ET4sLJpbZH
XeVX05Foex7PnrurDP6e+0H5TgqqcpQM17/ZlFcyKrJcHSCgV0ZDB3Sb8RLQSLns
8a+MOSbn1WZ7TkC7d/cWlKmasQRHQ2V/cWlGooyKNEPoGaEz8MbY0wn2spyIJwsB
dciERC6317VTXbiZdoD8QbAsT+tBvEHM2m2A7B7PQmHNehtyFNbSV5uZNodvv1uv
ZTnDa6IqpjFLb1b2HNFgwmaVPmmkLuy1l9PN+o6/DUnXKKBrfPAx4JOlqTKEQpWs
pnfacTE3sWkkmOSSFltAXfkXIJFKdS/hy5J/KQIDAQABo1AwTjAdBgNVHQ4EFgQU
zr1df/c9+NyTpmyiQO8g3a8NswYwHwYDVR0jBBgwFoAUzr1df/c9+NyTpmyiQO8g
3a8NswYwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEACG5shtMSDgCd
MNjOF+YmD+PX3Wy9J9zehgaDJ1K1oDvBbQFl7EOJl8lRMWITSws22Wxwh8UXVibL
sscKBp14dR3e7DdbwCVIX/JyrJyOaCfy2nNBdf1B06jYFsIHvP3vtBAb9bPNOTBQ
QE0Ztu9kCqgsmu0//sHuBEeA3d3E7wvDhlqRSxTLcLtgC1NSgkFvBw0JvwgpkX6s
M5WpSBZwZv8qpplxhFfqNy8Uf+xrpSW0pGfkHumehkQGC6/Ry7raganS0aHhDPK9
Z1bEJ2com1bFFAQsm9yIXrRVMGGCtihB2Au0Q4jpEjUbzWYM+ItZyvRAGRM6Qex6
s/jogMeRsw==
-----END CERTIFICATE-----
```
> `mykey.key`

```smartyconfig
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA0KtAOHKrcbLwWJXgRX7XSFyu4HHHpSty4bliv8ET4sLJpbZH
XeVX05Foex7PnrurDP6e+0H5TgqqcpQM17/ZlFcyKrJcHSCgV0ZDB3Sb8RLQSLns
8a+MOSbn1WZ7TkC7d/cWlKmasQRHQ2V/cWlGooyKNEPoGaEz8MbY0wn2spyIJwsB
dciERC6317VTXbiZdoD8QbAsT+tBvEHM2m2A7B7PQmHNehtyFNbSV5uZNodvv1uv
ZTnDa6IqpjFLb1b2HNFgwmaVPmmkLuy1l9PN+o6/DUnXKKBrfPAx4JOlqTKEQpWs
pnfacTE3sWkkmOSSFltAXfkXIJFKdS/hy5J/KQIDAQABAoIBAQDCd+bo9I0s8Fun
4z3Y5oYSDTZ5O/CY0O5GyXPrSzCSM4Cj7EWEj1mTdb9Ohv9tam7WNHHLrcd+4NfK
4ok5hLVs1vqM6h6IksB7taKATz+Jo0PzkzrsXvMqzERhEBo4aoGMIv2rXIkrEdas
S+pCsp8+nAWtAeBMCn0Slu65d16vQxwgfod6YZfvMKbvfhOIOShl9ejQ+JxVZcMw
Ti8sgvYmFUrdrEH3nCgptARwbx4QwlHGaw/cLGHdepfFsVaNQsEzc7m61fSO70m4
NYJv48ZgjOooF5AccbEcQW9IxxikwNc+wpFYy5vDGzrBwS5zLZQFpoyMWFhtWdjx
hbmNn1jlAoGBAPs0ZjqsfDrH5ja4dQIdu5ErOccsmoHfMMBftMRqNG5neQXEmoLc
Uz8WeQ/QDf302aTua6E9iSjd7gglbFukVwMndQ1Q8Rwxz10jkXfiE32lFnqK0csx
ltruU6hOeSGSJhtGWBuNrT93G2lmy23fSG6BqOzdU4rn/2GPXy5zaxM/AoGBANSm
/E96RcBUiI6rDVqKhY+7M1yjLB41JrErL9a0Qfa6kYnaXMr84pOqVN11IjhNNTgl
g1lwxlpXZcZh7rYu9b7EEMdiWrJDQV7OxLDHopqUWkQ+3MHwqs6CxchyCq7kv9Df
IKqat7Me6Cyeo0MqcW+UMxlCRBxKQ9jqC7hDfZuXAoGBAJmyS8ImerP0TtS4M08i
JfsCOY21qqs/hbKOXCm42W+be56d1fOvHngBJf0YzRbO0sNo5Q14ew04DEWLsCq5
+EsDv0hwd7VKfJd+BakV99ruQTyk5wutwaEeJK1bph12MD6L4aiqHJAyLeFldZ45
+TUzu8mA+XaJz+U/NXtUPvU9AoGBALtl9M+tdy6I0Fa50ujJTe5eEGNAwK5WNKTI
5D2XWNIvk/Yh4shXlux+nI8UnHV1RMMX++qkAYi3oE71GsKeG55jdk3fFQInVsJQ
APGw3FDRD8M4ip62ki+u+tEr/tIlcAyHtWfjNKO7RuubWVDlZFXqCiXmSdOMdsH/
bxiREW49AoGACWev/eOzBoQJCRN6EvU2OV0s3b6f1QsPvcaH0zc6bgbBFOGmJU8v
pXhD88tsu9exptLkGVoYZjR0n0QT/2Kkyu93jVDW/80P7VCz8DKYyAJDa4CVwZxO
MlobQSunSDKx/CCJhWkbytCyh1bngAtwSAYLXavYIlJbAzx6FvtAIw4=
-----END RSA PRIVATE KEY-----
```