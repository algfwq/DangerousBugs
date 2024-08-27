package main

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sync"
)

// 数据库操作
var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// 处理报错
var err error

// websocket连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
