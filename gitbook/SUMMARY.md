# Summary

* [文档介绍](README.md)
* [IRIS](iris.md)
* [IRIS 验证与授权](auth.md)
    * [IRIS 基础验证与授权(Basic Auth)](auth/basic.md)
    * [IRIS Auth2 验证与授权](auth/oauth2.md)
* [IRIS 缓存使用](cache.md)
    * [IRIS 客户端缓存](cache/client.md)
    * [IRIS 简单缓存](cache/simple.md)
* [IRIS 配置文件类型](conf.md)
* [其他函数转换成 IRIS 处理函数](convertHandlers.md)
* [IRIS cookies 使用](cookies.md)
* [IRIS 试验性功能](exper.md)
    * [IRIS 实现基于角色的HTTP权限控制](exper/casbin.md)
    * [IRIS 简单的aws监听API](exper/cloudwatch.md)
    * [IRIS CORS 跨域资源共享](exper/cors.md)
    * [IRIS csrf 防御](exper/csrf.md)
    * [IRIS JWT 使用](exper/jwt.md)
    * [IRIS 限制HTTP请求次数的中间件](exper/limitHandler.md)
    * [IRIS 跟踪事务](exper/newrelic.md)
    * [IRIS 系统和服务监控系统](exper/prometheus.md)
    * [IRIS 域名安全加密](exper/secure.md)
* [IRIS 文件服务](fileServer.md)
    * [IRIS 基础文件服务](fileServer/Basic.md)
    * [IRIS 基础单页面文件服务](fileServer/basicSPA.md)
    * [IRIS embedded单页面文件服务](fileServer/embeddedSPA.md)
    * [IRIS embedded单页面文件服务与路由](fileServer/embeddedSPAOR.md)
    * [IRIS embedding打包静态文件](fileServer/embeddingFiles.md)
    * [IRIS embedding打包静态文件Gzip格式](fileServer/embeddingGziped.md)
    * [IRIS favicon应用图标](fileServer/favicon.md)
    * [IRIS 文件下载](fileServer/sendFiles.md)
* [IRIS hello world小示例](helloWorld.md)
* [IRIS hero依赖注入与结构体转化](hero.md)
    * [IRIS hero基础](hero/basic.md)
    * [IRIS hero MVC](hero/overview.md)
    * [IRIS hero sessions](hero/sessions.md)
    * [IRIS hero智能路由](hero/smartContract.md)
* [IRIS http 服务监听](httpListening.md)
    * [IRIS 自定义http服务](httpListen/customHttpserver.md)
    * [IRIS 自定义监听器](httpListen/customListener.md)
    * [IRIS 优雅关闭服务](httpListen/gracefulShutdown.md)
    * [IRIS 配置与`host`配置获取关闭](httpListen/irisConfAndHostConf.md)
    * [IRIS 地址服务器监听](httpListen/listenAddr.md)
    * [IRIS 服务监听与`letsencrypt`加密](httpListen/listenLetsencrypt.md)
    * [IRIS 服务监听启用安全传输协议](httpListen/listenTls.md)
    * [IRIS 使用.sock文件服务监听](httpListen/listenUnix.md)
    * [IRIS 使用chan通知关闭服务](httpListen/notifyOnShutdown.md)
* [IRIS http request请求数据处理](httpRequest.md)
    * [IRIS 前置自定义读取请求(实现实现Decode接口)](request/customViaUnmarshaler.md)
    * [IRIS 获取引用者(extract referer)](request/extractReferer.md)
    * [IRIS 请求数据验证](request/jsonStructValidation.md)
    * [IRIS 前置自定义读取请求数据数据](request/readCustomPerType.md)
    * [IRIS 自定义结构体映射获取Form表单请求数据](request/readForm.md)
    * [IRIS 自定义结构体映射获取json格式请求数据](request/readJson.md)
    * [IRIS 自定义结构体映射获取xml格式请求数据](request/readXml.md)
    * [IRIS 请求日志记录配置(控制台)](request/requestLogger.md)
    * [IRIS 请求日志记录(文件)](request/requestLoggerFile.md)
    * [IRIS 文件上传示例](request/uploadFile.md)
    * [IRIS 多文件上传示例](request/uploadFiles.md)
* [IRIS http response模板,返回数据类型格式](reponseWriter.md)
    * [IRIS hero模板引擎](responseWriter/herotemplate.md)
    * [IRIS quicktemplate模板引擎](responseWriter/quickTemplate.md)
    * [IRIS sse应用](responseWriter/sse.md)
    * [IRIS sse第三方应用](responseWriter/sseThirdParty.md)
    * [IRIS 流写入](responseWriter/streamWriter.md)
    * [IRIS 请求事务](responseWriter/transactions.md)
    * [IRIS 以gzip格式写入数据](responseWriter/writeGzip.md)
    * [IRIS 数据返回类型](responseWriter/writeRest.md)
* [IRIS 一些杂项例子](miscellaneous.md)
* [IRIS MVC示例](mvc.md)
    * [IRIS MVC基础示例](mvc/basic.md)
    * [IRIS MVC hello world示例](mvc/helloWorld.md)
    * [IRIS MVC login示例](mvc/login.md)
    * [IRIS MVC中如何使用中间件](mvc/middleware.md)
    * [IRIS MVC overview示例](mvc/overview.md)
    * [IRIS MVC前置中间件使用](mvc/perMethod.md)
    * [IRIS session控制器](mvc/sessionController.md)
    * [IRIS MVC 单例控制器](mvc/singleton.md)
    * [IRIS MVC websocket](mvc/websocket.md)
    * [IRIS MVC ExecutionRules实现中间件](mvc/withoutCtxNext.md)
* [对象关系映射](orm.md)
    * [IRIS xorm使用](orm/xorm.md)
* [IRIS overview示例](overview.md)
* [API文档生成器](yaag.md)
    * [IRIS API文档生成器示例](yaag/iris.md)