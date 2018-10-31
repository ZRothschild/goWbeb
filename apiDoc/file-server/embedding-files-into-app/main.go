package main

import (
	"github.com/kataras/iris"
)

// 首先执行以下步骤：
// $ go get -u github.com/shuLhan/go-bindata/...
// $ go-bindata ./assets/...
// $ go build
// $ ./embedding-files-into-app
// 未使用"静态"文件，您可以删除"assets"文件夹并运行该示例。
// 详见 `file-server/embedding-gziped-files-into-app`
func newApp() *iris.Application {
	app := iris.New()
	app.StaticEmbedded("/static", "./assets", Asset, AssetNames)
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080/static/css/bootstrap.min.css
	// http://localhost:8080/static/js/jquery-2.1.1.js
	// http://localhost:8080/static/favicon.ico
	app.Run(iris.Addr(":8080"))
}
