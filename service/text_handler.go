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
	handlers = append(handlers, KeywordHandler)
	handlers = append(handlers, PixivHandler)
	handlers = append(handlers, SenpaiHandler)
}

func KeywordHandler(msg Request) {

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
					response := fmt.Sprintf("色图数量: %s", strconv.Itoa(len(eros)))
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
		response := fmt.Sprintf("添加%d张图片成功", len(events))
		SendGroupMsg(response, msg.GroupID)
	}
	m.Unlock()
}

func SenpaiHandler(msg Request) {

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
