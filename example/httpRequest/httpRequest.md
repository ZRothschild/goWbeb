# `iris http request`请求数据处理
这一篇章主要讲述，`iris`框架服务如何处理`http request`请求
## 概述
1. 获取引用者(`extract referer`)
2. 前置自定义读取请求数据数据处理
3. 前置自定义读取请求数据数据处理(框架实现了接口)
4. `iris`自定义结构体映射获取`Form`表单请求数据(`read form`)
5. `iris`自定义结构体映射获取`json`格式请求数据(`read-json`)
6. `iris`自定义结构体映射获取`json`格式请求数据(`read-json`,并自动验证)
7. `iris`自定义结构体映射获取`xml`格式请求数据(`read xml`)
8. 请求日志(`request-logger`)
9. `iris`单文件上传(`upload-file`)
10. `iris`多文件上传(`upload-files`)