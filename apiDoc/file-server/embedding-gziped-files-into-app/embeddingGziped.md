# embedding Gziped

### 示例代码

```go
package main

import (
	"github.com/kataras/iris"
)

//注意：与"embedding-files-into-app"示例压缩选择
//首先执行以下步骤：
// $ go get -u github.com/kataras/bindata/cmd/bindata
// $ bindata ./assets/...
// $ go build
// $ ./embedding-gziped-files-into-app
// 未使用"静态"文件，您可以删除"assets"文件夹并运行该示例。

func newApp() *iris.Application {
	app := iris.New()
	//注意`GzipAsset`和`GzipAssetNames`不同于`go-bindata`的`Asset`和`AssetNames，
	//这意味着你可以同时使用`go-bindata`和`bindata`工具，
	//`go-bindata`可以用于视图引擎的`Binary`方法
	//和带有`StaticEmbeddedGzip`的`bindata`（比使用`go-bindata`的StaticEmbeded快x8倍）。
	app.StaticEmbeddedGzip("/static", "./assets", GzipAsset, GzipAssetNames)
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080/static/css/bootstrap.min.css
	// http://localhost:8080/static/js/jquery-2.1.1.js
	// http://localhost:8080/static/favicon.ico
	app.Run(iris.Addr(":8080"))
}
```