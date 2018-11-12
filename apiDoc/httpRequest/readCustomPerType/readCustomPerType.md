# `iris`前置自定义读取请求数据数据处理
## 目录结构
> 主目录`readCustomPerType`
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
//Decode实现`kataras/iris/context＃BodyDecoder`可选接口
//任何go类型都可以实现，以便在读取请求的主体时进行自解码。
func (c *config) Decode(body []byte) error {
	return yaml.Unmarshal(body, c)
}

func handler(ctx iris.Context) {
	var c config
	// 注意：第二个参数是nil，因为我们的＆c实现了`context＃BodyDecoder`
	//优先于上下文#Unmarshaler（可以是读取请求主体全局选项）
	//参见`http_request/read-custom-via-unmarshaler/main.go`示例，了解如何使用上下文＃Unmarshaler。
	// 注意2：如果你因任何原因需要多次读取访问body
	//你应该通过`app.Run（...，iris.WithoutBodyConsumptionOnUnmarshal）消费body
	if err := ctx.UnmarshalBody(&c, nil); err != nil {
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

func TestReadCustomPerType(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	expectedResponse := `Received: main.config{Addr:"localhost:8080", ServerName:"Iris"}`
	e.POST("/").WithText("addr: localhost:8080\nserverName: Iris").Expect().
		Status(httptest.StatusOK).Body().Equal(expectedResponse)
}
```
## 提示
1. 记得是`RAW BODY`格式
2. 这个例子是需要`config`类型实现`Decode`接口