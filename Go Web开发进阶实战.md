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

主函数：```main.go```



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
   zap.ReplaceGlobals(lg)
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

### 初始化redis连接

### 注册路由

### 启动服务



# bulebell项目实战

## 课时30

用户表结构设计。

使用goland连接数据库。

构建数据表。在models文件夹下的creat_table里面

## 课时31

不使用自增id作为用户id的原因：

1、防止其它用户可以窃取到数据库信息。

2、数据库分区的时候便于管理。

![img](F:\goland\go_project\go_Web81\goWeb\bluebell\picture\img.png)

![img_1](F:\goland\go_project\go_Web81\goWeb\bluebell\picture\img_1.png)

## 课时32 注册业务流程

![image-20230703154907993](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20230703154907993.png)