# casbin 实现基于角色的 HTTP 权限控制

### [casbin](https://github.com/casbin/casbin)介绍 

[casbin](https://github.com/casbin/casbin)是由北大的一位博士生主导开发的一个基于`Go`语言的权限控制库。支持 `ACL`，`RBAC`，`ABAC` 等常用的访问控制模型。

[casbin](https://github.com/casbin/casbin)是`Golang`项目的强大而高效的开源访问控制库。 它支持基于各种访问控制模型实施授权。

[casbin](https://github.com/casbin/casbin)的核心是一套基于`PERM metamodel`(`Policy`, `Effect`, `Request`, `Matchers`)的`DSL`。`Casbin`从用这种`DSL`定义的配置
文件中读取访问控制模型，作为后续权限验证的基础

#### `Casbin`做了什么

1. 支持自定义请求的格式，默认的请求格式为`{subject, object, action}`。
2. 具有访问控制模型`model`和策略`policy`两个核心概念。
3. 支持`RBAC`中的多层角色继承，不止主体可以有角色，资源也可以具有角色。
4. 支持超级用户，如`root`或`Administrator`，超级用户可以不受授权策略的约束访问任意资源。
5. 支持多种内置的操作符，如`keyMatch`，方便对路径式的资源进行管理，如`/foo/bar`可以映射到`/foo*`

#### `Casbin`不做的事情

1. 身份认证`authentication`(即验证用户的用户名、密码)，`casbin`只负责访问控制。应该有其他专门的组件负责身份认证，然后由`casbin`进行访问
控制，二者是相互配合的关系。
2. 管理用户列表或角色列表。`Casbin`认为由项目自身来管理用户、角色列表更为合适，用户通常有他们的密码，但是`Casbin`的设计思想并不是把
它作为一个存储密码的容器。而是存储`RBAC`方案中用户和角色之间的映射关系。

### 配置示例

#### 模型与策略定制

```smartyconfig
//sub   "alice"// 想要访问资源的用户.
//obj  "data1" // 要访问的资源.
//act  "read"  // 用户对资源执行的操作.

# Request definition
[request_definition]
r = sub, obj, act

# Policy definition
[policy_definition]
p = sub, obj, act

# Policy effect
[policy_effect]
e = some(where (p.eft == allow))

# Matchers
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```



可以看到这个配置文件主要定义了`Request`和`Policy`的组成结构.`Policy effect`和`Matchers`则灵活的多，可以包含一些自定义的表达式
比如我们要加入一个名叫`root`的超级管理员，就可以这样写:

```smartyconfig
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act || r.sub == "root"
``` 
又比如我们可以用正则匹配来判断权限是否匹配:

```smartyconfig
[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
``` 

#### 具体规则设置

```smartyconfig
p, alice, data1, read
p, bob, data2, write
```
> 意思就是 alice 可以读 data1，bob 可以写 data2

### 示例 

> 模型与策略定制 `test.conf`

```smartyconfig
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

> 具体规则设置 `test.csv`

```smartyconfig
p, admin, domain1, data1, read
p, admin, domain1, data1, write
p, admin, domain2, data2, read
p, admin, domain2, data2, write
g, alice, admin, domain1
g, bob, admin, domain2
```

如上所示，alice 和 bob 分别是 domian1 和 domain2 的管理员

### iris 示例代码

#### 中间件格式 `错误返回forbidden`

> 模型与策略定制 `casbinmodel.conf`

```smartyconfig
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```
> 具体规则设置 `casbinpolicy.csv`

```smartyconfig
p, alice, /dataset1/*, GET
p, alice, /dataset1/resource1, POST
p, bob, /dataset2/resource1, *
p, bob, /dataset2/resource2, GET
p, bob, /dataset2/folder1/*, POST
```
> iris 代码示例

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/casbin/casbin"
	cm "github.com/iris-contrib/middleware/casbin"
)

// $ go get github.com/casbin/casbin
// $ go run main.go
// Enforcer映射模型和casbin服务的策略，我们也在main_test上使用此变量。
var Enforcer = casbin.NewEnforcer("casbinmodel.conf", "casbinpolicy.csv")

func newApp() *iris.Application {
	casbinMiddleware := cm.New(Enforcer)
	app := iris.New()
	app.Use(casbinMiddleware.ServeHTTP)
	app.Get("/", hi)
	app.Get("/dataset1/{p:path}", hi) // p, alice, /dataset1/*, GET
	app.Post("/dataset1/resource1", hi)
	app.Get("/dataset2/resource2", hi)
	app.Post("/dataset2/folder1/{p:path}", hi)
	app.Any("/dataset2/resource1", hi)
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

func hi(ctx iris.Context) {
	ctx.Writef("Hello %s", cm.Username(ctx.Request()))
}
```

#### 路由修饰模式  `错误返回403`

> 模型与策略定制 `casbinmodel.conf`

```smartyconfig
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```
> 具体规则设置 `casbinpolicy.csv`

```smartyconfig
p, alice, /dataset1/*, GET
p, alice, /dataset1/resource1, POST
p, bob, /dataset2/resource1, *
p, bob, /dataset2/resource2, GET
p, bob, /dataset2/folder1/*, POST
p, cathrin, /dataset2/resource2, GET
p, dataset1_admin, /dataset1/*, *
g, cathrin, dataset1_admin
```
> iris 代码示例

```go
package main

import (
	"github.com/kataras/iris"

	"github.com/casbin/casbin"
	cm "github.com/iris-contrib/middleware/casbin"
)

// $ go get github.com/casbin/casbin
// $ go run main.go
// Enforcer映射模型和casbin服务的策略，我们也在main_test上使用此变量。
var Enforcer = casbin.NewEnforcer("casbinmodel.conf", "casbinpolicy.csv")

func newApp() *iris.Application {
	casbinMiddleware := cm.New(Enforcer)
	app := iris.New()
	app.WrapRouter(casbinMiddleware.Wrapper())
	app.Get("/", hi)
	app.Any("/dataset1/{p:path}", hi) // p, dataset1_admin, /dataset1/*, * && p, alice, /dataset1/*, GET
	app.Post("/dataset1/resource1", hi)
	app.Get("/dataset2/resource2", hi)
	app.Post("/dataset2/folder1/{p:path}", hi)
	app.Any("/dataset2/resource1", hi)

	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

func hi(ctx iris.Context) {
	ctx.Writef("Hello %s", cm.Username(ctx.Request()))
}
```

### 提示

1. 以上的`go iris`都是使用`Basic Auth`,用`postman`测试请选择`Authorization`选项
2. `*.conf`文件是配置规则模型，`*.csv`是具体规则的体现，当然也可不使用这些东西，用户数据或者其他代替
3. 解释一下我对这些的理解

```smartyconfig
p, alice, /dataset1/*, GET         //alice 用户有对 method为GET路径满足 /dataset1/*的访问权限 下面同理
p, alice, /dataset1/resource1, POST
p, bob, /dataset2/resource1, *
p, bob, /dataset2/resource2, GET
p, bob, /dataset2/folder1/*, POST
p, cathrin, /dataset2/resource2, GET
p, dataset1_admin, /dataset1/*, *
g, cathrin, dataset1_admin  //cathrin用户属于dataset1_admin组，也就是dataset1_admin能访问的cathrin都能访问，反之不然
```

```smartyconfig
[request_definition]        //请求定义
r = sub, obj, act

[policy_definition]         //策略定义，也就是*.cvs文件 p 定义的格式
p = sub, obj, act

[role_definition]           //组定义，也就是*.cvs文件 g 定义的格式
g = _, _

[policy_effect]              
e = some(where (p.eft == allow))

[matchers]                 //满足条件
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
//请求用户与满足*.cvs p(策略)且满足g(组规则)且请求资源满足p(策略)规定资源  
```

[Go Web Iris中文网](https://www.studyiris.com/)