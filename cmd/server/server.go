// Package main server
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/riverchu/pkg/log"
)

func main() {
	r := gin.Default()

	registerRouter(r)

	err := r.Run(":7749")
	if err != nil {
		log.Error("gin server stopped: %s", err)
	}
}
