# `iris`自定义结构体映射获取`Form`表单请求数据
## 目录结构
> 主目录`readForm`

```html
    —— templates
        —— form.html
    —— main.go
```
## 代码示例
> `templates/form.html`

```html
<!DOCTYPE html>
<head>
<meta charset="utf-8">
</head>
<body>
	<form action="/form_action" method="post">
		Username: <input type="text" name="Username" /> <br />
		Mail: <input type="text" name="Mail" /> <br /> 
		Select one or more:  <br/> 
		<select multiple="multiple" name="mydata">
			<option value='one'>One</option>
			<option value='two'>Two</option>
			<option value='three'>Three</option>
			<option value='four'>Four</option>
		</select>
		<hr />
		<input type="submit" value="Send data"/>
	</form>
</body>
</html>
```
> `main.go`

```go
// package main包含一个关于如何使用ReadForm的示例，但使用相同的方法可以执行ReadJSON和ReadJSON
package main

import (
	"github.com/kataras/iris"
)

type Visitor struct {
	Username string
	Mail     string
	Data     []string `form:"mydata"`
}

func main() {
	app := iris.New()
	//设置视图html模板引擎
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))
	app.Get("/", func(ctx iris.Context) {
		if err := ctx.View("form.html"); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}
	})
	app.Post("/form_action", func(ctx iris.Context) {
		visitor := Visitor{}
		err := ctx.ReadForm(&visitor)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}
		ctx.Writef("Visitor: %#v", visitor)
	})
	app.Post("/post_value", func(ctx iris.Context) {
		username := ctx.PostValueDefault("Username", "iris")
		ctx.Writef("Username: %s", username)
	})
	app.Run(iris.Addr(":8080"))
}
```