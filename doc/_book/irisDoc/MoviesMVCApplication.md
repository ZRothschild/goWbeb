####MVC movie 应用程序

Iris有一个非常强大和[极快](https://github.com/kataras/iris/tree/master/_benchmarks)的MVC支持，你可以从方法函数返回任何类型的任何值，它将按预期发送到客户端。

* 如果是string ，那就是body。
* 如果string是第二个输出参数，那么它就是内容类型。
* 如果int是状态码。
* 如果错误而不是nil，那么(任何类型)响应将被省略，错误的文本将呈现400个错误请求。
* 如果(int, error)和error不是nil，那么响应结果将是错误的文本，状态码为int。
* 如果定制struct或interface{}或slice或map，那么它将被呈现为json，除非后面是字符串内容类型。
* 如果mvc。结果，然后它执行它的调度功能，所以好的设计模式可以用来分割模型的逻辑在需要的地方。

没有什么能阻止您使用自己喜欢的文件夹结构。 Iris是一个低级Web框架，它有MVC一流的支持，但它不限制你的文件夹结构，这是你的选择。

结构取决于您自己的需求。我们无法告诉您如何设计自己的应用程序，但您可以自由地仔细查看下面的一个用例示例;

![folder](images/folder_structure.png)

#####model层

让我们从我们的movie model 开始。

```go
    // file: datamodels/movie.go
    package datamodels
    // Movie 是我们的一个简单的数据结构体
    // 请注意公共标签（适用于我们web网络应用）
    // 应保存在“web / viewmodels / movie.go”等其他文件中
    //可以通过嵌入datamodels.Movie或
    //声明新字段但我们将使用此数据模型
    //作为我们应用程序中唯一的一个Movie模型，
    //为了摇摇欲坠。
    type Movie struct {
        ID     int64  `json:"id"`
        Name   string `json:"name"`
        Year   int    `json:"year"`
        Genre  string `json:"genre"`
        Poster string `json:"poster"`
    }
```
#####数据源/数据存储层

之后，我们继续为我们的Movie创建一个简单的内存存储。

```go
    // file: datasource/movies.go
    package datasource
    import "github.com/kataras/iris/_examples/mvc/overview/datamodels"
    // Movies is our imaginary data source.
    var Movies = map[int64]datamodels.Movie{
        1: {
            ID:     1,
            Name:   "Casablanca",
            Year:   1942,
            Genre:  "Romance",
            Poster: "https://iris-go.com/images/examples/mvc-movies/1.jpg",
        },
        2: {
            ID:     2,
            Name:   "Gone with the Wind",
            Year:   1939,
            Genre:  "Romance",
            Poster: "https://iris-go.com/images/examples/mvc-movies/2.jpg",
        },
        3: {
            ID:     3,
            Name:   "Citizen Kane",
            Year:   1941,
            Genre:  "Mystery",
            Poster: "https://iris-go.com/images/examples/mvc-movies/3.jpg",
        },
        4: {
            ID:     4,
            Name:   "The Wizard of Oz",
            Year:   1939,
            Genre:  "Fantasy",
            Poster: "https://iris-go.com/images/examples/mvc-movies/4.jpg",
        },
        5: {
            ID:     5,
            Name:   "North by Northwest",
            Year:   1959,
            Genre:  "Thriller",
            Poster: "https://iris-go.com/images/examples/mvc-movies/5.jpg",
        },
    }
```
#####Repositories

可以直接访问“数据源”并可以直接操作数据的层。

可选（因为您也可以在服务中使用它），但是该示例需要，我们创建一个存储库，一个存储库处理“低级”，直接访问Movies数据源。保留“存储库”，它是一个接口，因为它可能不同，它取决于您的应用程序开发的状态，即在生产中将使用一些真正的SQL查询或您用于查询数据的任何其他内容。

```go
    package repositories
    import (
        "errors"
        "sync"
        "github.com/kataras/iris/_examples/mvc/overview/datamodels"
    )
    // Query represents the visitor and action queries.
    type Query func(datamodels.Movie) bool
    // MovieRepository handles the basic operations of a movie entity/model.
    // It's an interface in order to be testable, i.e a memory movie repository or
    // a connected to an sql database.
    type MovieRepository interface {
        Exec(query Query, action Query, limit int, mode int) (ok bool)
        Select(query Query) (movie datamodels.Movie, found bool)
        SelectMany(query Query, limit int) (results []datamodels.Movie)
        InsertOrUpdate(movie datamodels.Movie) (updatedMovie datamodels.Movie, err error)
        Delete(query Query, limit int) (deleted bool)
    }
    // NewMovieRepository returns a new movie memory-based repository,
    // the one and only repository type in our example.
    func NewMovieRepository(source map[int64]datamodels.Movie) MovieRepository {
        return &movieMemoryRepository{source: source}
    }
    // movieMemoryRepository is a "MovieRepository"
    // which manages the movies using the memory data source (map).
    type movieMemoryRepository struct {
        source map[int64]datamodels.Movie
        mu     sync.RWMutex
    }
    const (
        // ReadOnlyMode will RLock(read) the data .
        ReadOnlyMode = iota
        // ReadWriteMode will Lock(read/write) the data.
        ReadWriteMode
    )
    func (r *movieMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
        loops := 0
        if mode == ReadOnlyMode {
            r.mu.RLock()
            defer r.mu.RUnlock()
        } else {
            r.mu.Lock()
            defer r.mu.Unlock()
        }
        for _, movie := range r.source {
            ok = query(movie)
            if ok {
                if action(movie) {
                    loops++
                    if actionLimit >= loops {
                        break // break
                    }
                }
            }
        }
        return
    }
    //选择接收查询功能
    //为内部的每个电影模型触发
    //我们想象中的数据源
    //当该函数返回true时，它会停止迭代。

    //它返回查询返回的最后一个已知“找到”值
    //和最后一个已知的电影模型
    //帮助呼叫者减少LOC。

    //它实际上是一个简单但非常聪明的原型函数
    //自从我第一次想到它以来，我一直在使用它，
    //希望你会发现它也很有用。
    func (r *movieMemoryRepository) Select(query Query) (movie datamodels.Movie, found bool) {
        found = r.Exec(query, func(m datamodels.Movie) bool {
            movie = m
            return true
        }, 1, ReadOnlyMode)

        //如果根本找不到的话,设置一个空的datamodels.Movie，
        if !found {
            movie = datamodels.Movie{}
        }
        return
    }
    // SelectMany与Select相同但返回一个或多个datamodels.Movie作为切片。
    //如果limit <= 0则返回所有内容
    func (r *movieMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.Movie) {
        r.Exec(query, func(m datamodels.Movie) bool {
            results = append(results, m)
            return true
        }, limit, ReadOnlyMode)
        return
    }
    // InsertOrUpdate将影片添加或更新到（内存）存储。
    // 返回新电影，如果有则返回错误。
    func (r *movieMemoryRepository) InsertOrUpdate(movie datamodels.Movie) (datamodels.Movie, error) {
        id := movie.ID
        if id == 0 { // Create new action
            var lastID int64
            //找到最大的ID，以便不重复
            //在制作应用中，您可以使用第三方
            //库以生成UUID作为字符串。
            r.mu.RLock()
            for _, item := range r.source {
                if item.ID > lastID {
                    lastID = item.ID
                }
            }
            r.mu.RUnlock()
            id = lastID + 1
            movie.ID = id
            // map-specific thing
            r.mu.Lock()
            r.source[id] = movie
            r.mu.Unlock()

            return movie, nil
        }
        //基于movie.ID更新动作，
        //这里我们将允许更新海报和流派，如果不是空的话。
        //或者我们可以做替换：
        // r.source [id] =电影
        //并评论下面的代码;
        current, exists := r.Select(func(m datamodels.Movie) bool {
            return m.ID == id
        })
        if !exists { //ID不是真实的，返回错误。
            return datamodels.Movie{}, errors.New("failed to update a nonexistent movie")
        }
        // 或者注释这些和r.source [id] = m进行纯替换
        if movie.Poster != "" {
            current.Poster = movie.Poster
        }
        if movie.Genre != "" {
            current.Genre = movie.Genre
        }
        // map-specific thing
        r.mu.Lock()
        r.source[id] = current
        r.mu.Unlock()
        return movie, nil
    }
    func (r *movieMemoryRepository) Delete(query Query, limit int) bool {
        return r.Exec(query, func(m datamodels.Movie) bool {
            delete(r.source, m.ID)
            return true
        }, limit, ReadWriteMode)
    }
```
#####服务层 service

service可以访问“存储库”和“模型”（甚至是“数据模型”，如果是简单应用程序）的函数的层。它应该包含大部分域逻辑。

我们需要一个服务来与我们的存储库在“高级”和存储/检索电影中进行通信，这将在下面的Web控制器上使用。

```go
    // file: services/movie_service.go
    package services
    import (
        "github.com/kataras/iris/_examples/mvc/overview/datamodels"
        "github.com/kataras/iris/_examples/mvc/overview/repositories"
    )
    // MovieService处理电影数据模型的一些CRUID操作。
    //这取决于影片库的动作。
    //这是将数据源与更高级别的组件分离。
    //因此，不同的存储库类型可以使用相同的逻辑，而无需任何更改。
    //它是一个界面，它在任何地方都被用作界面
    //因为我们可能需要在将来更改或尝试实验性的不同域逻辑。
    type MovieService interface {
        GetAll() []datamodels.Movie
        GetByID(id int64) (datamodels.Movie, bool)
        DeleteByID(id int64) bool
        UpdatePosterAndGenreByID(id int64, poster string, genre string) (datamodels.Movie, error)
    }
    // NewMovieService返回默认 movie service.
    func NewMovieService(repo repositories.MovieRepository) MovieService {
        return &movieService{
            repo: repo,
        }
    }
    type movieService struct {
        repo repositories.MovieRepository
    }
    // GetAll 获取所有的movie.
    func (s *movieService) GetAll() []datamodels.Movie {
        return s.repo.SelectMany(func(_ datamodels.Movie) bool {
            return true
        }, -1)
    }
    // GetByID 根据其ID返回一行。
    func (s *movieService) GetByID(id int64) (datamodels.Movie, bool) {
        return s.repo.Select(func(m datamodels.Movie) bool {
            return m.ID == id
        })
    }
    // UpdatePosterAndGenreByID更新电影的海报和流派。
    func (s *movieService) UpdatePosterAndGenreByID(id int64, poster string, genre string) (datamodels.Movie, error) {
        // update the movie and return it.
        return s.repo.InsertOrUpdate(datamodels.Movie{
            ID:     id,
            Poster: poster,
            Genre:  genre,
        })
    }
    // DeleteByID按ID删除电影。
    //
    //如果删除则返回true，否则返回false。
    func (s *movieService) DeleteByID(id int64) bool {
        return s.repo.Delete(func(m datamodels.Movie) bool {
            return m.ID == id
        }, 1)
    }
```
#####视图View Models

应该有视图模型，客户端将能够看到的结构。例：

```go
    import (
        "github.com/kataras/iris/_examples/mvc/overview/datamodels"
        "github.com/kataras/iris/context"
    )
    type Movie struct {
        datamodels.Movie
    }
    func (m Movie) IsValid() bool {
        /* 做一些检查，如果有效则返回true。.. */
        return m.ID > 0
    }
```

Iris能够将任何自定义数据结构转换为HTTP响应调度程序，因此从理论上讲，如果确实需要，则允许使用以下内容：

```go
    // Dispatch完成`kataras / iris / mvc＃Result`界面。
    //将“电影”作为受控的http响应发送。
    //如果其ID为零或更小，则返回404未找到错误
    //否则返回其json表示，
    //（就像控制器的函数默认为自定义类型一样）。
    //
    //不要过度，应用程序的逻辑不应该在这里。
    //在响应之前，这只是验证的又一步，
    //可以在这里添加简单的检查。
    //
    //这只是一个展示
    //想象一下设计更大的应用程序时此功能给出的潜力。
    //
    //调用控制器方法返回值的函数
    //是“电影”的类型。
    //例如`controllers / movie_controller.go＃GetBy`.
    func (m Movie) Dispatch(ctx context.Context) {
        if !m.IsValid() {
            ctx.NotFound()
            return
        }
        ctx.JSON(m, context.JSON{Indent: " "})
    }
```

但是，我们将使用“datamodels”作为唯一的模型包，因为Movie结构不包含任何敏感数据，客户端能够查看其所有字段，并且我们不需要任何额外的功能或验证。

##### 控制器 Controllers

处理Web请求，在服务和客户端之间架起桥梁。

而最重要的是，Iris来自哪里，是与MovieService进行通信的Controller。我们通常将所有与http相关的东西存储在一个名为“web”的不同文件夹中，这样所有控制器都可以在“web / controllers”中，
注意“通常”你也可以使用其他设计模式，这取决于你。

```go
    //file: web/controllers/movie_controller.go
    package controllers
    import (
        "errors"
        "github.com/kataras/iris/_examples/mvc/overview/datamodels"
        "github.com/kataras/iris/_examples/mvc/overview/services"
        "github.com/kataras/iris"
    )
    // MovieController
    type MovieController struct {
       //我们的MovieService，它是一个界面
        //从主应用程序绑定。
        Service services.MovieService
    }
    // 获取电影列表
    // curl -i http://localhost:8080/movies
    // 如果您有敏感数据，这是正确的方法：
    // func (c *MovieController) Get() (results []viewmodels.Movie) {
    //     data := c.Service.GetAll()
    //     for _, movie := range data {
    //         results = append(results, viewmodels.Movie{movie})
    //     }
    //     return
    // }
    // 否则只返回数据模型
    func (c *MovieController) Get() (results []datamodels.Movie) {
        return c.Service.GetAll()
    }
    //获取一部电影
    // Demo:
    // curl -i http://localhost:8080/movies/1
    func (c *MovieController) GetBy(id int64) (movie datamodels.Movie, found bool) {
        return c.Service.GetByID(id) // it will throw 404 if not found.
    }
    // 用put请求更新一部电影
    // Demo:
    // curl -i -X PUT -F "genre=Thriller" -F "poster=@/Users/kataras/Downloads/out.gif" http://localhost:8080/movies/1
    func (c *MovieController) PutBy(ctx iris.Context, id int64) (datamodels.Movie, error) {
        // get the request data for poster and genre
        file, info, err := ctx.FormFile("poster")
        if err != nil {
            return datamodels.Movie{}, errors.New("failed due form file 'poster' missing")
        }
        // 不需要文件所以关闭他
        file.Close()
        //想象一下，这是上传文件的网址......
        poster := info.Filename
        genre := ctx.FormValue("genre")
        return c.Service.UpdatePosterAndGenreByID(id, poster, genre)
    }
    // Delete请求删除一部电影
    // curl -i -X DELETE -u admin:password http://localhost:8080/movies/1
    func (c *MovieController) DeleteBy(id int64) interface{} {
        wasDel := c.Service.DeleteByID(id)
        if wasDel {
            // 返回删除的id
            return iris.Map{"deleted": id}
        }
    //在这里我们可以看到方法函数可以返回这两种类型中的任何一种（map或int），
        //我们不必将返回类型指定为特定类型。
        return iris.StatusBadRequest
    }
```

“web / middleware”中的一个中间件，用于动画示例。

```go
    // file: web/middleware/basicauth.go
    package middleware
    import "github.com/kataras/iris/middleware/basicauth"
    // 简单的授权验证
    var BasicAuth = basicauth.New(basicauth.Config{
        Users: map[string]string{
            "admin": "password",
        },
    })
```

最后我们的main.go.

```go
    // file: main.go
    package main
    import (
        "github.com/kataras/iris/_examples/mvc/overview/datasource"
        "github.com/kataras/iris/_examples/mvc/overview/repositories"
        "github.com/kataras/iris/_examples/mvc/overview/services"
        "github.com/kataras/iris/_examples/mvc/overview/web/controllers"
        "github.com/kataras/iris/_examples/mvc/overview/web/middleware"
        "github.com/kataras/iris"
        "github.com/kataras/iris/mvc"
    )
    func main() {
        app := iris.New()
        app.Logger().SetLevel("debug")
        //加载模板文件
        app.RegisterView(iris.HTML("./web/views", ".html"))
        // 注册控制器
        // mvc.New(app.Party("/movies")).Handle(new(controllers.MovieController))
       //您还可以拆分您编写的代码以配置mvc.Application
        //使用`mvc.Configure`方法，如下所示。
        mvc.Configure(app.Party("/movies"), movies)
        // http://localhost:8080/movies
        // http://localhost:8080/movies/1
        app.Run(
            //开启web服务
            iris.Addr("localhost:8080"),
            // 禁用更新
            iris.WithoutVersionChecker,
            // 按下CTRL / CMD + C时跳过错误的服务器：
            iris.WithoutServerError(iris.ErrServerClosed),
            //实现更快的json序列化和更多优化：
            iris.WithOptimizations,
        )
    }
    //注意mvc.Application，它不是iris.Application。
    func movies(app *mvc.Application) {
    //添加基本身份验证（admin：password）中间件
        //用于基于/电影的请求。
        app.Router.Use(middleware.BasicAuth)
        // 使用数据源中的一些（内存）数据创建我们的电影资源库。
        repo := repositories.NewMovieRepository(datasource.Movies)
        // 创建我们的电影服务，我们将它绑定到电影应用程序的依赖项。
        movieService := services.NewMovieService(repo)
        app.Register(movieService)
         //为我们的电影控制器服务
        //请注意，您可以为多个控制器提供服务
        //你也可以使用`movies.Party（relativePath）`或`movies.Clone（app.Party（...））创建子mvc应用程序
        app.Handle(new(controllers.MovieController))
    }
```