package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var err error
var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
