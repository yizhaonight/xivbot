package service

import (
	"log"
	"sync"
)

const REPEAT = 4

var repeatMap map[int64]MessageCounter

func init() {
	repeatMap = make(map[int64]MessageCounter)
	handlers = append(handlers, RepeatHandler)
}

func RepeatHandler(msg Request) {
	m := new(sync.Mutex)
	m.Lock()
	msgCounter := MessageCounter{
		Content: msg.Message,
		Count:   1,
	}
	v, ok := repeatMap[msg.GroupID]
	if !ok {
		repeatMap[msg.GroupID] = msgCounter
		log.Println(repeatMap)
		return
	}
	if v.Content == msg.Message {
		v.Count += 1
		if v.Count == REPEAT {
			v.Count = 0
			SendGroupMsg(msg.Message, msg.GroupID)
		}
		msgCounter.Count = v.Count
		repeatMap[msg.GroupID] = msgCounter
		log.Println(repeatMap)
		return
	}
	repeatMap[msg.GroupID] = msgCounter
	log.Println("Repeat Count:", repeatMap)
	m.Unlock()
}

type MessageCounter struct {
	Content string
	Count   int
}
