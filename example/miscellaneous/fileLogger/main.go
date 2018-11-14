package main

import (
	"os"
	"time"
	"github.com/kataras/iris"
)

//根据日期获取文件名，仅用于容易区分
func todayFilename() string {
	today := time.Now().Format("Jan 02 2006")
	return today + ".txt"
}

func newLogFile() *os.File {
	filename := todayFilename()
	//打开文件，如果服务器重新启动，这将附加到今天的文件。
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	f := newLogFile()
	defer f.Close()
	app := iris.New()
	//将文件作为记录器附加，请记住，iris的app logger只是一个io.Writer。
	//如果需要同时将日志写入文件和控制台，请使用以下代码。
	// app.Logger().SetOutput(io.MultiWriter(f, os.Stdout))
	app.Logger().SetOutput(f)
	app.Get("/ping", func(ctx iris.Context) {
		// for the sake of simplicity, in order see the logs at the ./_today_.txt
		ctx.Application().Logger().Infof("Request path: %s", ctx.Path())
		ctx.WriteString("pong")
	})
	//访问http://localhost:8080/ping
	//打开文件 ./logs{TODAY}.txt file.
	if err := app.Run(iris.Addr(":8080"), iris.WithoutBanner, iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
		app.Logger().Warn("Shutdown with error: " + err.Error())
	}
}