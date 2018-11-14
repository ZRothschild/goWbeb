// 文件: web/controllers/movie_controller.go
package controllers

import (
	"errors"
	"../../datamodels"
	"../../services"
	"github.com/kataras/iris"
)
// MovieController是我们的/movies controller。
type MovieController struct {
	//我们的MovieService，它是一个接口，从主应用程序绑定。
	Service services.MovieService
}

//获取电影的返回列表。
// 示例:
// curl -i http://localhost:8080/movies
//
// func (c *MovieController) Get() (results []viewmodels.Movie) {
// 	data := c.Service.GetAll()
// 	for _, movie := range data {
// 		results = append(results, viewmodels.Movie{movie})
// 	}
// 	return
// }
//否则只返回数据模型
func (c *MovieController) Get() (results []datamodels.Movie) {
	return c.Service.GetAll()
}

// MovieByID返回指定id的一部电影
// 例子:
// curl -i http://localhost:8080/movies/1
func (c *MovieController) GetBy(id int64) (movie datamodels.Movie, found bool) {
	return c.Service.GetByID(id) // it will throw 404 if not found.
}

// PutBy 更新指点电影数据.
// 例子:
// curl -i -X PUT -F "genre=Thriller" -F "poster=@/Users/kataras/Downloads/out.gif" http://localhost:8080/movies/1
func (c *MovieController) PutBy(ctx iris.Context, id int64) (datamodels.Movie, error) {
	// get the request data for poster and genre
	file, info, err := ctx.FormFile("poster")
	if err != nil {
		return datamodels.Movie{}, errors.New("failed due form file 'poster' missing")
	}
	//关闭文件。
	file.Close()
	//想象这是上传文件的网址...
	poster := info.Filename
	genre := ctx.FormValue("genre")
	return c.Service.UpdatePosterAndGenreByID(id, poster, genre)
}

// DeleteBy删除电影。
// 例子:
// curl -i -X DELETE -u admin:password http://localhost:8080/movies/1
func (c *MovieController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)
	if wasDel {
		//返回已删除电影的ID
		return iris.Map{"deleted": id}
	}
	//在这里我们可以看到方法函数可以返回这两种类型中的任何一种（map或int），
	//我们不必将返回类型指定为特定类型。
	return iris.StatusBadRequest
}