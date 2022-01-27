package service

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"xivbot/models"
	"xivbot/util"

	"github.com/google/uuid"
)

func init() {
	handlers = append(handlers, SenpaiHandler)
}

func SenpaiHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	if strings.Contains(msg.Message, "114514") || strings.Contains(msg.Message, "1919810") {
		senpai := models.Senpai{}
		senpais, err := senpai.Find()
		if err != nil {
			log.Println(err)
		}
		rand.Seed(time.Now().Unix())
		response, _ := ImageResponse(1, senpais)
		SendGroupMsg(response, msg.GroupID)
	}
	if strings.Contains(msg.Message, "添加臭图") {
		events, err := util.ParseEvent(msg.Message)
		if err != nil {
			log.Println(err)
		}
		for _, v := range events {
			if vv, okk := v["CQ"]; okk && vv == "image" {
				link := v["url"].(string)
				id := uuid.New().String()
				file := id + ".jpg"
				out, err := os.Create(Path + "senpai/" + file)
				if err != nil {
					log.Println(err)
					return
				}
				defer out.Close()
				resp, err := http.Get(link)
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					log.Println(err)
					return
				}
				senpai := models.Senpai{
					Src: Url + "/senpai/" + file,
				}
				err = senpai.Insert()
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
		response := Reply(msg.MessageID, fmt.Sprintf("添加%d张臭图成功", len(events)))
		SendGroupMsg(response, msg.GroupID)
	}
	m.Unlock()
}
