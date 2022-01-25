package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
		log.Println(string(r))
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

func SendGroupMsg(msg interface{}, groupID int64) {
	api := "/send_group_msg"
	url := Url + Port
	response := Message{
		GroupID: groupID,
		Message: msg,
	}
	r, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
	}
	_, err = http.Post(url+api, "application/json", bytes.NewBuffer(r))
	if err != nil {
		log.Println(err)
	}
}
