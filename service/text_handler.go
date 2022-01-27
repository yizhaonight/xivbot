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
					response := []CQMessage{
						{
							Type: "reply",
							Data: CQReply{
								ID: msg.MessageID,
							},
						},
						{
							Type: "text",
							Data: CQText{
								Text: fmt.Sprintf("色图数量: %s", strconv.Itoa(len(eros))),
							},
						},
					}
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
		response := []CQMessage{
			{
				Type: "reply",
				Data: CQReply{
					ID: msg.MessageID,
				},
			},
			{
				Type: "text",
				Data: CQText{
					Text: fmt.Sprintf("添加%d张色图成功", len(events)),
				},
			},
		}
		SendGroupMsg(response, msg.GroupID)
	}
	m.Unlock()
}

func SenpaiHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	if strings.Contains("114514", msg.Message) || strings.Contains("1919810", msg.Message) {
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
		response := []CQMessage{
			{
				Type: "reply",
				Data: CQReply{
					ID: msg.MessageID,
				},
			},
			{
				Type: "text",
				Data: CQText{
					Text: fmt.Sprintf("添加%d张臭图成功", len(events)),
				},
			},
		}
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
