# Go操作MySQL

```go
import (
   "database/sql"
   "fmt"
   _ "github.com/go-sql-driver/mysql" // init()
)
```

## 初始化数据库

```go
dsn := "root:961024@tcp(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
//打开数据库
sql.Open("mysql", dsn)
// 尝试与数据库建立连接
err = db.Ping()
	db.SetMaxOpenConns(100) //最大连接数
	db.SetMaxIdleConns(10)  //最大空闲连接数
```

## MySQL语句

```go
sqlStr := `select id, name, gender ,hobby from user_infos where id=?;`
sqlStr := "select id, name, gender ,hobby from user_infos where id>?"
sqlStr := "insert into user_infos(name,gender,hobby) values(?,?,?)"
sqlStr := "update user_infos set hobby = ?,name=? where id = ?"
sqlStr := "delete from user_infos where id = ?"
```

## 查询

### 单条

```go
err := db.QueryRow(sqlStr, 2).Scan(&u.id, &u.name, &u.hobby, &u.gender)
```

### 多条

```go
sqlStr := "select id, name, gender ,hobby from user_infos where id>?"
rows, err := db.Query(sqlStr, 1)
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.hobby, &u.gender)
		fmt.Printf("id:%d name:%s hobby:%s gender:%s \n", u.id, u.name, u.gender, u.hobby)
	}
```

## 插入、更新、删除

```go
ret, err := db.Exec(sqlStr, "年轻人", "男", "练死劲儿")
n, err = ret.RowsAffected() // 操作影响的行数
```

## 预处理

```
stmt, err := db.Prepare(sqlStr)
defer stmt.Close() //关闭命令
rows, err := stmt.Query(0)
defer rows.Close() //关闭数据
```

## 事务

```
tx, err := db.Begin() // 开启事务
tx.Rollback() // 操作失败则回滚事务
如果两次影响的行数都为1，则提交
tx.Commit() // 提交事务
```

# 小清单项目



```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

type ToDo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

var db *gorm.DB

func initMysql() (err error) {
	dsn := "root:961024@tcp(localhost:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	//连接数据库
	err := initMysql()
	if err != nil {
		fmt.Println("init mysql failed!,err:", err)
		return
	}
//模型绑定
//新建表todos
db.AutoMigrate(&ToDo{})

// 中间件
r := gin.Default()
//加载静态文件
r.Static("/static", "F:\\goland\\go_project\\go_web\\websrc\\web_25\\static")
//加载文件
r.LoadHTMLGlob("F:\\goland\\go_project\\go_web\\websrc\\web_25\\tmlplates/*")

r.GET("/", func(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
})
v1Group := r.Group("v1")
{
	v1Group.POST("/todo", func(c *gin.Context) {
		var todo ToDo
		// 1、 从请求中把数据拿出来
		c.BindJSON(&todo)

		//2、存入数据
		err := db.Create(&todo).Error
		if err != nil {
			fmt.Println("create err")
			return
		} else {
			//3、返回响应
			c.JSON(http.StatusOK, gin.H{
				"msg":    "2000",
				"status": "ok",
			})
		}
	})

	v1Group.GET("/todo", func(c *gin.Context) {
		var todo []ToDo
		db.Find(&todo)
		//显示数据
		c.JSON(http.StatusOK, todo)
	})
	v1Group.PUT("/todo/:id", func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"errors": "无效的id",
			})
			return
		}
		var todo ToDo
		//******************************************************//
		// 查询要修改的id，这个非常重要
		//如果没有这一项的话无法匹配到要修改的项，则会在每次更新的时候save两条空数据
		//First(&todo)而不是First(&ToDo)，不然会更改整个结构体，同样增加两条新数据
		err := db.Where("id = ?", id).First(&todo).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"errors": "无效的id",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":    "2000",
				"status": "ok",
			})
		}

		// 从请求中将数据拿出来存入结构体todo中
		c.BindJSON(&todo)
		// 更新数据
		err = db.Save(&todo).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, todo)
		}
	})
	v1Group.DELETE("/todo/:id", func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"errors": "无效的id",
			})
			return
		}
		err := db.Where("id = ? ", id).Delete(ToDo{}).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{id: "deleted"})
		}
	})
}
r.Run()
}
```
# Redis--开源的内存数据库



# zap日志库

## **基础**

**日志记录器：**

- 能够将事件**记录到文件**中，而不是应用程序控制台。
- 日志切割-能够根据文件大小、时间或间隔等来**切割**日志文件。
- 支持不同的**日志级别**。例如INFO，DEBUG，ERROR等。
- 能够**打印基本信息**，如调用文件/函数名和行号，日志时间等。

```defer.sync```：在程序退出前将缓冲区中的日志刷到磁盘上。

**两种类型的日志记录器：**

`SugaredLogger` 在性能很好但不是很关键的上下文中。



`Logger` 在每一微秒和每一次内存分配都很重要的上下文中。

- 通过调用`zap.NewProduction()`/`zap.NewDevelopment()`或者`zap.Example()`创建一个Logger。

- 唯一的区别在于它将**记录的信息不同**。
- 通过Logger调用Info/Error等。
- 默认情况下日志都会**打印到**应用程序的**console界面**。

**1、将日志写入文件中。**

使用`zap.New(…)`方法来手动传递所有配置。

```go
zapcore.Core`需要三个配置——`Encoder`，`WriteSyncer`，`LogLevel
```

**Encoder**:编码器(如何写入日志)。

**WriterSyncer** ：指定日志将写到哪里去。

**Log Level**：哪种级别的日志将被写入。



## **接收gin框架默认的日志并配置日志归档**



使用`gin.Default()`的同时是用到了gin框架内的**两个默认中间件**`Logger()`和`Recovery()`。

`Logger()`是把gin框架本身的日志输出到标准输出（我们本地开发调试时在终端输出的那些日志就是它的功劳）。

`Recovery()`是在程序出现panic的时候恢复现场并写入500响应的。

```go
r := gin.New() // gin.Default() 代替注册中间件
r.Use(GinLogger(logger), GinRecovery(logger, true))
```

# Viper--配置管理神器

处理所有类型的配置需求和格式。

- 从`JSON`、`TOML`、`YAML`、`HCL`、`envfile`和`Java properties`格式的配置文件读取配置信息
- 实时监控和重新读取配置文件（可选）
- 从环境变量中读取
- 从远程配置系统（etcd或Consul）读取并监控配置变化
- 从命令行参数读取配置
- 从buffer读取配置
- 显式配置值
- 设置默认值

Viper会按照下面的优先级。每个项目的优先级都高于它下面的项目:

- 显示调用`Set`设置值
- 命令行参数（flag）
- 环境变量
- 配置文件
- key/value存储
- 默认值

# 大型Web项目分层

## MVC 模式

MVC 模式代表 Model-View-Controller（模型-视图-控制器） 模式。这种模式用于应用程序的分层开发。

- **Model（模型）** - 模型代表一个存取数据的对象或 JAVA POJO。它也可以带有逻辑，在数据变化时更新控制器。
- **View（视图）** - 视图代表模型包含的数据的可视化。
- **Controller（控制器）** - 控制器作用于模型和视图上。它控制数据流向模型对象，并在数据变化时更新视图。它使视图与模型分离开。

<img src="https://www.runoob.com/wp-content/uploads/2014/08/1200px-ModelViewControllerDiagram2.svg_.png" alt="img" style="zoom:60%;" />

## 当前的分层

<img src="C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20230704103302249.png" alt="image-20230704103302249" style="zoom:80%;" />

## 通用脚手架

### 主函数：```main.go```

```go
func main() {
   //1、加载配置信息 viper
   if err := settings.Init(); err != nil {
      fmt.Printf("init settings failed,err:%v/n", err)
      return
   }
   //2、初始化日志 zap
   if err := logger.InitLogger(); err != nil {
      fmt.Printf("init logger failed,err:%v/n", err)
      return
   }
   //在程序退出前将缓冲区中的日志刷到磁盘上。
   defer zap.L().Sync()
   zap.L().Debug("logger init success...")
   //3、初始化MySQL连接
   if err := mysql.Init(); err != nil {
      fmt.Printf("init mysql failed,err:%v/n", err)
      return
   }
   defer mysql.Close()
   //4、初始化Redis连接
   if err := redis.InitClient(); err != nil {
      fmt.Printf("init redis failed,err:%v/n", err)
      return
   }
   fmt.Println("init redis succes")
   defer redis.Close()
   //5、注册路由
   r := routes.Setup()
   //6、启动服务
   srv := &http.Server{
      Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
      Handler: r,
   }

   go func() {
      // 开启一个goroutine启动服务
      if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
         log.Fatalf("listen: %s\n", err)
      }
   }()

   // 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
   quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
   // kill 默认会发送 syscall.SIGTERM 信号
   // kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
   // kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
   // signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
   <-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
   zap.L().Info("Shutdown Server ...")
   // 创建一个5秒超时的context
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   // 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
   if err := srv.Shutdown(ctx); err != nil {
      zap.L().Fatal("Server Shutdown: ", zap.Error(err))
   }
   zap.L().Info("Server exiting")
}
```

### 加载配置

创建```settings```文件夹，使用```viper```管理配置信息。创建```config.yaml```配置文件。

```go
func Init() (err error) {
   viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
   viper.SetConfigType("yaml")   // 指定配置文件类型
   viper.AddConfigPath(".")      // 指定查找配置文件的路径（这里使用相对路径）
   //viper.SetConfigFile("./conf/config.yaml") // 指定配置文件路径
   err = viper.ReadInConfig() // 读取配置信息
   if err != nil {            // 读取配置信息失败
      // 读取配置信息失败
      fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
      return
   }

   // 监控配置文件变化
   viper.WatchConfig()
   viper.OnConfigChange(func(e fsnotify.Event) {
      // 配置文件发生变更之后会调用的回调函数
      fmt.Println("Config file changed")
   })
   return
}
```

```yaml
app:
  name: "web_app"
  mode: "dev"
  port: 8081
log:
  level: "debug"
  filename: "web_bluebell.log"
  max_size: 200
  max_age: 30
  max_backups: 7
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "961024"
  dbname: "bluebell"
  max_open_conns: 200
  max_idle_conns: 50
redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0
  pool_size: 100
```

### 初始化日志

创建```logger```文件夹。zap来接收并记录gin框架默认的日志并配置日志归档

```go
// InitLogger 初始化Logger
func InitLogger() (err error) {
   writeSyncer := getLogWriter(
      viper.GetString("log.filename"),
      viper.GetInt("log.max_size"),
      viper.GetInt("log.max_backups"),
      viper.GetInt("log.max_age"),
   )
   encoder := getEncoder()
   var l = new(zapcore.Level)
   err = l.UnmarshalText([]byte(viper.GetString("log.level")))
   if err != nil {
      return
   }
   core := zapcore.NewCore(encoder, writeSyncer, l)

   lg := zap.New(core, zap.AddCaller())
   // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
   zap.ReplaceGlobals(lg)   //***********重要************
   return
}

func getEncoder() zapcore.Encoder {
   encoderConfig := zap.NewProductionEncoderConfig()
   encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
   encoderConfig.TimeKey = "time"
   encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
   encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
   encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
   return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
   lumberJackLogger := &lumberjack.Logger{
      Filename:   filename,
      MaxSize:    maxSize,
      MaxBackups: maxBackup,
      MaxAge:     maxAge,
   }
   return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
   return func(c *gin.Context) {
      start := time.Now()
      path := c.Request.URL.Path
      query := c.Request.URL.RawQuery
      c.Next()

      cost := time.Since(start)
      zap.L().Info(path,
         zap.Int("status", c.Writer.Status()),
         zap.String("method", c.Request.Method),
         zap.String("path", path),
         zap.String("query", query),
         zap.String("ip", c.ClientIP()),
         zap.String("user-agent", c.Request.UserAgent()),
         zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
         zap.Duration("cost", cost),
      )
   }
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
   return func(c *gin.Context) {
      defer func() {
         if err := recover(); err != nil {
            // Check for a broken connection, as it is not really a
            // condition that warrants a panic stack trace.
            var brokenPipe bool
            if ne, ok := err.(*net.OpError); ok {
               if se, ok := ne.Err.(*os.SyscallError); ok {
                  if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
                     brokenPipe = true
                  }
               }
            }

            httpRequest, _ := httputil.DumpRequest(c.Request, false)
            if brokenPipe {
               zap.L().Error(c.Request.URL.Path,
                  zap.Any("error", err),
                  zap.String("request", string(httpRequest)),
               )
               // If the connection is dead, we can't write a status to it.
               c.Error(err.(error)) // nolint: errcheck
               c.Abort()
               return
            }

            if stack {
               zap.L().Error("[Recovery from panic]",
                  zap.Any("error", err),
                  zap.String("request", string(httpRequest)),
                  zap.String("stack", string(debug.Stack())),
               )
            } else {
               zap.L().Error("[Recovery from panic]",
                  zap.Any("error", err),
                  zap.String("request", string(httpRequest)),
               )
            }
            c.AbortWithStatus(http.StatusInternalServerError)
         }
      }()
      c.Next()
   }
}
```

### 初始化MySQL连接

```go
import (
   "fmt"

   _ "github.com/go-sql-driver/mysql"
   "github.com/jmoiron/sqlx"
   "go.uber.org/zap"

   "github.com/spf13/viper"
)

// 定义一个全局对象db
var db *sqlx.DB

func Init() (err error) {
   dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
      viper.GetString("mysql.user"),
      viper.GetString("mysql.password"),
      viper.GetString("mysql.host"),
      viper.GetInt("mysql.port"),
      viper.GetString("mysql.dbname"),
   )
   // 也可以使用MustConnect连接不成功就panic
   db, err = sqlx.Connect("mysql", dsn)
   if err != nil {
      zap.L().Error("connect DB failed", zap.Error(err))
      return
   }
   db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
   db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
   return
}

func Close() {
   _ = db.Close()
}
```

### 初始化redis连接

```go
var rdb *redis.Client

func InitClient() (err error) {
   rdb = redis.NewClient(&redis.Options{
      Addr: fmt.Sprintf("%s:%d",
         viper.GetString("redis.host"),
         viper.GetInt("redis.port"),
      ),
      Password: viper.GetString("redis.password"), // no password set
      DB:       viper.GetInt("redis.db"),          // use default DB
      PoolSize: viper.GetInt("redis.pool_size"),   // 连接池大小
   })

   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()

   _, err = rdb.Ping(ctx).Result()
   return err
}

func Close() {
   _ = rdb.Close()
}
```

### 注册路由

新建一个```routes```

```go
func Setup() *gin.Engine {
   r := gin.New()

   r.Use(logger.GinLogger(), logger.GinRecovery(true))
   r.GET("/", func(c *gin.Context) {
      c.String(http.StatusOK, "ok")
   })
   return r
}
```

### 启动服务

优雅关机(在主函数zhong)

## web2 改进

```settings```将所有的配置信息都制作一个结构体变量，并且反序列化到结构体中，这样每次直接从结构体中取值即可。

定义一个全局的变量，将配置信息反序列化到 Conf 变量中。

直接在其它文件夹中调用这个结构体中的变量即可，避免频繁的调用配置信息使用```viper```库

```go
//直接寻找
viper.SetConfigFile("./conf/config.yaml") // 指定配置文件路径 可以代替下面两行。
//通过路径寻找这个名称的配置文件
viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
viper.SetConfigType("yaml")   // 指定配置文件类型(专用于从远程获取文件的配置信息时指定配置）etcd 
viper.AddConfigPath(".")      // 指定查找配置文件的路径（这里使用相对路径）
```

## os库

```os.Args[0]```指的是当前路径

```os.Args[1]```指的是当前目录下的文件夹

```go
if len(os.Args) < 2 {
   fmt.Println("need config file.eg: bluebell config.yaml")
   return
}
// 加载配置
if err := setting.Init(os.Args[1]); err != nil {
	fmt.Printf("load config failed, err:%v\n", err)
	return
}
```

## flag库 在任何路径下执行文件

[Go语言标准库flag基本使用 | 李文周的博客 (liwenzhou.com)](https://www.liwenzhou.com/posts/Go/flag/)

# bulebell项目实战

## 30用户表结构设计

用户表结构设计。

使用goland连接数据库。

构建数据表。在models文件夹下的creat_table里面

## 31雪花算法

不使用自增id作为用户id的原因：

1、防止其它用户可以窃取到数据库信息。

2、数据库分区分库可能会有id重复。

![img](F:\goland\go_project\go_Web81\goWeb\bluebell\picture\img.png)

![img_1](F:\goland\go_project\go_Web81\goWeb\bluebell\picture\img_1.png)

## 32 注册业务流程

![image-20230703154907993](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20230703154907993.png)

## 33注册业务流程

### ```routes```文件夹下添加注册post

```
// 注册
r.POST("/signup", controllers.SignUpHandler)
```

### ```controllers```文件夹下创建```user.go```文件

```go
// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
   //1、获取参数和参数校验
   //2、业务处理
   logic.SignUp()
   //3、返回响应
   c.JSON(http.StatusOK, "ok")
}
```

### ```logic```文件夹下创建```user.go```文件

```go
// 存放业务逻辑的代码

func SignUp() {
   //1、判断用户是否存在
   mysql.QueryUserByUsername()
   //2、生成UID
   snowflake.GenID()
   //3、保存进数据库
   mysql.InsertUser()
}
```

### ```mysql```文件夹下创建```user.go```

```go
// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

func QueryUserByUsername() {
}

func InsertUser() {
}
```

## 33请求参数的获取与校验

### ```params``` 文件夹 创建请求的参数结构体

```go
//定义请求的参数结构体

type ParmSignUp struct {
   Username   string `json:"username" binding:"required"`
   Password   string `json:"password" binding:"required"`
   RePassword string `json:"re_password" binding:"required,eqfield=Password"`
   // "required,eqfield=Password" 判断re_password == password
}
```

### ``` controllers```文件夹下的```user``` 

```go
// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
   //1、获取参数和参数校验
   p := new(models.ParmSignUp)
   if err := c.ShouldBindJSON(p); err != nil {
      //请求参数、有误直接返回响应
      // 记录日志
      zap.L().Error("SignUp with invalid param", zap.Error(err))
      // 判断err是不是validator.ValidationErrors 类型
      errs, ok := err.(validator.ValidationErrors)
      if !ok {
         c.JSON(http.StatusOK, gin.H{
            "msg": err.Error(),
         })
         return
      }
      c.JSON(http.StatusOK, gin.H{
         //"msg": "请求参数有误",
         //查看哪里有误
         //"msg": err.Error(),
         // 翻译
         // 移除不相干的
         "msg": removeTopStruct(errs.Translate(trans)),
      })
      return
   }

   //` `中使用 binding:"required" 替换以下功能
   // 手动对请求参数进行详细的业务规则校验
   //if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
   // zap.L().Error("SignUp with invalid param")
   // c.JSON(http.StatusOK, gin.H{
   //    "msg": "请求参数有误",
   // })
   // return
   //}

   fmt.Println(p)
   //2、业务处理
   if err := logic.SignUp(p); err != nil {
      c.JSON(http.StatusOK, gin.H{
         "msg": "注册失败",
      })
      return
   }
   //3、返回响应
   c.JSON(http.StatusOK, gin.H{
      "msg": "success",
   })
}
```

### ```logic```文件夹下的```user``` 注册

```go
// 存放业务逻辑的代码

func SignUp(p *models.ParmSignUp) (err error) {
   //1、判断用户是否存在
   err = mysql.CheckUserExist(p.Username) //绑定到结构体中了 所以直接调用mysql中语句进行检查
   if err != nil {
      // 数据库查询出错
      return err
   }
   //2、生成UID
   userID := snowflake.GenID()
   // 构造一个user实例
    // **********
	// 将网页中拿到的用户名密码反序列化到 models.ParmSignUp 结构体中
	// 再将 models.ParmSignUp 中的值赋给定义的 models.User 的结构体中
	// **********
   user := &models.User{
      UserID:   userID,
      Username: p.Username,
      Password: p.Password,
   }
   //3、保存进数据库
   return mysql.InsertUser(user)
}
```

### ```mysql```文件夹下的```user```

```go
const serect = "liwenzhou.com"

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

// 判断用户存不存在
func CheckUserExist(username string) (err error) {
   sqlStr := `select count(user_id) from user where username = ?`
   var count int
   if err := db.Get(&count, sqlStr, username); err != nil {
      return err
   }
   if count > 0 {
      return errors.New("用户已存在")
   }
   return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
   // 对密码进行加密
   user.Password = encryPassword(user.Password)
   // 执行SQL语句入库
   sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
   _, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
   return
}

// ************对密码进行加密*************
func encryPassword(oPassword string) string {
   h := md5.New()
   h.Write([]byte(serect))
   return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
```

### 附录：参数校验和翻译器

```controllers```文件夹下的```validator```

```go
// 翻译器
import (
   "fmt"
   "reflect"
   "strings"

   "github.com/gin-gonic/gin/binding"
   "github.com/go-playground/locales/en"
   "github.com/go-playground/locales/zh"
   ut "github.com/go-playground/universal-translator"
   "github.com/go-playground/validator/v10"
   enTranslations "github.com/go-playground/validator/v10/translations/en"
   zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器T
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
   // 修改gin框架中的Validator引擎属性，实现自定制
   if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		
      // ******Repassword  --> re_password  ***********
      // 注册一个获取json tag的自定义方法
      v.RegisterTagNameFunc(func(fld reflect.StructField) string {
         name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
         if name == "-" {
            return ""
         }
         return name
      })
	  // ***************************************
      zhT := zh.New() // 中文翻译器
      enT := en.New() // 英文翻译器

      // 第一个参数是备用（fallback）的语言环境
      // 后面的参数是应该支持的语言环境（支持多个）
      // uni := ut.New(zhT, zhT) 也是可以的
      uni := ut.New(enT, zhT, enT)

      // locale 通常取决于 http 请求头的 'Accept-Language'
      var ok bool
      // 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
      trans, ok = uni.GetTranslator(locale)
      if !ok {
         return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
      }

      // 注册翻译器
      switch locale {
      case "en":
         err = enTranslations.RegisterDefaultTranslations(v, trans)
      case "zh":
         err = zhTranslations.RegisterDefaultTranslations(v, trans)
      default:
         err = enTranslations.RegisterDefaultTranslations(v, trans)
      }
      return
   }
   return
}

type SignUpParam struct {
   Age        uint8  `json:"age" binding:"gte=1,lte=130"`
   Name       string `json:"name" binding:"required"`
   Email      string `json:"email" binding:"required,email"`
   Password   string `json:"password" binding:"required"`
   RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// removeTopStruct  去除提示信息中的结构体名称
// "ParmSignUp.re_password": "re_password为必填字段"
//
// re_password": "re_password为必填字段"
func removeTopStruct(fields map[string]string) map[string]string {
   res := map[string]string{}
   for field, err := range fields {
      res[field[strings.Index(field, ".")+1:]] = err
   }
   return res
}
```

## 36使用mode控制日志输出位置

```logger```下的```logger```文件

```go
func InitLogger(cfg *settings.LogConfig, mode string) (err error) {
    ...
var core zapcore.Core
if mode == "dev" {
   // 定义 进入开发模式，日志输出到终端 的变量
   consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
   //两种情况
    core = zapcore.NewTee(
      //输出到文件
      zapcore.NewCore(encoder, writeSyncer, l),
      //输出到终端
      zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
   )
} else {
   core = zapcore.NewCore(encoder, writeSyncer, l)
}
    ...
}
```

如果是发布者 ```release```模式，则不会输出到终端以下信息

```go
[GIN-debug] POST   /signup                   --> bluebell.com/bluebell/controllers.SignUpHandler (3 handlers)
[GIN-debug] GET    /                         --> bluebell.com/bluebell/routes.Setup.func1 (3 handlers)
```

```routes```文件下的```route```

```go
func Setup(mode string) *gin.Engine {
   if mode == gin.ReleaseMode {
      gin.SetMode(gin.ReleaseMode) // gin设置程发布者模式
   }
	...
}
```

```main.go```主函数中传入```settings.Conf.Mode```

## 37登录功能实现

### ```routes.routes```登录请求

```go
// 登录
r.POST("/login", controllers.LoginHandler)
```

### ```controllers.user```登录请求响应函数

获取参数和参数校验并 记录错误信息

业务处理

返回响应

```go
func LoginHandler(c *gin.Context) {
   //1、获取参数和参数校验
   p := new(models.ParmLogin)
   if err := c.ShouldBindJSON(p); err != nil {
      //请求参数、有误直接返回响应
      // 记录日志
      zap.L().Error("Login with invalid param", zap.Error(err))
      // 判断err是不是validator.ValidationErrors 类型
       // 获得应用程序的错误消息的所有信息
      errs, ok := err.(validator.ValidationErrors)
      if !ok {
         c.JSON(http.StatusOK, gin.H{
            "msg": err.Error(),
         })
         return
      }
      c.JSON(http.StatusOK, gin.H{
         "msg": removeTopStruct(errs.Translate(trans)),
      })
      return
   }
   //2、业务处理
   if err := logic.Login(p); err != nil {
      zap.L().Error("logic.Login failed", zap.String("username:", p.Username), zap.Error(err))
      c.JSON(http.StatusOK, gin.H{
         "msg": "用户名或密码错误",
      })
      return
   }
   //3、返回响应
   c.JSON(http.StatusOK, gin.H{
      "msg": "success",
   })
}
```

### ```logic.user```登录业务处理

获取到网页中请求的用户名和密码 加入到结构体中

判断与数据库中的是否相等

```go
func Login(p *models.ParmLogin) (err error) {
   user := &models.User{
      Username: p.Username,
      Password: p.Password,
   }
   if err := mysql.Login(user); err != nil {
      return err
   }
   return
}
```

### ```mysql.user```数据库登录时的操作

使用```sql```语句查看数据库中是否存在这个用户

判断密码是否正确

密码加密

```go
var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)
// Login
func Login(user *models.User) (err error) {
   oPassword := user.Password // 用户登录的密码
   sqlStr := `select user_id, username, password from user where username=?`
   err = db.Get(user, sqlStr, user.Username)
   // var ErrNoRows = errors.New("sql: no rows in result set")
   if err == sql.ErrNoRows {
      // 打印用户不存在 已经定义为全局变量
      return ErrorUserNotExist
   }
   if err != nil {
      // 查询数据库失败
      return err
   }
   // 判断密码是否正确
   password := encryptPassword(oPassword)
   if password != user.Password {
      // 如果密码不相等，打印用户名或密码错误 已经定义为全局变量
      return ErrorInvalidPassword
   }
   return
}
// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(serect))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
```

