####动态路由参数
Iris 拥有你从未遇到过的 简单的,强大的路由.

同时, Iris有自己的路径（就像编程语言一样），用于路由的路径语法及其路径参数解析和评估. 我们简称为"macros".

怎么样? 它计算了它的需求，如果没有需要任何特殊的正则表达式 那么它只是用低级路径语法注册路由，否则它预先编译正则表达式并添加必要的中间件。
这意味着相对于其他路由器或Web框架 您的性能成本为零 。

路径路径参数的标准macro类型
```go
+------------------------+
| {param:string}         |
+------------------------+

string 类型
任意字符串
+------------------------+
| {param:int}            |
+------------------------+
int 类型
支持数字 (0-9)(这里是组合123或者1234其他整数类型于此相同)

+------------------------+
| {param:long}           |
+------------------------+
int64 type
仅仅数字 (0-9)

+------------------------+
| {param:boolean}        |
+------------------------+
bool 类型
仅仅"1" 或者 "t" 或者 "T" 或者 "TRUE"或者 "true" 或者 "True"
或者 "0" 或者 "f" 或者 "F" 或者 "FALSE" 或者 "false" 或者 "False"

+------------------------+
| {param:alphabetical}   |
+------------------------+
alphabetical/letter (拼音或者字母)类型
letters only (大写或者小写)

+------------------------+
| {param:file}           |
+------------------------+
file 类型
letters (大写或者小写)
numbers (0-9)
underscore (_)
dash (-)
point (.)

没有空格 ！或其他字符

+------------------------+
| {param:path}           |
+------------------------+
path 类型
anything,应该是最后一部分，多个路径段,
示例: /path1/path2/path3 , ctx.Params().Get("param") == "/path1/path2/path3"
```

如果缺少类型，则参数的类型默认为字符串，因此{param} == {param：string}。

如果在该类型上找不到函数，则使用字符串macro类型的函数。

除了Iris提供基本类型和一些默认的“macro功能”你也可以注册自己的func！

注册命名路径参数功能
```go
app.Macros().Int.RegisterFunc("min", func(argument int) func(paramValue string) bool {
    // [...]
    return true 
    // -> true 意味着通过验证, false 表示无效的消息404，或者如果“其他500”被附加到macors语法，则内部服务器错误.
})
```
在 func(argument ...) 你可以有任何标准类型, 它将在服务器启动之前进行验证，因此不关心那里的任何性能成本，它在服务时运行的唯一事情就是返回func（paramValue string）bool。

{param:string equal(iris)} , "iris" 在这里是一个参数:
```go
app.Macros().String.RegisterFunc("equal", func(argument string) func(paramValue string) bool {
    return func(paramValue string){ return argument == paramValue }
})
```
示例代码:
```go
 app := iris.New()
// 您可以使用“string”类型，该类型对于可以是任何内容的单个路径参数有效
app.Get("/username/{name}", func(ctx iris.Context) {
    ctx.Writef("Hello %s", ctx.Params().Get("name"))
}) //  {name:string}

//注册我们的第一个附加到int macro类型的宏
// "min" = 当前函数名字
// "minValue" = 函数的参数
// func(string) bool = macro's 的路径参数评估器，这在服务时执行
// 用户使用min（...）macros参数函数请求包含：int macros类型的路径。
app.Macros().Int.RegisterFunc("min", func(minValue int) func(string) bool {
    // 在此之前做任何事情[...]
    //在这种情况下，我们不需要做任何事情
    return func(paramValue string) bool {
        n, err := strconv.Atoi(paramValue)
        if err != nil {
            return false
        }
        return n >= minValue
    }
})
// http://localhost:8080/profile/id>=1
// 这将抛出404，即使它被发现为路线 : /profile/0, /profile/blabla, /profile/-1
// macros 参数函数当然是可选的

app.Get("/profile/{id:int min(1)}", func(ctx iris.Context) {
    //  第二个参数是错误的，因为我们使用 macros 它总是为nil
    // 验证已经发生了.
    id, _ := ctx.Params().GetInt("id")
    ctx.Writef("Hello id: %d", id)
})

// 更改每个路由的macros评估程序的错误代码：
app.Get("/profile/{id:int min(1)}/friends/{friendid:int min(1) else 504}", func(ctx iris.Context) {
    id, _ := ctx.Params().GetInt("id")
    friendid, _ := ctx.Params().GetInt("friendid")
    ctx.Writef("Hello id: %d looking for friend id: ", id, friendid)
}) // 如果没有传递所有路由的macros，这将抛出504错误代码而不是404.

// http://localhost:8080/game/a-zA-Z/level/0-9
//记住，字母只是小写或大写字母。
app.Get("/game/{name:alphabetical}/level/{level:int}", func(ctx iris.Context) {
    ctx.Writef("name: %s | level: %s", ctx.Params().Get("name"), ctx.Params().Get("level"))
})

//让我们使用一个简单的自定义regexp来验证单个路径参数
//它的值只是小写字母。

// http://localhost:8080/lowercase/anylowercase
app.Get("/lowercase/{name:string regexp(^[a-z]+)}", func(ctx iris.Context) {
    ctx.Writef("name should be only lowercase, otherwise this handler will never executed: %s", ctx.Params().Get("name"))
})

// http://localhost:8080/single_file/app.js
app.Get("/single_file/{myfile:file}", func(ctx iris.Context) {
    ctx.Writef("file type validates if the parameter value has a form of a file name, got: %s", ctx.Params().Get("myfile"))
})

// http://localhost:8080/myfiles/any/directory/here/
// 这是唯一接受任意数量路径段的macro类型。
app.Get("/myfiles/{directory:path}", func(ctx iris.Context) {
    ctx.Writef("path type accepts any number of path segments, path after /myfiles/ is: %s", ctx.Params().Get("directory"))
})

app.Run(iris.Addr(":8080"))
}
```

路径参数名称应仅包含字母。不允许使用“_”和数字等符号。

最后，不要将ctx.Params（）与ctx.Values（）混淆。路径参数的值转到ctx.Params（）和上下文的本地存储 可以用来在处理程序和中间件之间进行通信，转到ctx.Values（）。