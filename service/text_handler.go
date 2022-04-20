package service

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"xivbot/util"

	"github.com/google/uuid"
)

func init() {
	handlers = append(handlers, TestHandler)
	handlers = append(handlers, KeywordHandler)
}

func KeywordHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	if ok, _ := regexp.MatchString(`^(?:!|！|\/|-|--)help$`, msg.Message); ok {
		response := Help
		SendGroupMsg(response, msg.GroupID)
		return
	}
	if ok, _ := regexp.MatchString(`^\#查询你群男同浓度$`, msg.Message); ok {
		response := Reply(msg.MessageID, "114514%")
		SendGroupMsg(response, msg.GroupID)
		return
	}

	// if strings.Contains(msg.Message, "原神") {
	// 	response := Reply(msg.MessageID, "原批爬")
	// 	SendGroupMsg(response, msg.GroupID)
	// 	return
	// }

	if strings.Contains(msg.Message, "granbluefantasy.jp") {
		response := Reply(msg.MessageID, "骑空士爬")
		SendGroupMsg(response, msg.GroupID)
		return
	}

	m.Unlock()
}

func TestHandler(msg Request) {
	if ok, _ := regexp.MatchString(`^(?:--测试|--test)`, msg.Message); ok {
		events, err := util.ParseEvent(msg.Message)
		if err != nil {
			log.Println(err)
		}
		for _, v := range events {
			if vv, okk := v["CQ"]; okk && vv == "image" {
				link := v["url"].(string)
				id := uuid.New().String()
				file := id + ".jpg"
				out, err := os.Create(Path + "test/" + file)
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
			}
		}
	}
}
