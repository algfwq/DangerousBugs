package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

var logFileServer *os.File

// GetMongoClient 返回全局MongoDB客户端实例
func GetMongoClient() *mongo.Client {
	clientOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

		// Connect to MongoDB
		var err error
		clientInstance, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal("连接数据库失败: ", err)
		}

		// Check the connection
		err = clientInstance.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal("数据库连接失败: ", err)
		}

		log.Println("Connected to MongoDB!")
	})
	return clientInstance
}

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
	// 初始化MongoDB连接
	GetMongoClient()

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

	// 关闭MongoDB连接
	if err := clientInstance.Disconnect(context.Background()); err != nil {
		log.Printf("关闭连接失败: %v", err)
	}
}
