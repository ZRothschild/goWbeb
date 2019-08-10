# `iris`多文件上传示例
## 目录结构
> 主目录`uploadFiles`

```html
    —— templates
        —— upload_form.html
    —— uploads
    —— main.go
```
## 示例代码
> `/templates/upload_form.html`

```html
{% raw %}
<html>
<head>
<title>Upload file</title>
</head>
<body>
	<form enctype="multipart/form-data"
		action="http://127.0.0.1:8080/upload" method="POST">
		<input type="file" name="uploadfile" multiple/> <input type="hidden"
			name="token" value="{{.}}" /> <input type="submit" value="upload" />
	</form>
</body>
</html>
{% endraw %}
```
> `main.go`

```go
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./templates", ".html"))
	//将upload_form.html提供给客户端
	app.Get("/upload", func(ctx iris.Context) {
		//创建一个令牌（可选）。
		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		//使用令牌渲染表单以供您使用
		ctx.View("upload_form.html", token)
	})
	//处理来自upload_form.html的请求
	app.Post("/upload", func(ctx iris.Context) {
		//上传任意数量的文件（表单输入中的multiple属性）。
		//第二个参数完全是可选的,它可以用来根据请求更改文件的名称，
		//在这个例子中，我们将展示如何使用它,通过在上传文件前加上当前用户的ip前缀
		ctx.UploadFormFiles("./uploads", beforeSave)
	})
	app.Post("/upload_manual", func(ctx iris.Context) {
		//获取通过iris.WithPostMaxMemory获取的最大上传值大小。
		maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
		err := ctx.Request().ParseMultipartForm(maxSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
			return
		}
		form := ctx.Request().MultipartForm
		files := form.File["files[]"]
		failures := 0
		for _, file := range files {
			_, err = saveUploadedFile(file, "./uploads")
			if err != nil {
				failures++
				ctx.Writef("failed to upload: %s\n", file.Filename)
			}
		}
		ctx.Writef("%d files uploaded", len(files)-failures)
	})
	//在http://localhost:8080启动服务器，上传限制为32 MB。
	app.Run(iris.Addr(":8080"), iris.WithPostMaxMemory(32<<20))
}

func saveUploadedFile(fh *multipart.FileHeader, destDirectory string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	out, err := os.OpenFile(filepath.Join(destDirectory, fh.Filename),
		os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, src)
}

func beforeSave(ctx iris.Context, file *multipart.FileHeader) {
	ip := ctx.RemoteAddr()
	//确保以某种方式格式化ip
	//可以用于文件名（简单情况）：
	ip = strings.Replace(ip, ".", "_", -1)
	ip = strings.Replace(ip, ":", "_", -1)
	//你可以使用time.Now，为文件添加前缀或后缀
	//基于当前时间戳。
	//即unixTime := time.Now().Unix()
	//使用$IP-为Filename添加前缀
	//不需要更多动作，内部上传者将使用此功能
	//将文件保存到"./uploads"文件夹中的名称。
	file.Filename = ip + "-" + file.Filename
}
```