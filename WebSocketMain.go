package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func WebSocketMain(w http.ResponseWriter, r *http.Request) {
	conn, errGrader := upgrader.Upgrade(w, r, nil)
	if errGrader != nil {
		log.Println("连接升级失败：", errGrader)
		return
	} else {
		// 创建一个包含所需信息的结构体
		response := map[string]string{"mode": "stop_loading"}

		// 将结构体编码为JSON字节数组
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("JSON编码失败：", err)
			return
		}

		// 发送JSON字节数组
		errSend := conn.WriteMessage(websocket.TextMessage, jsonResponse)
		if errSend != nil {
			log.Println("发送数据失败：", errSend)
			return
		}
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("关闭连接失败: %v", err)
		}
	}()

	for {
		messageType, p, errRead := conn.ReadMessage()
		if errRead != nil {
			log.Println("读取消息失败/连接关闭", errRead)
			return
		}
		log.Println("收到消息:", string(p))

		errSend := conn.WriteMessage(messageType, []byte("Hello, world!"))
		if errSend != nil {
			log.Println("发送数据失败：", errSend)
			return
		}
	}
}
