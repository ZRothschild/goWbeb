#### MVC 结合Websocket

```go
    package main
    import (
        "sync/atomic"
        "github.com/kataras/iris"
        "github.com/kataras/iris/mvc"
        "github.com/kataras/iris/websocket"
    )
    func main() {
        app := iris.New()
        // 加载模板
        app.RegisterView(iris.HTML("./views", ".html"))

        // 渲染视图.
        app.Get("/", func(ctx iris.Context) {
            ctx.View("index.html")
        })

        mvc.Configure(app.Party("/websocket"), configureMVC)

        // http://localhost:8080
        app.Run(iris.Addr(":8080"))
    }
    func configureMVC(m *mvc.Application) {
        ws := websocket.New(websocket.Config{})
        // http://localhost:8080/websocket/iris-ws.js
        m.Router.Any("/iris-ws.js", websocket.ClientHandler())

      //这将绑定ws.Upgrade的结果，这是一个websocket.Connection
        //由`m.Handle`服务的控制器。
        m.Register(ws.Upgrade)
        m.Handle(new(websocketController))
    }
    var visits uint64
    func increment() uint64 {
        return atomic.AddUint64(&visits, 1)
    }

    func decrement() uint64 {
        return atomic.AddUint64(&visits, ^uint64(0))
    }

    type websocketController struct {
        //注意你也可以使用匿名字段，无所谓，binder会找到它。
        //这是当前的websocket连接，每个客户端都有自己的* websocketController实例。
        Conn websocket.Connection
    }

    func (c *websocketController) onLeave(roomName string) {
        //visits--
        newCount := decrement()

        //这将在所有客户端上调用“visit”事件，当前客户端除外
        //（它不能因为它已经离开但是对于任何情况都使用这种类型的设计）
        c.Conn.To(websocket.Broadcast).Emit("visit", newCount)
    }

    func (c *websocketController) update() {
        // visits++
        newCount := increment()

       //这将在所有客户端上调用“visit”事件，包括当前事件
        //使用'newCount'变量。
        //
        //你有很多方法可以做到更快，例如你可以发送一个新的访问者
        //并且客户端可以自行增加，但这里我们只是“展示”websocket控制器。
        c.Conn.To(websocket.All).Emit("visit", newCount)
    }
    func (c *websocketController) Get( /* websocket.Connection could be lived here as well, it doesn't matter */ ) {
        c.Conn.OnLeave(c.onLeave)
        c.Conn.On("visit", c.update)

        // 在所有事件回调注册后调用它。
        c.Conn.Wait()
    }
```
