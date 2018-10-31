# embedding Files

#简介

此包将任何文件转换为可管理的Go源代码。 用于将二进制数据嵌入到go程序中。 在转换为原始字节切片之前，文件数据可选地进行gzip压缩。

它在go-bindata子目录中附带了一个命令行工具。 此工具提供一组命令行选项，用于自定义生成的输出。

## 代码示例 `main.go`

```go
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
```

### 提示

1. 先要安装`github.com/shuLhan/go-bindata/...`
2. 执行`go-bindata ./assets/...`会出现一个`bindata.go`文件
3. 再行`main.go`
