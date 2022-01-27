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

var msgMap map[int32]int

func Run(port string) {
	msgMap = make(map[int32]int)
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
		if msg.GroupID == 0 {
			return
		}
		m := new(sync.Mutex)
		m.Lock()
		if _, ok := msgMap[msg.MessageID]; ok {
			return
		}
		msgMap[msg.MessageID] = 1
		for _, v := range handlers {
			go v(msg)
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
