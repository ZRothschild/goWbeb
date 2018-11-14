# `iris`自定义结构体映射获取`json`格式请求数据并且自动验证
## 目录结构
> 主目录`readJsonStructValidation`

```html
    —— main.go
```
## 代码示例
> `main.go`

```go
//包main显示验证器（最新版本9）与Iris的集成。
//您可以在以下网址找到更多这样的示例：https//github.com/go-playground/validator/blob/v9/_examples
package main

import (
	"fmt"
	"github.com/kataras/iris"
	// $ go get gopkg.in/go-playground/validator.v9
	"gopkg.in/go-playground/validator.v9"
	//"gopkg.in/go-playground/validator.v8"
)
//用户包含用户信息。
type User struct {
	FirstName      string     `json:"fname"`
	LastName       string     `json:"lname"`
	Age            uint8      `json:"age" validate:"gte=0,lte=130"`
	Email          string     `json:"email" validate:"required,email"`
	FavouriteColor string     `json:"favColor" validate:"hexcolor|rgb|rgba"`
	Addresses      []*Address `json:"addresses" validate:"required,dive,required"`
}
//地址包含用户地址信息
type Address struct {
	Street string `json:"street" validate:"required"`
	City   string `json:"city" validate:"required"`
	Planet string `json:"planet" validate:"required"`
	Phone  string `json:"phone" validate:"required"`
}

// Use a single instance of Validate, it caches struct info.
//使用Validate的单个实例，它缓存struct info。
var validate *validator.Validate

func main() {
	validate = validator.New()
	//为'用户'注册验证
	//注意：只需要为'User'，validator注册一个非指针类型
	//在类型检查期间内部取消引用。
	validate.RegisterStructValidation(UserStructLevelValidation, User{})
	app := iris.New()
	app.Get("/user", func(ctx iris.Context) {
		fmt.Println("test")
		var user User
		if err := ctx.ReadJSON(&user); err != nil {
			// 处理错误
			fmt.Println(err)
		}
		//为错误的验证输入返回InvalidValidationError，nil或ValidationErrors（[] FieldError）
		err := validate.Struct(user)
		if err != nil {
			//只有在您的代码可以生成时才需要进行此检查
			//验证的无效值，例如与nil的接口
			//大多数包括我自己的值通常不会有这样的代码。
			if _, ok := err.(*validator.InvalidValidationError); ok {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return
			}
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.StructNamespace()) //注册自定义TagNameFunc时可能会有所不同
				fmt.Println(err.StructField()) //通过将alt名称传递给ReportError，如下所示
				fmt.Println(err.Tag())
				fmt.Println(err.ActualTag())
				fmt.Println(err.Kind())
				fmt.Println(err.Type())
				fmt.Println(err.Value())
				fmt.Println(err.Param())
				fmt.Println()
				//或者将它们收集为json对象
				//并通过ctx.JSON将收集的错误发送回客户端
				// {
				// 	"namespace":        err.Namespace(),
				// 	"field":            err.Field(),
				// 	"struct_namespace": err.StructNamespace(),
				// 	"struct_field":     err.StructField(),
				// 	"tag":              err.Tag(),
				// 	"actual_tag":       err.ActualTag(),
				// 	"kind":             err.Kind().String(),
				// 	"type":             err.Type().String(),
				// 	"value":            fmt.Sprintf("%v", err.Value()),
				// 	"param":            err.Param(),
				// }
			}
			ctx.WriteString("testasdd")
			//从这里你可以用你想要的任何语言创建自己的错误信息。
			return
		}
		//将用户保存到数据库
	})
	//使用Postman或其他什么来做POST请求
	//使用RAW BODY的http//localhost:8080/user
	/*
		{
			"fname": "",
			"lname": "",
			"age": 45,
			"email": "mail@example.com",
			"favColor": "#000",
			"addresses": [{
				"street": "Eavesdown Docks",
				"planet": "Persphone",
				"phone": "none",
				"city": "Unknown"
			}]
		}
	*/
	//Content-Type to application/json（可选，如果选择更好）。
	//由于空的`User.FirstName`（json中的fname），此请求将失败
	//和`User.LastName`（json中的lname）。
	//检查iris应用程序终端输出
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
// UserStructLevelValidation包含并非总是的自定义结构级别验证
//在字段验证级别有同样有意义。例如，此函数验证了这一点
//存在FirstName或LastName; 可以通过自定义字段验证完成，但随后
//必须将它添加到复制逻辑+开销的两个字段中，这样就可以了,仅验证一次。
//
//注意：您可能会问为什么我不能在验证器之外执行此操作，因为这样做
//挂钩到验证器，你可以与验证标签结合，但仍然有常见错误输出格式。
func UserStructLevelValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(User)
	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
		sl.ReportError(user.FirstName, "FirstName", "fname", "fnameorlname", "")
		sl.ReportError(user.LastName, "LastName", "lname", "fnameorlname", "")
	}
	//加上可以更多，甚至使用不同于“fnameorlname”的标签。
}
```