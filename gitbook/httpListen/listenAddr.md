# `IRIS`服务器监听地址

## `IRIS`服务器监听地址普通

### 目录结构
> 主目录`listenAddr`

```html
    —— main.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1>Hello World!</h1>")
	})
	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}
```
## `IRIS`服务器监听地址关闭错误
### 目录结构
> 主目录`omitServerErrors`

```html
    —— main.go
    —— main_test.go
```
### 代码示例
> `main.go`

```go
package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1>Hello World!</h1>")
	})

	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		// do something
	}
	// same as:
	// err := app.Run(iris.Addr(":8080"))
	// if err != nil && (err != iris.ErrServerClosed || err.Error() != iris.ErrServerClosed.Error()) {
	//     [...]
	// }
}
```
> `main_test.go`

```go
package main

import (
	"bytes"
	stdContext "context"
	"strings"
	"testing"
	"time"

	"github.com/kataras/iris"
)

func logger(app *iris.Application) *bytes.Buffer {
	buf := &bytes.Buffer{}

	app.Logger().SetOutput(buf)

	// disable the "Now running at...." in order to have a clean log of the error.
	// we could attach that on `Run` but better to keep things simple here.
	app.Configure(iris.WithoutStartupLog)
	return buf
}

func TestListenAddr(t *testing.T) {
	app := iris.New()
	// we keep the logger running as well but in a controlled way.
	log := logger(app)

	// close the server at 3-6 seconds
	go func() {
		time.Sleep(3 * time.Second)
		ctx, cancel := stdContext.WithTimeout(stdContext.TODO(), 3*time.Second)
		defer cancel()
		app.Shutdown(ctx)
	}()

	err := app.Run(iris.Addr(":9829"))
	// in this case the error should be logged and return as well.
	if err != iris.ErrServerClosed {
		t.Fatalf("expecting err to be `iris.ErrServerClosed` but got: %v", err)
	}

	expectedMessage := iris.ErrServerClosed.Error()

	if got := log.String(); !strings.Contains(got, expectedMessage) {
		t.Fatalf("expecting to log to contains the:\n'%s'\ninstead of:\n'%s'", expectedMessage, got)
	}

}

func TestListenAddrWithoutServerErr(t *testing.T) {
	app := iris.New()
	// we keep the logger running as well but in a controlled way.
	log := logger(app)

	// close the server at 3-6 seconds
	go func() {
		time.Sleep(3 * time.Second)
		ctx, cancel := stdContext.WithTimeout(stdContext.TODO(), 3*time.Second)
		defer cancel()
		app.Shutdown(ctx)
	}()

	// we disable the ErrServerClosed, so the error should be nil when server is closed by `app.Shutdown`.

	// so in this case the iris/http.ErrServerClosed should be NOT logged and NOT return.
	err := app.Run(iris.Addr(":9827"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		t.Fatalf("expecting err to be nil but got: %v", err)
	}

	if got := log.String(); got != "" {
		t.Fatalf("expecting to log nothing but logged: '%s'", got)
	}
}
```