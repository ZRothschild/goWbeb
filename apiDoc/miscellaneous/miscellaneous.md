# 一些杂项例子(`miscellaneous`)
## 文件日志记录(`file logger`)
### 代码示例
> 文件名称`main.go`
```go
package main

import (
	"os"
	"time"
	"github.com/kataras/iris"
)

//根据日期获取文件名，仅用于容易区分
func todayFilename() string {
	today := time.Now().Format("Jan 02 2006")
	return today + ".txt"
}

func newLogFile() *os.File {
	filename := todayFilename()
	//打开文件，如果服务器重新启动，这将附加到今天的文件。
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	f := newLogFile()
	defer f.Close()
	app := iris.New()
	//将文件作为记录器附加，请记住，iris的app logger只是一个io.Writer。
	//如果需要同时将日志写入文件和控制台，请使用以下代码。
	// app.Logger().SetOutput(io.MultiWriter(f, os.Stdout))
	app.Logger().SetOutput(f)
	app.Get("/ping", func(ctx iris.Context) {
		// for the sake of simplicity, in order see the logs at the ./_today_.txt
		ctx.Application().Logger().Infof("Request path: %s", ctx.Path())
		ctx.WriteString("pong")
	})
	//访问http://localhost:8080/ping
	//打开文件 ./logs{TODAY}.txt file.
	if err := app.Run(iris.Addr(":8080"), iris.WithoutBanner, iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
		app.Logger().Warn("Shutdown with error: " + err.Error())
	}
}
```
## 多语言切换(`i18n`)
### 代码示例
> 文件名称`main.go`
```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
)

func newApp() *iris.Application {
	app := iris.New()
	globalLocale := i18n.New(i18n.Config{
		Default:      "en-US",
		URLParameter: "lang",
		Languages: map[string]string{
			"en-US": "./locales/locale_en-US.ini",
			"el-GR": "./locales/locale_el-GR.ini",
			"zh-CN": "./locales/locale_zh-CN.ini"}})
	app.Use(globalLocale)
	app.Get("/", func(ctx iris.Context) {
		//它试图通过以下方式找到当前语言环境值：
		// ctx.Values().GetString("language")
		//如果是空的那么它尝试从配置中设置的URLParameter中查找
		//如果没找到的话
		//它试图通过("language")cookie找到语言环境值
		//如果没有找到，则将其设置为配置上设置的默认值

		// hi是键，'iris'是.ini文件中的%s,可以叫参数
		//第二个参数是可选的

		// hi := ctx.Translate("hi", "iris")
		// 或:
		hi := i18n.Translate(ctx, "hi", "iris")

		language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
		//返回'en-US'的形式
		//找到的第一个默认语言保存在名称为("language")的cookie中，
		//你可以通过改变：iris.TranslateLanguageContextKey的值来改变它
		ctx.Writef("From the language %s translated output: %s", language, hi)
	})

	multiLocale := i18n.New(i18n.Config{
		Default:      "en-US",
		URLParameter: "lang",
		Languages: map[string]string{
			"en-US": "./locales/locale_multi_first_en-US.ini, ./locales/locale_multi_second_en-US.ini",
			"el-GR": "./locales/locale_multi_first_el-GR.ini, ./locales/locale_multi_second_el-GR.ini"}})
	app.Get("/multi", multiLocale, func(ctx iris.Context) {
		language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())

		fromFirstFileValue := i18n.Translate(ctx, "key1")
		fromSecondFileValue := i18n.Translate(ctx, "key2")
		ctx.Writef("From the language: %s, translated output:\n%s=%s\n%s=%s",
			language, "key1", fromFirstFileValue,
			"key2", fromSecondFileValue)
	})
	return app
}

func main() {
	app := newApp()
	// 访问 http://localhost:8080/?lang=el-GR
	// 或 http://localhost:8080 (default is en-US)
	// 或 http://localhost:8080/?lang=zh-CN
	// 访问 http://localhost:8080/multi?lang=el-GR
	// 或 http://localhost:8080/multi (default is en-US)
	// 或 http://localhost:8080/multi?lang=en-US
	// 或使用cookie来设置语言.
	app.Run(iris.Addr(":8080"))
}
```
> 文件名称`main_test.go`(测试代码)
```go
package main

import (
	"fmt"
	"testing"
	"github.com/kataras/iris/httptest"
)

func TestI18n(t *testing.T) {
	app := newApp()
	expectedf := "From the language %s translated output: %s"
	var (
		elgr = fmt.Sprintf(expectedf, "el-GR", "γεια, iris")
		enus = fmt.Sprintf(expectedf, "en-US", "hello, iris")
		zhcn = fmt.Sprintf(expectedf, "zh-CN", "您好，iris")
		elgrMulti = fmt.Sprintf("From the language: %s, translated output:\n%s=%s\n%s=%s", "el-GR",
			"key1",
			"αυτό είναι μια τιμή από το πρώτο αρχείο: locale_multi_first",
			"key2",
			"αυτό είναι μια τιμή από το δεύτερο αρχείο μετάφρασης: locale_multi_second")
		enusMulti = fmt.Sprintf("From the language: %s, translated output:\n%s=%s\n%s=%s", "en-US",
			"key1",
			"this is a value from the first file: locale_multi_first",
			"key2",
			"this is a value from the second file: locale_multi_second")
	)
	e := httptest.New(t, app)
	// default is en-US
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal(enus)
	// default is en-US if lang query unable to be found
	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal(enus)
	e.GET("/").WithQueryString("lang=el-GR").Expect().Status(httptest.StatusOK).
		Body().Equal(elgr)
	e.GET("/").WithQueryString("lang=en-US").Expect().Status(httptest.StatusOK).
		Body().Equal(enus)
	e.GET("/").WithQueryString("lang=zh-CN").Expect().Status(httptest.StatusOK).
		Body().Equal(zhcn)
	e.GET("/multi").WithQueryString("lang=el-GR").Expect().Status(httptest.StatusOK).
		Body().Equal(elgrMulti)
	e.GET("/multi").WithQueryString("lang=en-US").Expect().Status(httptest.StatusOK).
		Body().Equal(enusMulti)
}
```
> 文件名称`locales/locale_el-GR.ini`(语言包)
```ini
hi = γεια, %s
```
> 文件名称`locales/locale_en-US.ini`(语言包)
```ini
hi = hello, %s
```
> 文件名称`locales/locale_zh-CN.ini`(语言包)
```ini
hi = 您好，%s
```
### 目录结构
> 主目录`i18n`
```html
    —— locales
        —— locale_el-GR.ini
        —— locale_en-US.ini
        —— locale_zh-CN.ini
    —— main.go
    —— main_test.go
```
## 代码性能测试(`pprof`)
### 代码示例
> 文件名称`main.go`
```go
//go中有pprof包来做代码的性能监控
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/pprof"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1> Please click <a href='/debug/pprof'>here</a>")
	})
	app.Any("/debug/pprof/{action:path}", pprof.New())
	app.Run(iris.Addr(":8080"))
}
```
## 图形验证码(`recaptcha`)
### 代码示例
> 文件名称`main.go`
```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recaptcha"
)

//密钥应通过https://www.google.com/recaptcha获取
const (
	recaptchaPublic = ""
	recaptchaSecret = ""
)

func showRecaptchaForm(ctx iris.Context, path string) {
	ctx.HTML(recaptcha.GetFormHTML(recaptchaPublic, path))
}

func main() {
	app := iris.New()
	// On both Get and Post on this example, so you can easly
	// use a single route to show a form and the main subject if recaptcha's validation result succeed.
	//在此示例的Get和Post上，您可以轻松
	//使用单个路径显示表单，并且重新验证结果的主要主题成功。
	app.HandleMany("GET POST", "/", func(ctx iris.Context) {
		if ctx.Method() == iris.MethodGet {
			showRecaptchaForm(ctx, "/")
			return
		}
		result := recaptcha.SiteFerify(ctx, recaptchaSecret)
		if !result.Success {
			/* 如果你想要或什么都不做，重定向到这里 */
			ctx.HTML("<b> failed please try again </b>")
			return
		}
		ctx.Writef("succeed.")
	})
	app.Run(iris.Addr(":8080"))
}
```
> 文件名称`custom_form/main.go`
```go
package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recaptcha"
)
//密钥应通过https://www.google.com/recaptcha获取
const (
	recaptchaPublic = "6Lf3WywUAAAAAKNfAm5DP2J5ahqedtZdHTYaKkJ6"
	recaptchaSecret = "6Lf3WywUAAAAAJpArb8nW_LCL_PuPuokmEABFfgw"
)

func main() {
	app := iris.New()
	r := recaptcha.New(recaptchaSecret)
	app.Get("/comment", showRecaptchaForm)
	//在主处理程序之前传递中间件或使用`recaptcha.SiteVerify`。
	app.Post("/comment", r, postComment)
	app.Run(iris.Addr(":8080"))
}

var htmlForm = `<form action="/comment" method="POST">
	    <script src="https://www.google.com/recaptcha/api.js"></script>
		<div class="g-recaptcha" data-sitekey="%s"></div>
    	<input type="submit" name="button" value="Verify">
</form>`

func showRecaptchaForm(ctx iris.Context) {
	contents := fmt.Sprintf(htmlForm, recaptchaPublic)
	ctx.HTML(contents)
}

func postComment(ctx iris.Context) {
	// [...]
	ctx.JSON(iris.Map{"success": true})
}
```
### 目录结构
> 主目录`recaptcha`
```html
    —— custom_form
        —— main.go
    —— main.go
```
## 异常回复(`recover`)
### 代码示例
> 文件名称`main.go`
```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	i := 0
	//让我们在下一个请求时模拟panic
	app.Get("/", func(ctx iris.Context) {
		i++
		if i%2 == 0 {
			panic("a panic here")
		}
		ctx.Writef("Hello, refresh one time more to get panic!")
	})
	// http://localhost:8080, 刷新5-6次
	app.Run(iris.Addr(":8080"))
}
//注意： app := iris.Default()而不是iris.New()自动使用恢复中间件。
```