package main

import (
	"fmt"

	"bluebell.com/bluebell/logger"
	"go.uber.org/zap"

	"bluebell.com/bluebell/settings"
)

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
	//4、初始化Redis连接
	//5、注册路由
	//6、启动服务
}
