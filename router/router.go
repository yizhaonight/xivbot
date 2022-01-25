package router

import "github.com/gin-gonic/gin"

type Handler interface {
	Handle() gin.HandlerFunc
}

func POST(r *gin.Engine, url string, handlers []Handler) {
	var h []gin.HandlerFunc
	for _, v := range handlers {
		h = append(h, v.Handle())
	}
	r.POST(url, h...)
}
