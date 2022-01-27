package models

import (
	"xivbot/util"
)

func init() {
	util.DB.AutoMigrate(&Reaction{})
}

type Reaction struct {
	Word     string
	Response string
	GroupID  int64
}

func (p Reaction) FindByGroupID() (data []Reaction, err error) {
	var reacts []Reaction
	util.DB.Where("group_id = ?", p.GroupID).Find(&reacts)
	data = append(data, reacts...)
	return
}

func (p Reaction) Insert() error {
	return util.DB.Create(&p).Error
}

func (p Reaction) DeleteByWord() error {
	return util.DB.Where("word = ?", p.Word).Delete(&p).Error
}
