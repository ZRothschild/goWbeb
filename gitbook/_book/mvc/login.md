# `MVC login`示例
## 目录结构
> 主目录`login`

```html
    —— datamodels
        —— user.go
    —— datasource
        —— users.go
    —— repositories
        —— user_repository.go
    —— services
        —— user_service.go
    —— web
        —— controllers
            —— user_controller.go
            —— users_controller.go
        —— middleware
            —— basicauth.go
        —— public
            —— css
                —— site.css
        —— viewmodels
            —— viewModels.md
        —— views
            —— shared
                —— error.html
                —— layout.html
            —— user
                —— login.html
                —— me.html
                —— register.html
   —— main.go  
```
## 代码示例
> 文件名称 `/datamodels/user.go`**User**结构体定义

```go
package datamodels

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)
//User是我们的用户示例模型。
//请注意标签（适用于我们的网络应用）
//应该保存在其他文件中，例如“web/viewmodels/user.go”
//可以通过嵌入datamodels.User或
//定义完全新的字段
//示例中，我们将使用此数据模型
//作为我们应用程序中唯一的一个用户模型。
type User struct {
	ID             int64     `json:"id" form:"id"`
	Firstname      string    `json:"firstname" form:"firstname"`
	Username       string    `json:"username" form:"username"`
	HashedPassword []byte    `json:"-" form:"-"`
	CreatedAt      time.Time `json:"created_at" form:"created_at"`
}
// IsValid可以做一些非常简单的“低级”数据验证
func (u User) IsValid() bool {
	return u.ID > 0
}
// GeneratePassword将根据我们为我们生成哈希密码
//用户的输入
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
// ValidatePassword将检查密码是否匹配
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
```
> 文件名称 `/datasource/user.go`数据资源，相当于数据库

```go
//文件: datasource/users.go
package datasource

import (
	"errors"
	"../datamodels"
)
//引擎来自何处获取数据，在这种情况下是用户。
type Engine uint32
const (
	//内存代表简单的内存位置;
	// map[int64]datamodels.User随时可以使用，这是我们在这个例子中的来源。
	Memory Engine = iota
	// Bolt for boltdb source location.
	Bolt
	// MySQL for mysql-compatible source location.
	MySQL
)
//为了简单起见，Load Users从内存中返回所有用户（空map）。
func LoadUsers(engine Engine) (map[int64]datamodels.User, error) {
	if engine != Memory {
		return nil, errors.New("for the shake of simplicity we're using a simple map as the data source")
	}
	return make(map[int64]datamodels.User), nil
}
```
> 文件名称 `/repositories/user_repository.go`对数据筛选，必要数据的仓库

```go
package repositories

import (
	"errors"
	"sync"
	"../datamodels"
)
// Query表示访问者和操作查询。
type Query func(datamodels.User) bool
// UserRepository处理用户实体/模型的基本操作。
//它是一个可测试的接口，即内存用户存储库或 连接到sql数据库。
type UserRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (user datamodels.User, found bool)
	SelectMany(query Query, limit int) (results []datamodels.User)
	InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
	Delete(query Query, limit int) (deleted bool)
}
// NewUserRepository返回一个新的基于用户内存的存储库，
//我们示例中唯一的存储库类型。
func NewUserRepository(source map[int64]datamodels.User) UserRepository {
	return &userMemoryRepository{source: source}
}
//userMemoryRepository是一个“UserRepository”
//使用内存数据源（map）管理用户。
type userMemoryRepository struct {
	source map[int64]datamodels.User
	mu     sync.RWMutex
}
const (
	// ReadOnlyMode将RLock(read) 数据。
	ReadOnlyMode = iota
	// ReadWriteMode将锁定(read/write)数据。
	ReadWriteMode
)
func (r *userMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0
	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}
	for _, user := range r.source {
		ok = query(user)
		if ok {
			if action(user) {
				loops++
				if actionLimit >= loops {
					break // break
				}
			}
		}
	}
	return
}
//Select接收查询方法
//为内部的每个用户模型触发查找我们想象中的数据源
//当该函数返回true时，它会停止迭代。
//它实际上是一个简单但非常游泳的原型函数
//自从我第一次想到它以来，我一直在使用它，
//希望你会发现它也很有用。
func (r *userMemoryRepository) Select(query Query) (user datamodels.User, found bool) {
	found = r.Exec(query, func(m datamodels.User) bool {
		user = m
		return true
	}, 1, ReadOnlyMode)
	//设置一个空的datamodels.User，如果根本找不到的话
	if !found {
		user = datamodels.User{}
	}
	return
}
// SelectMany与Select相同但返回一个或多个datamodels.User作为切片。
//如果limit <= 0则返回所有内容。
func (r *userMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.User) {
	r.Exec(query, func(m datamodels.User) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)
	return
}
// InsertOrUpdate将用户添加或更新到（内存）存储。
//返回新用户，如果有则返回错误。
func (r *userMemoryRepository) InsertOrUpdate(user datamodels.User) (datamodels.User, error) {
	id := user.ID
	if id == 0 {
		var lastID int64
		//找到最大的ID，以便不重复 在制作应用中，
		//您可以使用第三方库以生成UUID作为字符串。
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()
		id = lastID + 1
		user.ID = id
		r.mu.Lock()
		r.source[id] = user
		r.mu.Unlock()
		return user, nil
	}
	//基于user.ID更新操作，
	//这里我们将允许更新海报和流派，如果不是空的话。
	//或者我们可以做替换：
	// r.source [id] =user
	//的代码;
	current, exists := r.Select(func(m datamodels.User) bool {
		return m.ID == id
	})
	if !exists { // ID不是真实的，返回错误。
		return datamodels.User{}, errors.New("failed to update a nonexistent user")
	}
	//和r.source [id] =user 进行纯替换
	if user.Username != "" {
		current.Username = user.Username
	}
	if user.Firstname != "" {
		current.Firstname = user.Firstname
	}
	r.mu.Lock()
	r.source[id] = current
	r.mu.Unlock()
	return user, nil
}
func (r *userMemoryRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(m datamodels.User) bool {
		delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)
}
```
> 文件名称 `/services/user_service.go`业务逻辑代码

```go
package services

import (
	"errors"
	"../datamodels"
	"../repositories"
)
//UserService处理用户数据模型的CRUID操作，
//它取决于用户存储库的操作。
//这是将数据源与更高级别的组件分离。
//因此，不同的存储库类型可以使用相同的逻辑，而无需任何更改。
//它是一个接口，它在任何地方都被用作接口
//因为我们可能需要在将来更改或尝试实验性的不同域逻辑。
type UserService interface {
	GetAll() []datamodels.User
	GetByID(id int64) (datamodels.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	DeleteByID(id int64) bool
	Update(id int64, user datamodels.User) (datamodels.User, error)
	UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	UpdateUsername(id int64, newUsername string) (datamodels.User, error)
	Create(userPassword string, user datamodels.User) (datamodels.User, error)
}
// NewUserService返回默认用户服务
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
type userService struct {
	repo repositories.UserRepository
}
// GetAll返回所有用户。
func (s *userService) GetAll() []datamodels.User {
	return s.repo.SelectMany(func(_ datamodels.User) bool {
		return true
	}, -1)
}
// GetByID根据其id返回用户。
func (s *userService) GetByID(id int64) (datamodels.User, bool) {
	return s.repo.Select(func(m datamodels.User) bool {
		return m.ID == id
	})
}
//获取yUsernameAndPassword根据用户名和密码返回用户，
//用于身份验证。
func (s *userService) GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool) {
	if username == "" || userPassword == "" {
		return datamodels.User{}, false
	}
	return s.repo.Select(func(m datamodels.User) bool {
		if m.Username == username {
			hashed := m.HashedPassword
			if ok, _ := datamodels.ValidatePassword(userPassword, hashed); ok {
				return true
			}
		}
		return false
	})
}
//更新现有用户的每个字段的更新，
//通过公共API使用是不安全的
//但是我们将在web  controllers/user_controller.go#PutBy上使用它
//为了向您展示它是如何工作的。
func (s *userService) Update(id int64, user datamodels.User) (datamodels.User, error) {
	user.ID = id
	return s.repo.InsertOrUpdate(user)
}
// UpdatePassword更新用户的密码。
func (s *userService) UpdatePassword(id int64, newPassword string) (datamodels.User, error) {
	//更新用户并将其返回。
	hashed, err := datamodels.GeneratePassword(newPassword)
	if err != nil {
		return datamodels.User{}, err
	}
	return s.Update(id, datamodels.User{
		HashedPassword: hashed,
	})
}
// UpdateUsername更新用户的用户名
func (s *userService) UpdateUsername(id int64, newUsername string) (datamodels.User, error) {
	return s.Update(id, datamodels.User{
		Username: newUsername,
	})
}
//创建插入新用户，
// userPassword是客户端类型的密码
//它将在插入我们的存储库之前进行哈希处理
func (s *userService) Create(userPassword string, user datamodels.User) (datamodels.User, error) {
	if user.ID > 0 || userPassword == "" || user.Firstname == "" || user.Username == "" {
		return datamodels.User{}, errors.New("unable to create this user")
	}
	hashed, err := datamodels.GeneratePassword(userPassword)
	if err != nil {
		return datamodels.User{}, err
	}
	user.HashedPassword = hashed
	return s.repo.InsertOrUpdate(user)
}
// DeleteByID按其id删除用户。
//如果删除则返回true，否则返回false。
func (s *userService) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m datamodels.User) bool {
		return m.ID == id
	}, 1)
}
```
> 文件名称 `/web/controllers/user_controller.go`

```go
// 文件: controllers/user_controller.go
package controllers

import (
	"../../datamodels"
	"../../services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)
// UserController是我们的/用户控制器。
// UserController负责处理以下请求：
// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// GET 				/user/me
//所有HTTP方法 /user/logout
type UserController struct {
	//每个请求都由Iris自动绑定上下文，
	//记住，每次传入请求时，iris每次都会创建一个新的UserController，
	//所以所有字段都是默认的请求范围，只能设置依赖注入
	//自定义字段，如服务，对所有请求都是相同的（静态绑定）
	//和依赖于当前上下文的会话（动态绑定）。
	Ctx iris.Context
	//我们的UserService，它是一个接口
	//从主应用程序绑定。
	Service services.UserService
	//Session，使用来自main.go的依赖注入绑定
	Session *sessions.Session
}
const userIDKey = "UserID"
func (c *UserController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}
func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}
func (c *UserController) logout() {
	c.Session.Destroy()
}
var registerStaticView = mvc.View{
	Name: "user/register.html",
	Data: iris.Map{"Title": "User Registration"},
}
// GetRegister 处理 GET: http://localhost:8080/user/register.
func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}
	return registerStaticView
}
// PostRegister 处理 POST: http://localhost:8080/user/register.
func (c *UserController) PostRegister() mvc.Result {
	//从表单中获取名字，用户名和密码
	var (
		firstname = c.Ctx.FormValue("firstname")
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)
	//创建新用户，密码将由服务进行哈希处理
	u, err := c.Service.Create(password, datamodels.User{
		Username:  username,
		Firstname: firstname,
	})
	//将用户的id设置为此会话，即使err！= nil，
	//零id无关紧要因为.getCurrentUserID()检查它。
	//如果错误！= nil那么它将被显示，见下面的mvc.Response.Err：err
	c.Session.Set(userIDKey, u.ID)
	return mvc.Response{
		//如果不是nil，则会显示此错误
		Err: err,
		//从定向 /user/me.
		Path: "/user/me",
		//当从POST重定向到GET请求时，您应该使用此HTTP状态代码，
		//但是如果你有一些（复杂的）选择
		//在线搜索甚至是HTTP RFC。
		//状态“查看其他”RFC 7231，但虹膜可以自动修复它
		//但很高兴知道你可以设置自定义代码;
		//代码：303，
	}
}
var loginStaticView = mvc.View{
	Name: "user/login.html",
	Data: iris.Map{"Title": "User Login"},
}
// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.logout()
	}
	return loginStaticView
}
// PostLogin handles
// PostLogin处理POST: http://localhost:8080/user/register.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)
	u, found := c.Service.GetByUsernameAndPassword(username, password)
	if !found {
		return mvc.Response{
			Path: "/user/register",
		}
	}
	c.Session.Set(userIDKey, u.ID)
	return mvc.Response{
		Path: "/user/me",
	}
}
// GetMe 处理P GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		//如果没有登录，则将用户重定向到登录页面。
		return mvc.Response{Path: "/user/login"}
	}
	u, found := c.Service.GetByID(c.getCurrentUserID())
	if !found {
		//如果session存在但由于某种原因用户不存在于“数据库”中
		//然后注销并重新执行该函数，它会将客户端重定向到
		// /user/login页面。
		c.logout()
		return c.GetMe()
	}
	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User":  u,
		},
	}
}
// AnyLogout处理 All/AnyHTTP 方法：http://localhost:8080/user/logout
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}
	c.Ctx.Redirect("/user/login")
}
```
> 文件名称 `/web/controllers/users_controller.go`

```go
package controllers

import (
	"../../datamodels"
	"../../services"
	"github.com/kataras/iris"
)
// UsersController是我们 /users API控制器。
// GET				/users  | get all
// GET				/users/{id:long} | 获取通过 id
// PUT				/users/{id:long} | 修改通过 id
// DELETE			/users/{id:long} | 删除通过 id
//需要基本身份验证。
type UsersController struct {
	//可选：Iris在每个请求中自动绑定上下文，
	//记住，每次传入请求时，iris每次都会创建一个新的UserController，
	//所以所有字段都是默认的请求范围，只能设置依赖注入
	//自定义字段，如Service，对所有请求都是相同的（静态绑定）。
	Ctx iris.Context
	//我们的UserService，它是一个接口
	//从主应用程序绑定。
	Service services.UserService
}
//获取用户的返回列表。
// 示例:
// curl -i -u admin:password http://localhost:8080/users
// func (c *UsersController) Get() (results []viewmodels.User) {
// 	data := c.Service.GetAll()
//
// 	for _, user := range data {
// 		results = append(results, viewmodels.User{user})
// 	}
// 	return
// }
//否则只返回数据模型
func (c *UsersController) Get() (results []datamodels.User) {
	return c.Service.GetAll()
}
// GetBy 返回指定一个id用户
// 示例:
// curl -i -u admin:password http://localhost:8080/users/1
func (c *UsersController) GetBy(id int64) (user datamodels.User, found bool) {
	u, found := c.Service.GetByID(id)
	if !found {
		//此信息将被绑定到
		// main.go -> app.OnAnyErrorCode -> NotFound -> shared/error.html -> .Message text.
		c.Ctx.Values().Set("message", "User couldn't be found!")
	}
	return u, found // it will throw/emit 404 if found == false.
}
// PutBy 修改指定用户.
// 示例:
// curl -i -X PUT -u admin:password -F "username=kataras"
// -F "password=rawPasswordIsNotSafeIfOrNotHTTPs_You_Should_Use_A_client_side_lib_for_hash_as_well"
// http://localhost:8080/users/1
func (c *UsersController) PutBy(id int64) (datamodels.User, error) {
	// username := c.Ctx.FormValue("username")
	// password := c.Ctx.FormValue("password")
	u := datamodels.User{}
	if err := c.Ctx.ReadForm(&u); err != nil {
		return u, err
	}
	return c.Service.Update(id, u)
}
// DeleteBy 删除指定用户
// Demo:
// curl -i -X DELETE -u admin:password http://localhost:8080/users/1
func (c *UsersController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)
	if wasDel {
		//返回已删除用户的ID
		return map[string]interface{}{"deleted": id}
	}
	//在这里我们可以看到一个方法函数
	//可以返回这两种类型中的任何一种（map或int），
	//我们不必将返回类型指定为特定类型。
	return iris.StatusBadRequest //等同于 400.
}
```
> 文件名称 `/web/middleware/basicauth.go`

```go
// 文件: middleware/basicauth.go
package middleware

import "github.com/kataras/iris/middleware/basicauth"

// BasicAuth中间件示例。
var BasicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{
		"admin": "password",
	},
})
```
> 文件名称 `/web/public/css/site.css`

```css
/* Bordered form */
form {
    border: 3px solid #f1f1f1;
}
/* Full-width inputs */
input[type=text], input[type=password] {
    width: 100%;
    padding: 12px 20px;
    margin: 8px 0;
    display: inline-block;
    border: 1px solid #ccc;
    box-sizing: border-box;
}
/* Set a style for all buttons */
button {
    background-color: #4CAF50;
    color: white;
    padding: 14px 20px;
    margin: 8px 0;
    border: none;
    cursor: pointer;
    width: 100%;
}
/* Add a hover effect for buttons */
button:hover {
    opacity: 0.8;
}
/* Extra style for the cancel button (red) */
.cancelbtn {
    width: auto;
    padding: 10px 18px;
    background-color: #f44336;
}
/* Center the container */
/* Add padding to containers */
.container {
    padding: 16px;
}
/* The "Forgot password" text */
span.psw {
    float: right;
    padding-top: 16px;
}
/* Change styles for span and cancel button on extra small screens */
@media screen and (max-width: 300px) {
    span.psw {
        display: block;
        float: none;
    }
    .cancelbtn {
        width: 100%;
    }
}
```
> 文件名称 `/web/viewmodels/viewModels.md`

```markdown
应该有视图模型，客户端将能够看到的结构例：
import (
    "github.com/kataras/iris/_examples/mvc/login/datamodels"
    "github.com/kataras/iris/context"
)

type User struct {
    datamodels.User
}
func (m User) IsValid() bool {
    /*做一些检查，如果有效则返回true ... */
    return m.ID > 0
}

Iris能够将任何自定义数据结构转换为HTTP响应调度程序，
从理论上讲，如果真的有必要，可以使用以下内容;

//Dispatch实现`kataras/iris/mvc＃Result`接口。
//将`User` 作为受控的http响应发送。
//如果其ID为零或更小，则返回404未找到错误
//否则返回其json表示，
//（就像控制器的函数默认为自定义类型一样）。
//不要过度，应用程序的逻辑不应该在这里。可以在这里添加简单的检查验证
//这只是一个展示
//想象一下设计更大的应用程序时此功能將很有帮助。
//调用控制器方法返回值的函数
//是`User` 的类型。
func (m User) Dispatch(ctx context.Context) {
    if !m.IsValid() {
        ctx.NotFound()
        return
    }
    ctx.JSON(m, context.JSON{Indent: " "})
}
但是，我们将使用“datamodels”作为唯一的一个模型包，因为
User结构不包含任何敏感数据，客户端可以查看其所有字段我们内部不需要任何额外的功能或验证。
```
> 文件名称 `/web/views/shared/error.html`

```html
{% raw %}
<h1>错误</h1>
<h2>处理您的请求时发生错误。</h2>
<h3>{{.Message}}</h3>
<footer>
    <h2>Sitemap</h2>
    <a href="http://localhost:8080/user/register">/user/register</a><br/>
    <a href="http://localhost:8080/user/login">/user/login</a><br/>
    <a href="http://localhost:8080/user/logout">/user/logout</a><br/>
    <a href="http://localhost:8080/user/me">/user/me</a><br/>
    <h3>requires authentication</h3><br/>
    <a href="http://localhost:8080/users">/users</a><br/>
    <a href="http://localhost:8080/users/1">/users/{id}</a><br/>
</footer>
{% endraw %}
```
> 文件名称 `shared/layout.html`

```html
{% raw %}
<html>
<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" type="text/css" href="/public/css/site.css" />
</head>
<body>
    {{ yield }}
</body>
</html>
{% endraw %}
```
> 文件名称 `/web/views/user/login.html`

```html
<form action="/user/login" method="POST">
    <div class="container">
        <label><b>Username</b></label>
        <input type="text" placeholder="Enter Username" name="username" required>
        <label><b>Password</b></label>
        <input type="password" placeholder="Enter Password" name="password" required>
        <button type="submit">Login</button>
    </div>
</form>
```
> 文件名称 `/web/views/user/me.html`

```html
{% raw %}
<p>
    Welcome back <strong>{{.User.Firstname}}</strong>!
</p>
{% endraw %}
```
> 文件名称 `/web/views/user/register.html`

```html
<form action="/user/register" method="POST">
    <div class="container">
        <label><b>Firstname</b></label>
        <input type="text" placeholder="Enter Firstname" name="firstname" required>
        <label><b>Username</b></label>
        <input type="text" placeholder="Enter Username" name="username" required>
        <label><b>Password</b></label>
        <input type="password" placeholder="Enter Password" name="password" required>
        <button type="submit">Register</button>
    </div>
</form>
```
![目录结构](../img/folder_structure.png)