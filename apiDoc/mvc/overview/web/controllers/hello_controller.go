//文件: web/controllers/hello_controller.go
package controllers

import (
	"errors"
	"github.com/kataras/iris/mvc"
)

// HelloController是我们的示例控制器
//它处理 GET：/hello和GET：/hello/{name}
type HelloController struct{}

var helloView = mvc.View{
	Name: "hello/index.html",
	Data: map[string]interface{}{
		"Title":     "Hello Page",
		"MyMessage": "Welcome to my awesome website",
	},
}
//Get将返回带有绑定数据的预定义视图
//`mvc.Result`只是一个带有'Dispatch`功能的接口。
//`mvc.Response`和`mvc.View`是内置的结果类型调度程序
//你甚至可以创建自定义响应调度程序
//实现`github.com/kataras/iris/hero＃Result`接口。
func (c *HelloController) Get() mvc.Result {
	return helloView
}

//您可以定义标准错误，以便在应用中的任何位置重复使用
var errBadName = errors.New("bad name")

//您可以将其作为错误返回
//使用mvc.Response包装此错误，使其成为mvc.Result兼容类型。
var badName = mvc.Response{Err: errBadName, Code: 400}

// GetBy 返回 "Hello {name}" response.
// 示例:
// curl -i http://localhost:8080/hello/iris
// curl -i http://localhost:8080/hello/anything
func (c *HelloController) GetBy(name string) mvc.Result {
	if name != "iris" {
		return badName
		// 或
		// GetBy(name string) (mvc.Result, error) {
		//	return nil, errBadName
		// }
	}
	// return mvc.Response{Text: "Hello " + name} 或:
	return mvc.View{
		Name: "hello/name.html",
		Data: name,
	}
}