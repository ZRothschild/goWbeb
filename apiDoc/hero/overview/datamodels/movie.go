//文件: datamodels/movie.go
package datamodels

//Movie是我们的示例数据结构。
//注意结构体字段可导出
//应该保存在其他文件中，例如web/viewmodels/movie.go
//可以通过嵌入datamodels.Movie或 声明新字段但我们将使用此数据模型
//作为我们应用程序中唯一的一个Movie模型，为了不冲突。
type Movie struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
	Poster string `json:"poster"`
}