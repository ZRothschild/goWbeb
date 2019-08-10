####Sessions
Iris提供快速，功能齐全且易于使用的会话管理器。

Iris sessions使用 kataras/iris/sessions 包.

###概述
```go
    import "github.com/kataras/iris/sessions"
    sess := sessions.Start(http.ResponseWriter, *http.Request)
    sess.
      ID() string
      Get(string) interface{}
      HasFlash() bool
      GetFlash(string) interface{}
      GetFlashString(string) string
      GetString(key string) string
      GetInt(key string) (int, error)
      GetInt64(key string) (int64, error)
      GetFloat32(key string) (float32, error)
      GetFloat64(key string) (float64, error)
      GetBoolean(key string) (bool, error)
      GetAll() map[string]interface{}
      GetFlashes() map[string]interface{}
      VisitAll(cb func(k string, v interface{}))
      Set(string, interface{})
      SetImmutable(key string, value interface{})
      SetFlash(string, interface{})
      Delete(string)
      Clear()
      ClearFlashes()
```

此示例将说明如何存储会话中的数据。
除了Iris之外，您不需要任何第三方库，但如果您想要使用任何东西，请记住Iris与标准库完全兼容。按下即可找到更详细的示例

在此示例中，我们将仅允许经过身份验证的用户在/ secret时段查看我们的秘密消息。要获得访问权限，首先必须访问/登录才能获得有效的会话cookie，并将其登录。此外，他可以访问/注销以撤消对我们的秘密消息的访问权限。

```go
    // sessions.go
    package main
    import (
        "github.com/kataras/iris"
        "github.com/kataras/iris/sessions"
    )
    var (
        cookieNameForSessionID = "mycookiesessionnameid"
        sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
    )
    func secret(ctx iris.Context) {
        //验证用户授权
        if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
            ctx.StatusCode(iris.StatusForbidden)
            return
        }
        //输出消息
        ctx.WriteString("The cake is a lie!")
    }
    func login(ctx iris.Context) {
        session := sess.Start(ctx)
        // 在里执行验证
        // ...
        //把验证状态保存为true
        session.Set("authenticated", true)
    }
    func logout(ctx iris.Context) {
        session := sess.Start(ctx)
        // 撤消用户身份验证
        session.Set("authenticated", false)
    }
    func main() {
        app := iris.New()
        app.Get("/secret", secret)
        app.Get("/login", login)
        app.Get("/logout", logout)
        app.Run(iris.Addr(":8080"))
    }
```
> $ go run sessions.go
 
> $ curl -s http://localhost:8080/secret
 Forbidden
 
> $ curl -s -I http://localhost:8080/login
 Set-Cookie: mysessionid=MTQ4NzE5Mz...
 
> $ curl -s --cookie "mysessionid=MTQ4NzE5Mz..." http://localhost:8080/secret
 The cake is a lie!`

####后端存储

有时您需要一个后端存储，即文件存储或redis存储，这将使您的会话数据在服务器重新启动时保持不变。

使用单个调用注册数据库非常容易.UseDatabase（数据库）。

让我们看一个使用快速键值存储螺栓db的简单示例。[bolt](https://github.com/boltdb/bolt)

```go
    package main
    import (
        "time"
        "github.com/kataras/iris"
        "github.com/kataras/iris/sessions"
        "github.com/kataras/iris/sessions/sessiondb/boltdb"
    )
    func main() {
        db, _ := boltdb.New("./sessions/sessions.db", 0666, "users")
        // 使用不同的go协程来同步数据库
        db.Async(true)
        // 按下control + C / cmd + C时关闭并解锁数据库
        iris.RegisterOnInterrupt(func() {
            db.Close()
        })
        sess := sessions.New(sessions.Config{
            Cookie:  "sessionscookieid",
            Expires: 45 * time.Minute, // 0 代表忽略
        })
        // 重要:
        sess.UseDatabase(db)
        // 其余的代码保持不变
        app := iris.New()
        app.Get("/", func(ctx iris.Context) {
            ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
        })
        app.Get("/set", func(ctx iris.Context) {
            s := sess.Start(ctx)
            //设置
            s.Set("name", "iris")
            //测试如果在这里设置
            ctx.Writef("All ok session setted to: %s", s.GetString("name"))
        })
        app.Get("/set/{key}/{value}", func(ctx iris.Context) {
            key, value := ctx.Params().Get("key"), ctx.Params().Get("value")
            s := sess.Start(ctx)
            // 设置会话值
            s.Set(key, value)
            // 测试如果在这里设置
            ctx.Writef("All ok session setted to: %s", s.GetString(key))
        })
        app.Get("/get", func(ctx iris.Context) {
            // 获取一个特定的键，如字符串，如果没有找到只返回一个空字符串
            name := sess.Start(ctx).GetString("name")
            ctx.Writef("The name on the /set was: %s", name)
        })
        app.Get("/get/{key}", func(ctx iris.Context) {
            // 获取一个特定的键，如字符串，如果没有找到只返回一个空字符串
            name := sess.Start(ctx).GetString(ctx.Params().Get("key"))
            ctx.Writef("The name on the /set was: %s", name)
        })
        app.Get("/delete", func(ctx iris.Context) {
            // 删除一个具体的值
            sess.Start(ctx).Delete("name")
        })
        app.Get("/clear", func(ctx iris.Context) {
            // 删除所有条目
            sess.Start(ctx).Clear()
        })
        app.Get("/destroy", func(ctx iris.Context) {
            //destroy，删除整个会话数据和cookie
            sess.Destroy(ctx)
        })
        app.Get("/update", func(ctx iris.Context) {
            //更新过期日期与新日期
            sess.ShiftExpiration(ctx)
        })
        app.Run(iris.Addr(":8080"))
    }
```
创建自定义后端会话存储

您可以通过实现Database接口来创建自己的后端存储。
```go
    type Database interface {
        Load(sid string) returns struct {
            //值包含整个内存存储，此存储
            //包含从内存调用更新的当前值，
            //会话数据（键和值）。这条路
            //数据库可以访问整个会话的数据
            // 每次。
            Values memstore.Store
            //在插入时它包含到期日期时间
            //在更新时它包含新的到期日期时间（如果更新或旧的）
            //在删除时它将为零
            //明确时它将为零
            //在销毁时它将为零
            Lifetime LifeTime
        }
        Sync(accepts struct {
            //值包含整个内存存储，此存储
            //包含从内存调用更新的当前值，
            //会话数据（键和值）。这条路
            //数据库可以访问整个会话的数据每次。
            Values memstore.Store

            //在插入时它包含到期日期时间
            //在更新时它包含新的到期日期时间（如果更新或旧的）
            //在删除时它将为零
            //明确时它将为零
            //在销毁时它将为零
            Lifetime LifeTime
        })
    }
```
这就是boltdb会话数据库的样子
```go
    package boltdb
    import (
        "bytes"
        "os"
        "path/filepath"
        "runtime"
        "time"
        "github.com/boltdb/bolt"
        "github.com/kataras/golog"
        "github.com/kataras/iris/core/errors"
        "github.com/kataras/iris/sessions"
    )
    // DefaultFileMode用作默认数据库的“fileMode”
    //用于创建会话目录路径，打开和写入
    //会话boltdb（基于文件）存储。
    var (
        DefaultFileMode = 0666
    )
    //数据库BoltDB（基于文件）会话存储。
    type Database struct {
        table []byte
        //服务是下划线BoltDB数据库连接，
        //它在`New`或`NewFromDB`初始化。
        //可用于获取统计数据。
        Service *bolt.DB
        async   bool
    }
    var (
        //当path或tableName为空时，ErrOptionsMissing在`New`上返回。
        ErrOptionsMissing = errors.New("required options are missing")
    )
    // New创建并返回一个新的BoltDB（基于文件）存储
    //基于“路径”的实例。
    //路径应包括文件名和目录（也称为fullpath），即sessions / store.db。
    //它将删除任何旧的会话文件。
    func New(path string, fileMode os.FileMode, bucketName string) (*Database, error) {
        if path == "" || bucketName == "" {
            return nil, ErrOptionsMissing
        }
        if fileMode <= 0 {
            fileMode = os.FileMode(DefaultFileMode)
        }
        // create directories if necessary
        if err := os.MkdirAll(filepath.Dir(path), fileMode); err != nil {
            golog.Errorf("error while trying to create the necessary directories for %s: %v", path, err)
            return nil, err
        }
        service, err := bolt.Open(path, 0600,
            &bolt.Options{Timeout: 15 * time.Second},
        )
        if err != nil {
            golog.Errorf("unable to initialize the BoltDB-based session database: %v", err)
            return nil, err
        }
        return NewFromDB(service, bucketName)
    }
    // NewFromDB与`New`相同，但接受已创建的自定义boltdb连接。
    func NewFromDB(service *bolt.DB, bucketName string) (*Database, error) {
        if bucketName == "" {
            return nil, ErrOptionsMissing
        }
        bucket := []byte(bucketName)
        service.Update(func(tx *bolt.Tx) (err error) {
            _, err = tx.CreateBucketIfNotExists(bucket)
            return
        })
        db := &Database{table: bucket, Service: service}
        runtime.SetFinalizer(db, closeDB)
        return db, db.Cleanup()
    }
    //清理会删除所有无效（已过期）的会话条目，
    //它也会在“新”上自动调用。
    func (db *Database) Cleanup() error {
        err := db.Service.Update(func(tx *bolt.Tx) error {
            b := db.getBucket(tx)
            c := b.Cursor()
            for k, v := c.First(); k != nil; k, v = c.Next() {
                if len(k) == 0 { // empty key, continue to the next pair
                    continue
                }
                storeDB, err := sessions.DecodeRemoteStore(v)
                if err != nil {
                    continue
                }
                if storeDB.Lifetime.HasExpired() {
                    if err := c.Delete(); err != nil {
                        golog.Warnf("troubles when cleanup a session remote store from BoltDB: %v", err)
                    }
                }
            }
            return nil
        })
        return err
    }
    // Async如果为true，那么它将使用不同的
    //go协程来更新BoltDB（基于文件的）存储。
    func (db *Database) Async(useGoRoutines bool) *Database {
        db.async = useGoRoutines
        return db
    }
    //加载来自BoltDB（基于文件）会话存储的会话。
    func (db *Database) Load(sid string) (storeDB sessions.RemoteStore) {
        bsid := []byte(sid)
        err := db.Service.View(func(tx *bolt.Tx) (err error) {
            // db.getSessBucket(tx, sid)
            b := db.getBucket(tx)
            c := b.Cursor()
            for k, v := c.First(); k != nil; k, v = c.Next() {
                if len(k) == 0 { // empty key, continue to the next pair
                    continue
                }
                if bytes.Equal(k, bsid) { // session id should be the name of the key-value pair
                    storeDB, err = sessions.DecodeRemoteStore(v) // decode the whole value, as a remote store
                    break
                }
            }
            return
        })
        if err != nil {
            golog.Errorf("error while trying to load from the remote store: %v", err)
        }
        return
    }
    // Sync同步数据库与会话（内存）存储。
    func (db *Database) Sync(p sessions.SyncPayload) {
        if db.async {
            go db.sync(p)
        } else {
            db.sync(p)
        }
    }
    func (db *Database) sync(p sessions.SyncPayload) {
        bsid := []byte(p.SessionID)
        if p.Action == sessions.ActionDestroy {
            if err := db.destroy(bsid); err != nil {
                golog.Errorf("error while destroying a session(%s) from boltdb: %v",
                    p.SessionID, err)
            }
            return
        }
        s, err := p.Store.Serialize()
        if err != nil {
            golog.Errorf("error while serializing the remote store: %v", err)
        }
        err = db.Service.Update(func(tx *bolt.Tx) error {
            return db.getBucket(tx).Put(bsid, s)
        })
        if err != nil {
            golog.Errorf("error while writing the session bucket: %v", err)
        }
    }
    func (db *Database) destroy(bsid []byte) error {
        return db.Service.Update(func(tx *bolt.Tx) error {
            return db.getBucket(tx).Delete(bsid)
        })
    }
    func (db *Database) getBucket(tx *bolt.Tx) *bolt.Bucket {
        return tx.Bucket(db.table)
    }
    // Len报告存储到此BoltDB表的会话数。
    func (db *Database) Len() (num int) {
        db.Service.View(func(tx *bolt.Tx) error {
            // Assume bucket exists and has keys
            b := db.getBucket(tx)
            if b == nil {
                return nil
            }
            b.ForEach(func([]byte, []byte) error {
                num++
                return nil
            })
            return nil
        })
        return
    }
    // 关闭BoltDB连接。FUNC
    func (db *Database) Close() error {
        return closeDB(db)
    }
    func closeDB(db *Database) error {
        err := db.Service.Close()
        if err != nil {
            golog.Warnf("closing the BoltDB connection: %v", err)
        }
        return err
    }
```