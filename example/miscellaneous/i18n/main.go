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