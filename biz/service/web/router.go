package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/riverchu/pkg/websocket"

	"github.com/riverchu/worker/biz/handler"
)

func registerRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/ws", ws.WSHanlder(handler.WSHandle))
}
