### Websockets
WebSocket是一种通过TCP连接实现双向持久通信通道的协议。
它可用于聊天，股票行情，游戏等应用程序，您可以在Web应用程序中使用实时功能

[查看代码示例.](https://github.com/kataras/iris/tree/master/_examples/websocket)

#### 何时使用它

需要直接使用套接字连接时，请使用WebSockets。
例如，您可能需要实时游戏的最佳性能。
 
#### 如何使用
```go
import the "github.com/kataras/iris/websocket"
```
* import "github.com/kataras/iris/websocket"
* 配置websockets包
* 接收websocket包
* 发送和接收消息

```go
    func main() {
        ws := websocket.New(websocket.Config{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        })
    }
```

#### 完整的配置
```go
    //配置websocket服务器配置
    // 所有这些都是可选的
    type Config struct {
        // IDGenerator用于创建（以及稍后设置）
        //每个传入的websocket连接（客户端）的ID。
        //请求是一个参数，您可以使用它来生成ID（例如，来自标题）。
        //如果为空，则由DefaultIDGenerator生成ID：randomString（64）
        IDGenerator func(ctx context.Context) string
        Error       func(w http.ResponseWriter, r *http.Request, status int, reason error)
        CheckOrigin func(r *http.Request) bool
        // HandshakeTimeout指定握手完成的持续时间。
        HandshakeTimeout time.Duration
        //允许WriteTimeout时间向连接写入消息。
        // 0表示没有超时。
        //默认值为0
        WriteTimeout time.Duration
        //允许ReadTimeout时间从连接中读取消息。
        // 0表示没有超时。
        //默认值为0
        ReadTimeout time.Duration
       // PongTimeout允许从连接中读取下一个pong消息。
        //默认值为60 * time.Second
        PongTimeout time.Duration
        // PingPeriod将ping消息发送到此期间的连接。必须小于PongTimeout。
        //默认值为60 * time.Second
        PingPeriod time.Duration
         // MaxMessageSize连接允许的最大消息大小。
        //默认值为1024
        MaxMessageSize int64
       // BinaryMessages将其设置为true，以表示二进制数据消息而不是utf-8文本
        //兼容，如果您想使用Connection的EmitMessage将自定义二进制数据发送到客户端，就像本机服务器 - 客户端通信一样。
        //默认为false
        BinaryMessages bool
        // ReadBufferSize是下划线阅读器的缓冲区大小
        //默认值为4096
        ReadBufferSize int
        // WriteBufferSize是下划线编写器的缓冲区大小
        //默认值为4096
        WriteBufferSize int
        // EnableCompression指定服务器是否应尝试协商每个
        //消息压缩（RFC 7692）。将此值设置为true则不会
        //保证支持压缩。目前只有“没有背景
        //支持“接管”模式。
        EnableCompression bool
        //子协议按顺序指定服务器支持的协议
        //偏好。如果设置了此字段，则Upgrade方法通过使用协议选择此列表中的第一个匹配来协商子协议
        //客户要求。
        Subprotocols []string
    }
    Accept WebSocket requests & send & receive messages
    import (
        "github.com/kataras/iris"
        "github.com/kataras/iris/websocket"
    )
    func main() {
        ws := websocket.New(websocket.Config{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        })
        ws.OnConnection(handleConnection)
        app := iris.New()
        // 在端点上注册服务器。
        // 请参阅websockets.html中的内联JavaScript代码，此端点用于连接到服务器。
        app.Get("/echo", ws.Handler())
        //提供javascript built'n客户端库，
        //请参阅weboskcets.html脚本标记，使用此路径。
        app.Any("/iris-ws.js", func(ctx iris.Context) {
            ctx.Write(websocket.ClientSource)
        })
    }
    func handleConnection(c websocket.Connection) {
        //从浏览器中读取事件
        c.On("chat", func(msg string) {
            // 将消息打印到控制台，c .Context（）是iris的http上下文。
            fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
            //将消息写回客户端消息所有者：
            // c.Emit("chat", msg)
            c.To(websocket.Broadcast).Emit("chat", msg)
        })
    }
```
