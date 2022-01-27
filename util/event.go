package util

import (
	"errors"
	"regexp"
	"strings"
)

// event is something like [CQ:face,id=123] as string.
func ParseEvent(event string) (events []map[string]interface{}, err error) {
	reg := regexp.MustCompile(`\[(.*?)\]`)
	matches := reg.FindAll([]byte(event), -1)
	if len(matches) == 0 {
		return
	}
	for _, v := range matches {
		e := make(map[string]interface{})
		match := string(v)
		match = match[1 : len(match)-1]
		split := strings.Split(match, ",")
		cq := split[0]
		cqSplit := strings.Split(cq, ":")
		if len(cqSplit) < 2 {
			err = errors.New("CQ error")
			return
		}
		e["CQ"] = cqSplit[1]
		for i := 1; i < len(split); i++ {
			ele := strings.Split(split[i], "=")
			e[ele[0]] = ele[1]
		}
		events = append(events, e)
	}
	return
}
