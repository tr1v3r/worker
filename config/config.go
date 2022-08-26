package config

import (
	"sync"

	"github.com/gorilla/websocket"
)

var mu sync.RWMutex
var wsConnMap = map[ /*token*/ string]*websocket.Conn{}

func GetWSConn(token string) *websocket.Conn {
	mu.RLock()
	defer mu.RUnlock()
	return wsConnMap[token]
}

func SetWSConn(token string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	wsConnMap[token] = conn
}
