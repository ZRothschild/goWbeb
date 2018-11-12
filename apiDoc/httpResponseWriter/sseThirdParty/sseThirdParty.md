# `SSE`第三方包使用
## 目录结构
> 主目录`sseThirdParty`
```html
    —— main.go
```
## 代码示例
> `main.go`
```go
package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/r3labs/sse"
)

//首先安装sse第三方软件包（如果您不喜欢这种方法，可以使用其他软件包或继续使用sse示例）
// $ go get -u github.com/r3labs/sse
func main() {
	app := iris.New()
	s := sse.New()
	/*
		这会在调度程序内部创建一个新流。
		没有消费者，发布消息,这个频道什么都不做。
		启动iris处理程序后，客户端可以连接到此流
		通过将stream指定为url参数，如下所示：
		http://localhost:8080/events?stream=messages
	*/
	s.CreateStream("messages")
	app.Any("/events", iris.FromStd(s.HTTPHandler))
	go func() {
		//您设计何时向客户端发送消息，
		//这里我们等待5秒钟发送第一条消息
		//为了给你时间打开一个浏览器窗口..
		time.Sleep(5 * time.Second)
		//将有效负载发布到流。
		s.Publish("messages", &sse.Event{
			Data: []byte("ping"),
		})
		time.Sleep(3 * time.Second)
		s.Publish("messages", &sse.Event{
			Data: []byte("second message"),
		})
		time.Sleep(2 * time.Second)
		s.Publish("messages", &sse.Event{
			Data: []byte("third message"),
		})
	}()
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
/*对于golang SSE客户端，您可以查看：https://github.com/r3labs/sse#example-client */
```