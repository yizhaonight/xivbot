package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"xivbot/models"
	"xivbot/util"

	"github.com/google/uuid"
)

func init() {
	handlers = append(handlers, TestHandler)
	handlers = append(handlers, KeywordHandler)
	handlers = append(handlers, PixivHandler)
	handlers = append(handlers, SenpaiHandler)
}

func KeywordHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	if ok, _ := regexp.MatchString(`^(?:!|！|\/|-|--)help$`, msg.Message); ok {
		response := Help
		SendGroupMsg(response, msg.GroupID)
	}

	var react = models.Reaction{GroupID: msg.GroupID}
	reacts, err := react.FindByGroupID()
	if err != nil {
		log.Println(err)
	}
	ok, _ := regexp.MatchString(`^有人说`, msg.Message)
	if (strings.Contains(msg.Message, "回复")) && ok {
		split := strings.Split(msg.Message, "有人说")
		condition := strings.Split(split[1], "回复")[0]
		reply := strings.Split(split[1], "回复")[1]
		for _, v := range reacts {
			if v.Word == condition {
				response := Reply(msg.MessageID, "规则已存在")
				SendGroupMsg(response, msg.GroupID)
				return
			}
		}
		r := models.Reaction{
			Word:     condition,
			Response: reply,
			GroupID:  msg.GroupID,
		}
		err := r.Insert()
		if err != nil {
			log.Println(err)
		}
		response := Reply(msg.MessageID, "添加规则成功")
		SendGroupMsg(response, msg.GroupID)
		return
	}

	if ok, _ := regexp.MatchString(`^查看规则$`, msg.Message); ok {
		var r string
		if len(reacts) == 0 {
			r = "当前没有生效的规则"
		} else {
			for i, v := range reacts {
				r += fmt.Sprintf("%d. 有人说%s回复%s\n", i+1, v.Word, v.Response)
			}
		}
		response := Reply(msg.MessageID, r)
		SendGroupMsg(response, msg.GroupID)
		return
	}

	if ok, _ := regexp.MatchString(`^删除规则`, msg.Message); ok {
		var r string
		split := strings.Split(msg.Message, "删除规则")
		condition := split[1]
		if condition == "" {
			return
		}
		c := 0
		for _, v := range reacts {
			if v.Word == condition {
				c += 1
			}
		}
		if c == 0 {
			r = "该规则不存在"
		} else {
			target := models.Reaction{
				Word:    condition,
				GroupID: msg.GroupID,
			}
			err := target.DeleteByWord()
			if err != nil {
				log.Println(err)
			}
			r = "删除成功"
		}
		response := Reply(msg.MessageID, r)
		SendGroupMsg(response, msg.GroupID)
		return
	}
	for _, v := range reacts {
		if strings.Contains(msg.Message, v.Word) {
			response := Reply(msg.MessageID, v.Response)
			SendGroupMsg(response, msg.GroupID)
		}
	}
	m.Unlock()
}

func PixivHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	// Query pixiv pictures
	if ok, _ := regexp.MatchString(`^(!|！)色图`, msg.Message); ok {
		var response interface{}
		split := strings.Split(msg.Message, " ")
		ero := models.Ero{}
		eros, err := ero.Find()
		if err != nil {
			log.Println(err)
		}
		if len(split) > 1 {
			if okk, _ := regexp.MatchString(`^--`, split[1]); okk {
				cmd := strings.Split(split[1], "--")[1]
				if cmd == "count" {
					response := Reply(msg.MessageID, fmt.Sprintf("色图数量: %d", len(eros)))
					SendGroupMsg(response, msg.GroupID)
					return
				}
			} else if count, e := strconv.Atoi(split[1]); e == nil {
				for i := 0; i < count; i++ {
					response, err = ImageResponse(count, eros)
					if err != nil {
						response = err.Error()
					}
				}
			}
		} else if len(split) == 1 {
			response, err = ImageResponse(1, eros)
			if err != nil {
				log.Println(err)
				return
			}
		}
		SendGroupMsg(response, msg.GroupID)
	}

	// Add pixiv pictures
	if strings.Contains(msg.Message, "添加色图") {
		events, err := util.ParseEvent(msg.Message)
		if err != nil {
			log.Println(err)
		}
		for _, v := range events {
			if vv, okk := v["CQ"]; okk && vv == "image" {
				link := v["url"].(string)
				id := uuid.New().String()
				file := id + ".jpg"
				out, err := os.Create(Path + "pixiv/" + file)
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
				ero := models.Ero{
					Src: Url + "/pixiv/" + file,
				}
				err = ero.Insert()
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
		response := Reply(msg.MessageID, fmt.Sprintf("添加%d张色图成功", len(events)))
		SendGroupMsg(response, msg.GroupID)
	}
	m.Unlock()
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
