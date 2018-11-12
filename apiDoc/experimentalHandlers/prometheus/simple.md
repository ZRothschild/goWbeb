# prometheus
## prometheus介绍
访问[prometheus.io](prometheus.io)获取完整的文档，示例和指南。

`Prometheus`是一个云原生计算基础项目，是一个系统和服务监控系统。它以给定的时间间隔从配置的目标收集指标，评估规则表达式，显示结果，
并且如果观察到某些条件为真，则可以触发警报。

与其他监测系统相比，普罗米修斯的主要区别特征是：
1. 多维数据模型（由度量标准名称和键/值维度集定义的时间序列）
2. 灵活的查询语言，以利用此维度
3. 不依赖于分布式存储; 单个服务器节点是自治的
4. 时间序列集合通过`HTTP`上的拉模型进行
5. 通过中间网关支持推送时间序列
6. 通过服务发现或静态配置发现目标
7. 多种图形和仪表板支持模式
8. 支持分层和水平联合
## 目录结构
> 主目录`simple`

```html
    —— main.go
```
## 代码示例 
> `main.go`

```go
package main

import (
	"math/rand"
	"time"
	"github.com/kataras/iris"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	app := iris.New()
	m := prometheusMiddleware.New("serviceName", 300, 1200, 5000)
	app.Use(m.ServeHTTP)
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		//错误代码处理程序不与其他路由共享相同的中间件，所以单独执行错误
		m.ServeHTTP(ctx)
		ctx.Writef("Not Found")
	})
	app.Get("/", func(ctx iris.Context) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		ctx.Writef("Slept for %d milliseconds", sleep)
	})
	app.Get("/metrics", iris.FromStd(prometheus.Handler()))
	// http://localhost:8080/
	// http://localhost:8080/anotfound
	// http://localhost:8080/metrics
	app.Run(iris.Addr(":8080"))
}
```