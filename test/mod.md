# 不需要翻墙之`Go Modules`模块版本控制介绍 `Go goproxy`代理模块镜像设置

## `Go goproxy` 代理设置

> 好多go语言爱好者，初次体验go语言就被下载依赖包难住了。执行`go get github.com/kataras/iris/v12@latest`然后等了许久，最终下载请求超时。
> 很多包工具都是需要翻墙才能获取到的，然后你又不能科学上网，只能去GitHub找找有没有对应的包，然后手动下载。

### 为了解决不用翻墙就能下载我们需要的go包我们可以设置代理

> 下面是一些代理提供企业或机构

1. proxy.golang.org - 官方项目-由Google运行-Go团队构建的默认Go模块代理
2. gocenter.io - 免费社区项目-由JFrog Artifactory运行-主要的Go模块存储库
3. mirrors.tencent.com/go - 商业项目-由腾讯云运营-Go模块备用代理
4. mirrors.aliyun.com/goproxy -商业项目-由阿里云运营-Go模块备用代理
5. goproxy.cn - 开源项目-由七牛云运营-中国最受信任的Go模块代理
6. goproxy.io - 开源项目-由中国Golang贡献者俱乐部运营-Go模块的全球代理
7. Athens - 开源项目-自托管-Go模块数据存储和代理
8. athens.azurefd.net - 开源项目-由Microsoft运营-运行Athens的托管模块代理
9. Goproxy - 开源项目-自托管-极简的Go模块代理处理程序
10.THUMBAI - 开源项目-自托管-Go mod代理服务器和Go vanity导入路径服务器

### go goproxy 设置代理

1. window 操作

```bash
# go 代理仓库
$ go env -w GOPROXY=https://goproxy.cn,direct 

# go 私有仓库 用逗号分隔
$ go env -w GOPRIVATE=*.corp.com,github.com/secret/repo

```

2. mac 与 Linux 操作

```bash

# go 代理仓库
$ export GOPROXY=https://goproxy.cn

# go 私有仓库 用逗号分隔 如果需要使用自己或公司的私有代码库，则需要使用GOPRIVATE配置私有库或项目的地址
$ export GOPRIVATE=*.corp.com,github.com/secret/repo

```

> 其实就是设置环境变化，设置好你就不需要翻墙了

## `Go Modules`模块版本控制

> 我这里只是简单介绍而已如果要深入了解请移步[这里](https://github.com/golang/go/wiki/Modules)

### 开启模块版本控制

> go 1.14说明

1. window 开启模块版本控制可用 `go env -w GO111MODULE=on` mac|linux 开启模块版本控制可用 `export GO111MODULE=on`
2. 创建你的项目目录，进入你创建的项目目录 初始化自己的模块 `go mod init [目录名/项目名]`，会新生成go.mod记住它的数据
```bash
$ mkdir -p /tmp/scratchpad/repo
$ cd /tmp/scratchpad/repo

$ go mod init github.com/my/repo

go: creating new go.mod: module github.com/my/repo

```

3. 下载你需要的依赖包 `go get github.com/kataras/iris/v12@latest` 然后看看你的go.mod文件，和最开始的对比一下


## 其余的一些介绍

### GO111MODULE 有三个值可以设置

1. 在GOPATH内部-默认为旧的1.10行为（忽略模块）
2. 在GOPATH之外，而在带有go.mod- 的文件树中-默认为模块行为

> auto 未设置或auto-上面的默认行为， 如果找到任何go.mod，即使在GOPATH内部也启用模块模式。（在Go 1.13之前，GO111MODULE=auto永远不会在GOPATH中启用模块模式）
> on 不管目录位置如何，都强制支持模块，当显式启用模块感知模式（通过设置GO111MODULE=on）时，如果不存在go.mod文件，则大多数模块命令的功能将受到更多限制
> off 不管目录位置如何，都强制关闭模块支持

### go mod 常用命令

1. go mod tidy 删除不必要的依赖，添加OS, architecture, and build tags组合所需要的依赖
2. go mod vendor 可选步骤，用于建立vendor文件夹，用于vendor机制的包管理
3. go mod init 将go项目初始化成module-mode，使用go modules进行依赖管理
4. go mod verify 校验go.sum记录的依赖信息是否正确
5. go list -u -m all 要查看所有直接和间接依赖项的可用次要和补丁升级
6. go get example.com/package 要将依赖关系及其所有依赖关系升级到最新版本
7. go get -u example.com/package 要将依赖关系及其所有依赖关系升级到最新版本
8. go test all 在升级或降级任何依赖项之后，您可能想要对构建中的所有软件包（包括直接和间接依赖项）再次运行测试以检查不兼容性

#### 模块由Go源文件树定义，该go.mod文件在树的根目录中。模块源代码可能位于GOPATH之外。有四种指令：module，require，replace，exclude。