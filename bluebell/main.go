package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bluebell.com/bluebell/controllers"

	"bluebell.com/bluebell/pkg/snowflake"

	"fmt"

	"bluebell.com/bluebell/dao/mysql"
	"bluebell.com/bluebell/dao/redis"
	"bluebell.com/bluebell/routes"

	"bluebell.com/bluebell/logger"
	"go.uber.org/zap"

	"bluebell.com/bluebell/settings"
)

// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
func main() {
	//if len(os.Args) < 2 {
	//	fmt.Println("need config file.eg: bluebell config.yaml")
	//	return
	//}
	//1、加载配置信息 viper
	//运行的目录在bluebell下，故需要去它的下一级目录找 os.Args[1]
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v/n", err)
		return
	}
	//2、初始化日志 zap
	if err := logger.InitLogger(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed,err:%v/n", err)
		return
	}
	//在程序退出前将缓冲区中的日志刷到磁盘上。
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	//3、初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
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
	// 雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	fmt.Println(snowflake.GenID())

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init calidator trans failed,err:%v\n", err)
		return
	}
	//5、注册路由
	r := routes.Setup(settings.Conf.Mode)
	//6、启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
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
