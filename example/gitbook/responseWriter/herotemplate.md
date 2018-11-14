# Hero
## Hero介绍
`Hero`是一个方便，快速和强大的`go`模板引擎，可以预编译`html`模板以获取代码。
 
> 特征

1. 高性能。
2. 使用方便。
3. 强大。 模板扩展和包含支持。
4. 文件更改时自动编译。
## 目录结构

> 主目录`herotemplate`

```html
    —— template
        —— index.html
        —— user.html
        —— userlist.html
        —— userlistwriter.html
    —— main.go
```
## 代码示例
> `main.go`

```go
package main

import (
	"bytes"
	"./template"
	"github.com/kataras/iris"
)
// $ go get -u github.com/shiyanhui/hero/hero
// $ go run main.go
// 了解更多 https://github.com/shiyanhui/hero
func main() {
	app := iris.New()
	app.Get("/users", func(ctx iris.Context) {
		ctx.Gzip(true)
		ctx.ContentType("text/html")
		var userList = []string{
			"Alice",
			"Bob",
			"Tom",
		}
		//最好使用buffer sync.Pool
		//Hero(github.com/shiyanhui/hero)为此导出了GetBuffer和PutBuffer

		// buffer := hero.GetBuffer()
		// defer hero.PutBuffer(buffer)
		// buffer := new(bytes.Buffer)
		// template.UserList(userList, buffer)
		// ctx.Write(buffer.Bytes())

		//使用io.Writer进行自动缓冲管理（hero built-in 缓冲池），
		// iris context通过其ResponseWriter实现io.Writer
		//这是标准http.ResponseWriter的增强版本
		//但仍然100％兼容，GzipResponseWriter也同样兼容：
		// _, err := template.UserListToWriter(userList, ctx.GzipResponseWriter())
		buffer := new(bytes.Buffer)
		template.UserList(userList, buffer)
		_, err := ctx.Write(buffer.Bytes())
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}
	})
	app.Run(iris.Addr(":8080"))
}
```
> `/template/index.html`

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>
    <body>
        <%@ body { %>
        <% } %>
    </body>
</html>
```
> 文件名称`template/user.html`

```html
<li>
    <%= user %>
</li>
```
> 文件名称`template/userlist.html`

```html
<%: func UserList(userList []string, buffer *bytes.Buffer) %>
<%~ "index.html" %>
<%@ body { %>
    <% for _, user := range userList { %>
        <ul>
            <%+ "user.html" %>
        </ul>
    <% } %>
<% } %>
```
> `/template/userlistwriter.html`

```html
<%: func UserListToWriter(userList []string, w io.Writer) (int, error)%>
<%~ "index.html" %>
<%@ body { %>
    <% for _, user := range userList { %>
        <ul>
            <%+ "user.html" %>
        </ul>
    <% } %>
<% } %>
```
## 提示
1. 引入相应的包，建好所有所有文件
2. 执行相应的命令，所有的`html`文件都会生成相应的`.go`文件
3. [传送们](https://github.com/shiyanhui/hero)