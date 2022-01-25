package util

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// event is something like [CQ:face,id=123] as string.
func ParseEvent(event string) (eventMap map[string]interface{}, err error) {
	eventMap = make(map[string]interface{})
	reg := regexp.MustCompile(`\[(.*?)\]`)
	processed := reg.FindStringSubmatch(event)
	if len(processed) == 0 {
		return
	}
	split := strings.Split(processed[1], ",")
	cq := split[0]
	cqSplit := strings.Split(cq, ":")
	if len(cqSplit) < 2 {
		err = errors.New("CQ error")
		return
	}
	eventMap["CQ"] = cqSplit[1]
	for i := 1; i < len(split); i++ {
		ele := strings.Split(split[i], "=")
		eventMap[ele[0]] = ele[1]
	}
	fmt.Println(eventMap)
	return
}
