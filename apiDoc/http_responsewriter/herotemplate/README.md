# Hero Template Example

此文件夹包含原始英雄示例的`Iris`版本：[https://github.com/shiyanhui/hero/tree/master/examples/app](https://github.com/shiyanhui/hero/tree/master/examples/app)
`Iris`与`net/http` 100％兼容，因此您无需更改任何其他内容,除了原始示例中的处理程序输入。

从:
```go
if _, err := w.Write(buffer.Bytes()); err != nil {
// and
template.UserListToWriter(userList, w)
```
到: 
```go
if _, err := ctx.Write(buffer.Bytes()); err != nil {
// and
template.UserListToWriter(userList, ctx)
```
如此容易.

了解更多信息，请访问：[https://github.com/shiyanhui/hero](https://github.com/shiyanhui/hero)