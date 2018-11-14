# `iris`前置自定义读取请求(不需要自己实现实现`Decode`接口)
## 目录结构
> 主目录`readCustomViaUnmarshaler`

```html
    —— main.go
    —— main_test.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"gopkg.in/yaml.v2"
	"github.com/kataras/iris"
)

func main() {
	app := newApp()
	//使用Postman或其他什么来做POST请求
	//（但是你总是可以自由地使用app.Get和GET http方法请求来读取请求值）
	//使用RAW BODY到http//localhost:8080：
	/*
		addr: localhost:8080
		serverName: Iris
	*/
	//响应应该是：
	//收到: main.config{Addr:"localhost:8080", ServerName:"Iris"}
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}

func newApp() *iris.Application {
	app := iris.New()
	app.Post("/", handler)
	return app
}

//简单的yaml内容，请访问https://github.com/go-yaml/yaml阅读更多内容
type config struct {
	Addr       string `yaml:"addr"`
	ServerName string `yaml:"serverName"`
}

/*
type myBodyDecoder struct{}

var DefaultBodyDecoder = myBodyDecoder{}

//在我们的例子中实现`kataras/iris/context＃Unmarshaler`
//我们将使用最简单的`context＃UnmarshalerFunc`来传递yaml.Unmarshal。
//可以用作：ctx.UnmarshalBody(＆c.DefaultBodyDecoder)

func (r *myBodyDecoder) Unmarshal(data []byte, outPtr interface{}) error {
	return yaml.Unmarshal(data, outPtr)
}
*/

func handler(ctx iris.Context) {
	var c config
	// 注意:yaml.Unmarshal已经实现了`context＃Unmarshaler`
	//所以我们可以直接使用它，比如json.Unmarshal(ctx.ReadJSON)，xml.Unmarshal(ctx.ReadXML)
	//以及遵循最佳实践并符合Go标准的每个库。
	//注意2：如果你因任何原因需要多次读取访问body
	//你应该通过`app.Run（...，iris.WithoutBodyConsumptionOnUnmarshal）消费body
	if err := ctx.UnmarshalBody(&c, iris.UnmarshalerFunc(yaml.Unmarshal)); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Writef("Received: %#+v", c)
}
```
> `main_test.go`

```go
package main

import (
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestReadCustomViaUnmarshaler(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	expectedResponse := `Received: main.config{Addr:"localhost:8080", ServerName:"Iris"}`
	e.POST("/").WithText("addr: localhost:8080\nserverName: Iris").Expect().
		Status(httptest.StatusOK).Body().Equal(expectedResponse)
}
```