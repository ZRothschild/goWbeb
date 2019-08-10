package main

import (
	"github.com/kataras/iris"
	"fmt"
	"github.com/kataras/iris/context"
)

func main(){
	app := iris.New()
	app.UseGlobal(BeginRequest)
	app.PartyFunc("/users", func(users iris.Party) {
		users.Use(myAuthMiddlewareHandler)

		// http://localhost:8080/users/42/profile
		users.Get("/{id:int}/profile", userProfileHandler)
		// http://localhost:8080/users/messages/1
		users.Get("/messages/{id:int}", userMessageHandler)
	})

	app.Run(iris.Addr(":8080"))
}
func myAuthMiddlewareHandler(ctx iris.Context){
	ctx.WriteString("Authentication failed\n\r")

	ctx.Next()
}
func userProfileHandler(ctx iris.Context) {//

	id:=ctx.Params().Get("id")
	ctx.WriteString(id)
}

func userMessageHandler(ctx iris.Context){
	id:=ctx.Params().Get("id")
	ctx.Handlers()[3](ctx)
	ctx.WriteString(id+"message")
}
func add(ctx context.Context)  {
	ctx.Handlers()[3](ctx)//定义再handler里面的位置
	fmt.Println("addHandler is runing")

}
func BeginRequest(ctx context.Context){

	fmt.Println("yes it is")
	ctx.AddHandler(add)

	ctx.Next()

}