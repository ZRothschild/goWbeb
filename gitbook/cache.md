# cache(缓存)
## 为什么要用缓存
- 请求更快：通过将内容缓存在本地浏览器或距离最近的缓存服务器（如`CDN`），在不影响网站交互的前提下可以大大加快网站加载速度
- 节省带宽：对于已缓存的文件，可以减少请求带宽甚至无需请求网络
- 降低服务器压力：在大量用户并发请求的情况下，服务器的性能受到限制，此时将一些静态资源放置在网络的多个节点，
可以起到均衡负载的作用，降低服务器的压力

## 缓存分类
1. 缓存分为服务端侧（`server side`，比如 `Nginx、Apache`）。服务端缓存又分为代理服务器缓存和反向代理服务器缓存
（也叫网关缓存，比如 `Nginx`反向代理、`nginx`缓存、`Squid`等）
2. 客户端侧（`client side`，比如`web browser`）。 例如 浏览器缓存