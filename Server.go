package main

import (
	"log"
	"net/http"
	"os"
)

var logFileServer *os.File

func init() {
	// 初始化日志文件
	logFileServer, err = os.OpenFile("log_server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}

	// 设置日志输出到文件
	log.SetOutput(logFileServer)

	// 记录日志文件初始化成功的消息
	log.Println("日志文件初始化成功")
}

func main() {
	//结束时执行
	defer mainExit()

	// 创建HTTP服务器
	http.HandleFunc("/main/", WebSocketMain)
	http.HandleFunc("/download/", downloadHandler)
	log.Println("服务器开启在端口:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainExit() {
	// 确保日志文件在程序结束时关闭
	// 注意：因为 init() 函数没有直接的方法来在程序结束时执行代码，
	// 通常需要在main()中或程序的其他适当位置显式调用关闭文件的操作。
	if err := logFileServer.Close(); err != nil {
		log.Printf("关闭logFile失败: %v", err)
	}
}
