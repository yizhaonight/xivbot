package service

import (
	"log"
	"xivbot/util"
)

func init() {
	handlers = append(handlers, AtBotHandler)
}

func AtBotHandler(msg Request) {
	if msg.UserID != 752992425 {
		return
	}
	events, err := util.ParseEvent(msg.Message)
	if err != nil {
		log.Println(err)
	}
	if len(events) == 1 && events[0]["CQ"] == "at" && events[0]["qq"] == "2757030031" {
		response := Reply(msg.MessageID, "ご主人様大好き⭐")
		SendGroupMsg(response, msg.GroupID)
		return
	}
}
