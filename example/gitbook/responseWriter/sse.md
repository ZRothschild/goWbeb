# `SSE`
## `SSE`(`Server-sent Events`)介绍
`SSE`（`Server-sent Events`是`WebSocket`的一种轻量代替方案，使用`HTTP`协议

严格地说，`HTTP`协议是没有办法做服务器推送的，但是当服务器向客户端声明接下来要发送流信息时，客户端就会保持连接打开，`SSE`使用的就是这种原理
### 1.`SSE`能做什么？
理论上，`SSE`和`WebSocket`做的是同一件事情。当你需要用新数据局部更新网络应用时，`SSE`可以做到不需要用户执行任何操作，便可以完成

举例我们要做一个统计系统的管理后台，我们想知道统计数据的实时情况。类似这种更新频繁、 低延迟的场景，`SSE`可以完全满足。

其他一些应用场景：例如邮箱服务的新邮件提醒，微博的新消息推送、管理后台的一些操作实时同步等，`SSE`都是不错的选择

`SSE`是单向通道，只能服务器向客户端发送消息，如果客户端需要向服务器发送消息，则需要一个新的`HTTP`请求。这对比`WebSocket`的双工通道来说，
会有更大的开销。这么一来的话就会存在一个「什么时候才需要关心这个差异？」的问题，如果平均每秒会向服务器发送一次消息的话，那应该选择`WebSocket`。
如果一分钟仅 5 - 6 次的话，其实这个差异并不大
### 2.`SSE`优势
1. 实现一个完整的服务仅需要少量的代码
2. 可以在现有的服务中使用，不需要启动一个新的服务
3. 可以用任何一种服务端语言中使用
4. 基于`HTTP/HTTPS`协议，可以直接运行于现有的代理服务器和认证技术

## 目录结构
> 主目录`sse`

```html
    —— main.go
    —— optional.sse.mini.js.html
```
## 代码示例
> `main.go`

```go
// Package main显示如何通过代理通过SSE向客户端发送连续事件消息。
//阅读详情：https://www.w3schools.com/htmL/html5_serversentevents.asp
// https://robots.thoughtbot.com/writing-a-server-sent-events-server-in-go
package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	//注意: 由于某种原因，最新的vscode-go语言扩展不能提供足够智能帮助（参数文档并转到定义功能）
	//对于`iris.Context`别名，因此如果您使用VS Code，则导入`Context`的原始导入路径，它将执行此操作：
	"github.com/kataras/iris/context"
)
//Broker拥有开放的客户端连接
//在其Notifier频道上侦听传入事件
//并将事件数据广播到所有已注册的连接
type Broker struct {
	//主要事件收集例程将事件推送到此频道
	Notifier chan []byte
	//新的客户端连接
	newClients chan chan []byte
	//关闭客户端连接
	closingClients chan chan []byte
	//客户端连接注册表
	clients map[chan []byte]bool
}

// NewBroker返回一个新的代理工厂
func NewBroker() *Broker {
	b := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	//设置它正在运行 - 收听和广播事件
	go b.listen()
	return b
}

//听取不同的频道并采取相应应对
func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			//新客户端已连接
			//注册他们的消息频道
			b.clients[s] = true
			golog.Infof("Client added. %d registered clients", len(b.clients))
		case s := <-b.closingClients:
			//客户端已离线，我们希望停止向其发送消息。
			delete(b.clients, s)
			golog.Warnf("Removed client. %d registered clients", len(b.clients))
		case event := <-b.Notifier:
			//我们从外面得到了一个新事件
			//向所有连接的客户端发送事件
			for clientMessageChan := range b.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (b *Broker) ServeHTTP(ctx context.Context) {
	//确保编写器支持刷新
	flusher, ok := ctx.ResponseWriter().Flusher()
	if !ok {
		ctx.StatusCode(iris.StatusHTTPVersionNotSupported)
		ctx.WriteString("Streaming unsupported!")
		return
	}
	//设置与事件流相关的header，如果发送纯文本，则可以省略“application/json”
	//如果你开发了一个go客户端，你必须设置：“Accept”：“application/json，text/event-stream”header
	ctx.ContentType("application/json, text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	//我们还添加了跨源资源共享标头，以便不同域上的浏览器仍然可以连接
	ctx.Header("Access-Control-Allow-Origin", "*")
	//每个连接都使用Broker的连接注册表注册自己的消息通道
	messageChan := make(chan []byte)
	//通知我们有新连接的Broker
	b.newClients <- messageChan
	//监听连接关闭以及整个请求处理程序链退出时（此处理程序）并取消注册messageChan。
	ctx.OnClose(func() {
		//从已连接客户端的map中删除此客户端,当这个处理程序退出时
		b.closingClients <- messageChan
	})
	//阻止等待在此连接的消息上广播的消息
	for {
		//写入ResponseWriter
		// Server Sent Events兼容
		ctx.Writef("data: %s\n\n", <-messageChan)
		//或json：data：{obj}
		//立即刷新数据而不是稍后缓冲它
		flusher.Flush()
	}
}

type event struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

const script = `<script type="text/javascript">
if(typeof(EventSource) !== "undefined") {
	console.log("server-sent events supported");
	var client = new EventSource("http://localhost:8080/events");
	var index = 1;
	client.onmessage = function (evt) {
		console.log(evt);
		// it's not required that you send and receive JSON, you can just output the "evt.data" as well.
		dataJSON = JSON.parse(evt.data)
		var table = document.getElementById("messagesTable");
		var row = table.insertRow(index);
		var cellTimestamp = row.insertCell(0);
		var cellMessage = row.insertCell(1);
		cellTimestamp.innerHTML = dataJSON.timestamp;
		cellMessage.innerHTML = dataJSON.message;
		index++;
		window.scrollTo(0,document.body.scrollHeight);
	};
} else {
	document.getElementById("header").innerHTML = "<h2>SSE not supported by this client-protocol</h2>";
}
</script>`

func main() {
	broker := NewBroker()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			now := time.Now()
			evt := event{
				Timestamp: now.Unix(),
				Message:   fmt.Sprintf("Hello at %s", now.Format(time.RFC1123)),
			}
			evtBytes, err := json.Marshal(evt)
			if err != nil {
				golog.Error(err)
				continue
			}
			broker.Notifier <- evtBytes
		}
	}()
	app := iris.New()
	app.Get("/", func(ctx context.Context) {
		ctx.HTML(
			`<html><head><title>SSE</title>` + script + `</head>
				<body>
					<h1 id="header">Waiting for messages...</h1>
					<table id="messagesTable" border="1">
						<tr>
							<th>Timestamp (server)</th>
							<th>Message</th>
						</tr>
					</table>
				</body>
			 </html>`)
	})
	app.Get("/events", broker.ServeHTTP)
	// http://localhost:8080
	// http://localhost:8080/events
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```
> `optional.sse.mini.js.html`

```html
<!-- 你可以把它放到你最喜欢的浏览器 -->
<html>
<head>
    <title>SSE(javascript side)</title>
    <script type="text/javascript">
        var client = new EventSource("http://localhost:8080/events")
        client.onmessage = function (evt) {
            console.log(evt)
        }
    </script>
</head>
<body>
    <h1>打开浏览器控制台（F12）并观察传入的事件消息</h1>
</body>
</html>
```
