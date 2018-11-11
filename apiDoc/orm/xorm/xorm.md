# `xorm`包使用
## 示例代码
```go
//包主显示如何在您的Web应用程序中使用orm
//它只是插入一列并选择第一列。
package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)
/*
	go get -u github.com/mattn/go-sqlite3
	go get -u github.com/go-xorm/xorm
	如果您使用的是win64并且无法安装go-sqlite3：
		1.下载：https：//sourceforge.net/projects/mingw-w64/files/latest/download
		2.选择“x86_x64”和“posix”
		3.添加C:\Program Files\mingw-w64\x86_64-7.1.0-posix-seh-rt_v5-rev1\mingw64\bin
		到你的PATH env变量。
	手册: http://xorm.io/docs/
*/
//User是我们的用户表结构。
type User struct {
	ID        int64  // xorm默认自动递增
	Version   string `xorm:"varchar(200)"`
	Salt      string
	Username  string
	Password  string    `xorm:"varchar(200)"`
	Languages string    `xorm:"varchar(200)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func main() {
	app := iris.New()
	orm, err := xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}
	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})
	err = orm.Sync2(new(User))
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}
	app.Get("/insert", func(ctx iris.Context) {
		user := &User{Username: "kataras", Salt: "hash---", Password: "hashed", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		orm.Insert(user)
		ctx.Writef("user inserted: %#v", user)
	})
	app.Get("/get", func(ctx iris.Context) {
		user := User{ID: 1}
		if ok, _ := orm.Get(&user); ok {
			ctx.Writef("user found: %#v", user)
		}
	})
	// http://localhost:8080/insert
	// http://localhost:8080/get
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```