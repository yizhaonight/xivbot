package service

import (
	"fmt"
	"io"
	"log"
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
	handlers = append(handlers, PixivHandler)
}

func PixivHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	// Query pixiv pictures
	if ok, _ := regexp.MatchString(`^(!|！)\s*色图`, msg.Message); ok {
		var response interface{}
		split := strings.Fields(msg.Message)
		ero := models.Ero{}
		eros, err := ero.Find()
		if err != nil {
			log.Println(err)
		}
		if len(split) > 1 && split[1] == "色图" {
			split = split[1:]
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
		return
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
					Src:        Url + "/pixiv/" + file,
					CreateTime: time.Now().UnixMilli(),
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
