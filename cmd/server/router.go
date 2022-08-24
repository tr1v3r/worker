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
	r.GET("/ws", WSHanlder())
}

var upgrader = websocket.Upgrader{}

func WSHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Warn("read message fail: %s", err)
				break
			}

			var resp string
			switch string(message) {
			case "ping":
				resp = "pong"
			default:
				resp = string(message)
			}

			err = ws.WriteMessage(mt, []byte(resp))
			if err != nil {
				log.Warn("write message fail: %s", err)
				break
			}
		}
	}
}
