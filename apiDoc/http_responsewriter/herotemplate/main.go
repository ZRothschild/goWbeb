package main

import (
	"bytes"
	"./template"
	"github.com/kataras/iris"
)
// $ go get -u github.com/shiyanhui/hero/hero
// $ go run app.go
// 了解更多 https://github.com/shiyanhui/hero/hero
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
		//Hero(github.com/shiyanhui/hero/hero)为此导出了GetBuffer和PutBuffer

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