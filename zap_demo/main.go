package main

import (
	"go.uber.org/zap"
	"net/http"
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

func main() {
	InitLogger()
	defer logger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.bilibili.com")
}

func InitLogger() {
	logger, _ = zap.NewProduction()
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		logger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}
