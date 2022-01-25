package service

import (
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
			}
		}
		rand.Seed(time.Now().UnixNano())
		response := CQMessage{
			Type: "image",
			Data: CQImage{
				File: eros[rand.Intn(len(eros))],
			},
		}
		SendGroupMsg(response, msg.GroupID)
	}

	// Add pixiv pictures
	if strings.Contains(msg.Message, "添加色图") {
		eventMap, err := util.ParseEvent(msg.Message)
		if err != nil {
			log.Println(err)
		}
		if _, okk := eventMap["CQ"]; okk {
			link := eventMap["url"].(string)
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
			response := "添加成功"
			SendGroupMsg(response, msg.GroupID)
		}
	}
	m.Unlock()
}

func SenpaiHandler(msg Request) {

}
