package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
	"github.com/kataras/iris"
)

const maxSize = 5 << 20 // 5MB

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./templates", ".html"))
	//将upload_form.html提供给客户端。
	app.Get("/upload", func(ctx iris.Context) {
		//创建一个令牌（可选）。
		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		//使用令牌渲染表单以供您使用。
		//ctx.ViewData(""，token)
		//或者在`View`方法中添加第二个参数。
		//令牌将作为{{.}}传递到模板中。
		ctx.View("upload_form.html", token)
	})
	/* 继续之前阅读。
	0.默认发布最大大小为32MB，您可以使用`app.Run`中的`iris.WithPostMaxMemory(maxSize)`配置器扩展它以读取更多数据，
	请注意，这不足以满足您的需求，请阅读以下内容。
	1.检查大小的更快方法是使用`ctx.GetContentLength()`来返回整个请求的大小.(加上一个逻辑数字,如2MB甚至10MB，其余大小如headers).
	你可以创建一个中间件使其适应任何必要的处理程序
	myLimiter := func(ctx iris.Context) {
		if ctx.GetContentLength() > maxSize { // + 2 << 20 {
			ctx.StatusCode(iris.StatusRequestEntityTooLarge)
			return
		}
		ctx.Next()
	}
	app.Post("/upload", myLimiter, myUploadHandler)
	大多数客户端将设置"Content-Length"header(如浏览器),但确保任何客户端始终更好无法发送您的服务器无法或不想处理的数据.
	这可以使用`app.Use(LimitRequestBodySize(maxSize))`(作为app或路由中间件)
	或者`ctx.SetMaxRequestBodySize(maxSize)`来限制基于特定处理程序内部的自定义逻辑的请求，它们是相同的，
	参见下文。
	2.您可以使用`ctx.SetMaxRequestBodySize(maxSize)`强制限制处理程序内的请求体大小，
	如果传入的数据较大（大多数客户端将其作为"连接重置"接收），这将强制关闭连接，
	使用它来确保客户端不会发送服务器不能或不想接受的数据作为后备。
	app.Post("/upload", iris.LimitRequestBodySize(maxSize), myUploadHandler)
	或
	app.Post("/upload", func(ctx iris.Context){
		ctx.SetMaxRequestBodySize(maxSize)
		// [...]
	})
	3.另一种方法是接收数据并检查`ctx.FormFile`的第二个返回值的'Size`值，即`info.Size`，这会给你确切的文件大小，
	而不是整个传入的请求数据长度。
	app.Post("/", func(ctx iris.Context){
		file, info, err := ctx.FormFile("uploadfile")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer file.Close()
		if info.Size > maxSize {
			ctx.StatusCode(iris.StatusRequestEntityTooLarge)
			return
		}
		// [...]
	})
	*/
	//处理来自upload_form.html的请求数据处理
	app.Post("/upload", iris.LimitRequestBodySize(maxSize+1<<20), func(ctx iris.Context) {
		// Get the file from the request.
		file, info, err := ctx.FormFile("uploadfile")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer file.Close()
		fname := info.Filename
		//创建一个具有相同名称的文件
		//假设你有一个名为'uploads'的文件夹
		out, err := os.OpenFile("./uploads/"+fname,
			os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
			return
		}
		defer out.Close()
		io.Copy(out, file)
	})
	//在http//localhost:8080启动服务器，上传限制为5MB。
	app.Run(iris.Addr(":8080") /* 0.*/, iris.WithPostMaxMemory(maxSize))
}