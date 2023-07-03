package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// // 设置日志记录器，设置日志的保存位置
//
//	func SetupLogger() {
//		logFileLocation, _ := os.OpenFile("F:\\goland\\go_project\\go_Web81\\goWeb\\zap_demo/test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//		log.SetOutput(logFileLocation)
//	}
//
// // 使用Logger
//
//	func simpleHttpGet(url string) {
//		resp, err := http.Get(url)
//		if err != nil {
//			log.Printf("Error fetching url %s : %s", url, err.Error())
//		} else {
//			log.Printf("Status Code for %s : %s", url, resp.Status)
//			resp.Body.Close()
//		}
//	}
//
//	func main() {
//		SetupLogger()
//		simpleHttpGet("www.baidu.com")
//		simpleHttpGet("http://www.bilibili.com")
//	}
var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	//defer logger.Sync()
	//for i := 0; i < 100000; i++ {
	//	logger.Info("test for log rotate...")
	//}
	//simpleHttpGet("www.baidu.com")
	//simpleHttpGet("http://www.bilibili.com")
	r := gin.New() // gin.Default() 代替注册中间件
	r.Use(GinLogger(logger), GinRecovery(logger, true))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello q1mi!")
	})
	r.Run()
}

func InitLogger() {
	//  logger, _ := zap.NewProduction()
	//	sugarLogger = logger.Sugar()

	//Log Level：哪种级别的日志将被写入。
	//我们将修改上述部分中的Logger代码，并重写InitLogger()方法
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	// zapcore.Core`需要三个配置——`Encoder`，`WriteSyncer`，`LogLevel
	//**Encoder**:编码器(如何写入日志)。
	//**WriterSyncer** ：指定日志将写到哪里去。
	//**Log Level**：哪种级别的日志将被写入。
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	//使用zap.new()手动传输所有配置
	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

// Encoder:编码器(如何写入日志)
func getEncoder() zapcore.Encoder {
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// WriterSyncer ：指定日志将写到哪里去
func getLogWriter() zapcore.WriteSyncer {
	//file, _ := os.OpenFile("F:\\goland\\go_project\\go_Web81\\goWeb\\zap_demo\\test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	//return zapcore.AddSync(file)
	// 日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     // M
		MaxBackups: 5,     // 最大备份数量
		MaxAge:     30,    // 最大备份天数
		Compress:   false, // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}
func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
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

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
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
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
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
