// file: repositories/movie_repository.go

package repositories

import (
	"errors"
	"sync"

	"github.com/kataras/iris/_examples/hero/overview/datamodels"
)

// Query表示访问者和操作查询
type Query func(datamodels.Movie) bool

// MovieRepository处理电影实体/模型的基本操作。
//它是一个可测试的接口，即一个内存电影库或连接到sql数据库。
type MovieRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)

	Select(query Query) (movie datamodels.Movie, found bool)
	SelectMany(query Query, limit int) (results []datamodels.Movie)

	InsertOrUpdate(movie datamodels.Movie) (updatedMovie datamodels.Movie, err error)
	Delete(query Query, limit int) (deleted bool)
}

//NewMovieRepository返回一个新的基于电影内存的repository，
//我们示例中唯一的repository类型。
func NewMovieRepository(source map[int64]datamodels.Movie) MovieRepository {
	return &movieMemoryRepository{source: source}
}

// movieMemoryRepository is a "MovieRepository"
// which manages the movies using the memory data source (map).

// movieMemoryRepository是一个MovieRepository
//使用内存数据源（map）管理电影。
type movieMemoryRepository struct {
	source map[int64]datamodels.Movie
	mu     sync.RWMutex
}

const (
	// ReadOnlyMode将RLock（读取）数据。
	ReadOnlyMode = iota
	// ReadWriteMode将锁定（读/写）数据。
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

// Select receives a query function
// which is fired for every single movie model inside
// our imaginary data source.
// When that function returns true then it stops the iteration.
//
// It returns the query's return last known "found" value
// and the last known movie model
// to help callers to reduce the LOC.
//
// It's actually a simple but very clever prototype function
// I'm using everywhere since I firstly think of it,
// hope you'll find it very useful as well.

//Select接收查询函数
//为内部的每个电影模型触发
//我们想象中的数据源
//当该函数返回true时，它会停止迭代。
//
//它返回查询返回的最后一个已知“找到”值
//和最后一个已知的电影模型
//帮助呼叫者减少LOC。
//
//它实际上是一个简单但非常聪明的原型函数
//自从我第一次想到它以来，我一直在使用它，
//希望你会发现它也很有用。
func (r *movieMemoryRepository) Select(query Query) (movie datamodels.Movie, found bool) {
	found = r.Exec(query, func(m datamodels.Movie) bool {
		movie = m
		return true
	}, 1, ReadOnlyMode)

	//设置一个空的datamodels.Movie，如果根本找不到的话。
	if !found {
		movie = datamodels.Movie{}
	}

	return
}

// SelectMany与Select相同但返回一个或多个datamodels.Movie作为切片。
//如果limit <= 0则返回所有内容。
func (r *movieMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.Movie) {
	r.Exec(query, func(m datamodels.Movie) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)

	return
}

// InsertOrUpdate将movie添加或更新到map中存储。
//返回新电影，如果有则返回错误。
func (r *movieMemoryRepository) InsertOrUpdate(movie datamodels.Movie) (datamodels.Movie, error) {
	id := movie.ID
	if id == 0 { // Create new action  //创建新记录
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
	// r.source [id] =movie
	//向下面的代码一样;
	current, exists := r.Select(func(m datamodels.Movie) bool {
		return m.ID == id
	})
	if !exists { // ID不是真实的，返回错误。
		return datamodels.Movie{}, errors.New("failed to update a nonexistent movie")
	}
	//或注释这些和r.source[id] = m进行纯替换
	if movie.Poster != "" {
		current.Poster = movie.Poster
	}
	if movie.Genre != "" {
		current.Genre = movie.Genre
	}
	//锁定数据
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
