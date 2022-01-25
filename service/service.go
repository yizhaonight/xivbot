package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

var handlers []func(Request)

func Run(port string) {
	engine := gin.Default()
	engine.POST("/", Handler())
	engine.Run(port)
}

func Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg Request
		r, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(r, &msg)
		if err != nil {
			log.Println(err)
		}
		m := new(sync.Mutex)
		m.Lock()
		for _, v := range handlers {
			v(msg)
		}
		m.Unlock()
	}
}
