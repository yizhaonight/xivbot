package service

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"xivbot/models"
	"xivbot/util"

	_ "github.com/google/uuid"
)

func init() {
	handlers = append(handlers, TextHandler)
	handlers = append(handlers, PixivHandler)
	handlers = append(handlers, SenpaiHandler)
}

func TextHandler(msg Request) {

}

func PixivHandler(msg Request) {
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

	if ok, _ := regexp.MatchString(`(^添加色图)|(添加色图$)`, msg.Message); ok {
		eventMap, err := util.ParseEvent(msg.Message)
		if err != nil {
			log.Println(err)
		}
		if _, okk := eventMap["CQ"]; okk {
			//link := eventMap["url"]
			//id := uuid.New()

		}
	}
}

func SenpaiHandler(msg Request) {

}
