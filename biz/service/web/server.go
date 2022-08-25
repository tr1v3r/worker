package web

import (
	"github.com/gin-gonic/gin"
	"github.com/riverchu/pkg/log"

	"github.com/riverchu/worker/config"
)

func Serve() {
	r := gin.Default()

	registerRouter(r)

	err := r.Run(":" + config.WebServerPort())
	if err != nil {
		log.Error("gin server stopped: %s", err)
	}
}
