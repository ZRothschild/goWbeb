# YAAG：golang web API 文档生成器

### YAAG介绍

Golang非常适合开发Web应用程序。 人们已经创建了许多优秀的Web框架，Web帮助程序库。 
如果我们考虑Golang中的整个Web应用程序生态系统，似乎很完整了，但我们却少了一个编写API文档的中间件。
因此，我们为基于Golang的Web应用程序创建了第一个API文档生成器。


大多数Web服务都将其API暴露给移动或第三方开发人员。 记录它们有点痛苦。 
我们正在努力减轻痛苦，至少对于您不必向世界公开您的文档的内部项目。 
YAAG生成简单的基于引导程序的API文档，无需编写任何注释。


YAAG是一个中间件。 您必须在路线中添加YAAG处理程序，您就完成了。 
继续使用POSTMAN，Curl或任何客户端调用您的API，YAAG将继续更新API Doc html。
注意：我们还生成一个包含所有API调用数据的json文件）