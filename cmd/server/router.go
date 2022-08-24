package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/riverchu/pkg/log"
)

func registerRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ws", WSHanlder(handle))
}

var upgrader = new(websocket.Upgrader)

func WSHanlder(handle func([]byte) []byte) gin.HandlerFunc {
	return WSHanlderWithUpgrader(upgrader, handle)
}

func WSHanlderWithUpgrader(upgrader *websocket.Upgrader, handle func([]byte) []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()
		for {
			mt, msg, err := ws.ReadMessage()
			if err != nil {
				log.Warn("read message fail: %s", err)
				break
			}

			err = ws.WriteMessage(mt, handle(msg))
			if err != nil {
				log.Warn("write message fail: %s", err)
				break
			}
		}
	}

}

func handle(msg []byte) []byte {
	switch string(msg) {
	case "ping":
		return []byte("pong")
	default:
		return msg
	}
}
