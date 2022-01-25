package models

import (
	"xivbot/util"
)

func init() {
	util.DB.AutoMigrate(&Ero{})
}

type Ero struct {
	Src string
}

func (p Ero) FindOne() error {
	return util.DB.Find(&p).Error
}

func (p Ero) Find() (data []string, err error) {
	var eroList []Ero
	util.DB.Find(&eroList)
	for _, v := range eroList {
		data = append(data, v.Src)
	}
	return
}

func (p Ero) Insert() error {
	return util.DB.Create(&p).Error
}
