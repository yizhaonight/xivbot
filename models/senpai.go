package models

import (
	"xivbot/util"
)

func init() {
	util.DB.AutoMigrate(&Senpai{})
}

type Senpai struct {
	Src string
}

func (p Senpai) FindOne() error {
	return util.DB.Find(&p).Error
}

func (p Senpai) Find() (data []string, err error) {
	var senpaiList []Senpai
	util.DB.Find(&senpaiList)
	for _, v := range senpaiList {
		data = append(data, v.Src)
	}
	return
}

func (p Senpai) Insert() error {
	return util.DB.Create(&p).Error
}
