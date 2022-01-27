package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

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

func ImageResponse(count int, l []string) (response []CQMessage, err error) {
	if count > 5 {
		err = errors.New("仅支持最多5张")
		return
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		response = append(response, CQMessage{
			Type: "image",
			Data: CQImage{
				File: l[rand.Intn(len(l))],
			},
		})
	}
	return
}

func Reply(msgID int32, text string) []CQMessage {
	return []CQMessage{
		{
			Type: "reply",
			Data: CQReply{
				ID: msgID,
			},
		},
		{
			Type: "text",
			Data: CQText{
				Text: text,
			},
		},
	}
}
