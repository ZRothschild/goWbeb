###Iris 框架的特性
* 专注于高性能
* 简单流畅的API   
* 高扩展性
* 强大的路由和中间件生态系统
  
  - 使用iris独特的表达主义路径解释器构建RESTful API
  - 动态路径参数化或通配符路由与静态路由不冲突
  - 使用重定向选项从URL中删除尾部斜杠
  - 使用虚拟主机和子域名变得容易
  - 分组API和静态或甚至动态子域名
  - net / http和negroni-like处理程序通过iris.FromStd兼容
  - 针对任意Http请求错误 定义处理函数
  - 支持事务和回滚
  - 支持响应缓存
  - 使用简单的函数嵌入资源并与go-bindata 保持兼容 
  - mvc
  
* 上下文
   - 高度可扩展的试图渲染(目前支持markdown,json,xml，jsonp等等)
   - 正文绑定器和发送HTTP响应的便捷功能
   - 限制请求正文
   - 提供静态资源或嵌入式资产
   - 本地化i18N 
   - 压缩（Gzip是内置的）
* 身份验证
   - Basic Authentication
   - OAuth, OAuth2 (支持27个以上的热门网站)
   - JWT
*服务器
   - 通过TLS提供服务时，自动安装和提供来自https://letsencrypt.org的证书
   - 默认为关闭状态
   - 在关闭，错误或中断事件时注册
   - 连接多个服务器，完全兼容 net/http#Server   
* 视图系统.支持五种模板隐隐 完全兼容 html/template   

* Websocket库，其API类似于socket.io [如果你愿意，你仍然可以使用你最喜欢的]

* 热重启
* Typescript集成 + Web IDE
* Iris是最具特色的网络框架之一，并非所有功能都在这里.
* 如果你发现我遗漏了一些东西,请发送邮箱到go-iris@qq.com