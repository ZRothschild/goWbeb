# `iris`自定义结构体映射获取`xml`格式请求数据
## 代码示例
### 主目录`readXml`
> 文件名称`main.go`
```go
package main

import (
	"encoding/xml"
	"github.com/kataras/iris"
	"fmt"
)

func main() {
	app := newApp()
	//使用Postman或其他什么来做POST请求
	//使用RAW BODY到http//localhost:8080/
	/*
		<person name="Winston Churchill" age="90">
			<description>Description of this person, the body of this inner element.</description>
		</person>
	*/
	//和Content-Type到application/xml（可选但最好设置）
	//响应应该是：
	// 接收: main.person{XMLName:xml.Name{Space:"", Local:"person"}, Name:"Winston Churchill", Age:90,
	// Description:"Description of this person, the body of this inner element."}
	app.Run(iris.Addr(":8080"), iris.WithOptimizations)
}

func newApp() *iris.Application {
	app := iris.New()
	app.Post("/", handler)
	return app
}

//简单的xml例子，请访问https://golang.org/pkg/encoding/xml
type person struct {
	XMLName     xml.Name `xml:"person"`      //元素名称
	Name        string   `xml:"name,attr"`   //，attr属性。
	Age         int      `xml:"age,attr"`    //，attr属性。
	Description string   `xml:"description"` //内部元素名称，值是它的主体。
}

func handler(ctx iris.Context) {
	fmt.Println(ctx.GetCurrentRoute())
	var p person
	if err := ctx.ReadXML(&p); err != nil {
		fmt.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Writef("Received: %#+v", p)
}
```
> 文件名称`main_test.go`
```go
package main

import (
	"testing"
	"github.com/kataras/iris/httptest"
)

func TestReadXML(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)
	expectedResponse := `Received: main.person{XMLName:xml.Name{Space:"", Local:"person"}, Name:"Winston Churchill", Age:90, Description:"Description of this person, the body of this inner element."}`
	send := `<person name="Winston Churchill" age="90"><description>Description of this person, the body of this inner element.</description></person>`
	e.POST("/").WithText(send).Expect().
		Status(httptest.StatusOK).Body().Equal(expectedResponse)
}
```
## 目录结构
> 主目录`readXml`
```html
    —— main.go
    —— main_test.go
```