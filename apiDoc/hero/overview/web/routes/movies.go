// file: web/routes/movie.go
package routes

import (
	"errors"
	"../../datamodels"
	"../../services"
	"github.com/kataras/iris"
)

//Movies返回电影列表
// 例子:
// curl -i http://localhost:8080/movies
func Movies(service services.MovieService) (results []datamodels.Movie) {
	return service.GetAll()
}

// MovieByID返回指定id的一部电影
// 例子:
// curl -i http://localhost:8080/movies/1
func MovieByID(service services.MovieService, id int64) (movie datamodels.Movie, found bool) {
	return service.GetByID(id) //如果未找到返回404.
}

// UpdateMovieByID更新电影。
// 例子:
// curl -i -X PUT -F "genre=Thriller" -F "poster=@/Users/kataras/Downloads/out.gif" http://localhost:8080/movies/1
func UpdateMovieByID(ctx iris.Context, service services.MovieService, id int64) (datamodels.Movie, error) {
	//获取海报和流派的请求数据
	file, info, err := ctx.FormFile("poster")
	if err != nil {
		return datamodels.Movie{}, errors.New("failed due form file 'poster' missing")
	}
	//关闭文件。
	file.Close()
	//想象这是上传文件的网址...
	poster := info.Filename
	genre := ctx.FormValue("genre")
	return service.UpdatePosterAndGenreByID(id, poster, genre)
}

// DeleteMovieByID删除电影。
// 例子:
// curl -i -X DELETE -u admin:password http://localhost:8080/movies/1
func DeleteMovieByID(service services.MovieService, id int64) interface{} {
	wasDel := service.DeleteByID(id)
	if wasDel {
		//返回已删除电影的ID
		return iris.Map{"deleted": id}
	}
	//在这里我们可以看到方法函数可以返回这两种类型中的任何一种（map或int），
	//我们不必将返回类型指定为特定类型。
	return iris.StatusBadRequest
}