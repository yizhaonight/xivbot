package service

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"xivbot/models"
	"xivbot/util"
)

func init() {
	handlers = append(handlers, RuleHandler)
}

func RuleHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
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
		_, e := util.ParseEvent(condition)
		if e == nil {
			return
		}
		_, e = util.ParseEvent(reply)
		if e == nil {
			return
		}
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
