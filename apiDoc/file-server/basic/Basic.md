# Basic

### 代码示例

```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Favicon("./assets/favicon.ico")
	//启用gzip，可选：
	//如果在`StaticXXX`处理程序之前使用的话内容字节范围功能消失了。
	//推荐：特别关闭大文件 当服务器内存不足时 打开中等大小的文件
	//或者对于大型文件（如果它们已经压缩），
	// i.e "zippedDir/file.gz"
	// app.Use(iris.Gzip)
	//第一个参数是请求路径，第二个参数是系统目录

	// app.StaticWeb("/css", "./assets/css")
	// app.StaticWeb("/js", "./assets/js")
	app.StaticWeb("/static", "./assets")
	// http://localhost:8080/static/css/main.css
	// http://localhost:8080/static/js/jquery-2.1.1.js
	// http://localhost:8080/static/favicon.ico
	app.Run(iris.Addr(":8080"))
	//路由不允许.StaticWeb("/"，"./ assets")
	//要了解如何包装路由器以实现
	//根路径上的通配符，请参阅"single-page-application"。
}
```