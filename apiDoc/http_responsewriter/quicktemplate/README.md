首先，安装[quicktemplate](https://github.com/valyala/quicktemplate) package and [quicktemplate compiler](https://github.com/valyala/quicktemplate/tree/master/qtc)

```sh
go get -u github.com/valyala/quicktemplate
go get -u github.com/valyala/quicktemplate/qtc
```

该示例已经为您编译了Go代码，因此:
```sh
go run main.go # http://localhost:8080
```

但是下面有一条说明，可以在https://github.com/valyala/quicktemplate找到完整的文档。

将模板文件保存到扩展名`*.qtpl`下的`templates`文件夹中，打开终端并在此文件夹中运行`qtc`。

如果一切顺利，`*.qtpl.go`文件必须出现在`templates`文件夹中。 这些文件包含所有`* .qtpl`文件的Go代码。

> 请记住，每次更改`/templates/*.qtpl`文件时，都必须运行`qtc`命令并重新构建应用程序。