// file: web/routes/hello.go
package routes

import (
	"errors"
	"github.com/kataras/iris/hero"
)

var helloView = hero.View{
	Name: "hello/index.html",
	Data: map[string]interface{}{
		"Title":     "Hello Page",
		"MyMessage": "Welcome to my awesome website",
	},
}

// Hello将返回带有绑定数据的预定义视图。
//`hero.Result`只是一个带有'Dispatch`功能的接口。
//`hero.Response`和`hero.View`是内置的结果类型调度程序
//你甚至可以创建自定义响应调度程序
//实现`github.com/kataras/iris/hero＃Result`接口。
func Hello() hero.Result {
	return helloView
}

//您可以定义标准错误，以便在应用中的任何位置重复使用。
var errBadName = errors.New("bad name")

//您可以将其作为错误返回，甚至其他更好类型的返回
//用hero.Response包装此错误，使其成为hero.Result兼容类型。
var badName = hero.Response{Err: errBadName, Code: 400}

// HelloName 返回 "Hello {name}".
// 例子:
// curl -i http://localhost:8080/hello/iris
// curl -i http://localhost:8080/hello/anything
func HelloName(name string) hero.Result {
	if name != "iris" {
		return badName
	}
	// 或者返回 hero.Response{Text: "Hello " + name}:
	return hero.View{
		Name: "hello/name.html",
		Data: name,
	}
}