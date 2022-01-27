package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	handlers = append(handlers, WikiHandler)
}

func WikiHandler(msg Request) {
	if ok, _ := regexp.MatchString(`(?:!|！)(?i)wiki`, msg.Message); ok {
		var response string
		split := strings.Split(msg.Message, " ")
		if len(split) == 1 {
			return
		}
		prefix := "中文名"
		encoded := fmt.Sprintf("{%s:\"%s\"}", url.QueryEscape(prefix), url.QueryEscape(split[1]))
		url := "http://ff14.huijiwiki.com/api/rest_v1/namespace/data?filter=" + encoded
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.StatusCode)
		if resp.StatusCode != 200 {
			return
		}
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		log.Println(string(r))
		var data Data
		err = json.Unmarshal(r, &data)
		if err != nil {
			log.Println(err)
		}
		embedded := data.Embedded
		for _, v := range embedded {
			switch v.DataType {
			default:
				response = fmt.Sprintf("名称: %s\n描述: %s", v.Name, v.Description)
			case "Item":
				response += FormatItem(v)
			}
		}
		SendGroupMsg(response, msg.GroupID)
	}
}

func FormatItem(data EmbeddedData) (response string) {
	response = fmt.Sprintf(
		"名称: %s\n日文名: %s\n描述: %s\n品级: %s\n来源: \n",
		data.Name,
		data.JPName,
		data.Description,
		strconv.Itoa(data.ILevel),
	)
	if len(data.Source.Collect) == 0 {
		return
	}
	response += "1. 采集: \n"
	for _, vv := range data.Source.Collect {
		response += fmt.Sprintf(
			"   职业: %s\n   等级: %s\n   星级: %s\n",
			vv.Job,
			strconv.Itoa(vv.Level),
			strconv.Itoa(vv.Star),
		)
	}
	if len(data.Source.RetainerAdventure) == 0 {
		return
	}
	response += "2. 雇员探险: \n"
	for i, vvv := range data.Source.RetainerAdventure {
		response += fmt.Sprintf(
			"%d) 任务名: %s\n    职业: %s\n    等级: %d\n    探险币: %d\n    时长: %d\n    经验值: %d\n",
			i+1,
			vvv.Name,
			vvv.Job,
			vvv.Level,
			vvv.Coin,
			vvv.Time,
			vvv.Exp,
		)
	}
	return
}
