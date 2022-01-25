package service

import "github.com/gin-gonic/gin"

var handlers []func(Request)

func Run(port string) {
	engine := gin.Default()
	engine.POST("/", Handler())
	engine.Run(port)
}

func Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
