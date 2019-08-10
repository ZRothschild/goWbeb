####包装路由器

443/5000
非常罕见，您可能永远不需要它，但无论如何您都需要它。(以备不时之需)

有时您需要覆盖或决定是否将在传入请求上执行路由器。 如果你以前有过使用net / http和其他web框架的经验，
这个函数会熟悉你（它有net / http中间件的形式，但它不接受下一个处理程序，而是接受Router作为函数 是否被执行）。

```go
    // WrapperFunc用作预期的输入参数签名
    //用于WrapRouter。 它是一个“低级”签名，与net / http兼容。
    //它用于运行或不运行基于自定义逻辑的路由器。
    type WrapperFunc func(w http.ResponseWriter, r *http.Request, firstNextIsTheRouter http.HandlerFunc)
    // WrapRouter在主路由器的顶部添加了一个包装器。
    //通常，当需要使用像CORS这样的中间件包装整个应用程序时，它对第三方中间件很有用。
    //开发人员可以添加多个包装器，这些包装器的执行从最后到第一个。
    //这意味着第二个包装器将包装第一个，依此类推。
    //在构建之前。
    func WrapRouter(wrapperFunc WrapperFunc)
```
Iris的路由器基于HTTP方法搜索其路由，路由器包装器可以覆盖该行为并执行自定义代码。
示例代码:

```go
func main() {
    app := iris.New()
    app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
        ctx.HTML("<b>Resource Not found</b>")
    })
    app.Get("/", func(ctx iris.Context) {
        ctx.ServeFile("./public/index.html", false)
    })
    app.Get("/profile/{username}", func(ctx iris.Context) {
        ctx.Writef("Hello %s", ctx.Params().Get("username"))
    })
   //提供来自根“/”的文件，如果我们使用.StaticWeb它可以覆盖
   //由于下划线需要通配符，所有路由。
   //在这里，我们将看到如何绕过这种行为
   //通过创建一个新的文件服务器处理程序和
   //为路由器设置包装器（如“低级”中间件）
   //为了手动检查我们是否想要正常处理路由器
   //或者执行文件服务器处理程序。
   //使用.StaticHandler
   //它与StaticWeb相同，但不是
   //注册路由，它只返回处理程序。
    fileServer := app.StaticHandler("./public", false, false)
    //使用本机net / http处理程序包装路由器。
    //如果url不包含任何“。” （即：.css，.js ......）
    //（取决于应用程序，您可能需要添加更多文件服务器异常），
    //然后处理程序将执行负责的路由器
   //注册路线（看“/”和“/ profile / {username}”）
   //如果没有，那么它将根据根“/”路径提供文件。
    app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
        path := r.URL.Path
       //请注意，如果path的后缀为“index.html”，则会自动重定向到“/”，
       //所以我们的第一个处理程序将被执行。
        if !strings.Contains(path, ".") { 
            //如果它不是资源，那么就像正常情况一样继续使用路由器. <-- IMPORTANT
            router(w, r)
            return
        }
        //获取并释放上下文，以便使用它来执行我们的文件服务器
        //记得：我们使用net / http.Handler，因为我们在路由器本身之前处于“低级别”。
        ctx := app.ContextPool.Acquire(w, r)
        fileServer(ctx)
        app.ContextPool.Release(ctx)
    })
    // http://localhost:8080
    // http://localhost:8080/index.html
    // http://localhost:8080/app.js
    // http://localhost:8080/css/main.css
    // http://localhost:8080/profile/anyusername
    app.Run(iris.Addr(":8080"))
   //注意：在这个例子中我们只看到一个用例，
   //你可能想要.WrapRouter或.Downgrade以绕过虹膜的默认路由器，即：
   //您也可以使用该方法设置自定义代理。
   //如果您只想在除root之外的其他路径上提供静态文件
   //你可以使用StaticWeb, i.e:
    //                          .StaticWeb("/static", "./public")
    // ________________________________requestPath, systemPath
}
```