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