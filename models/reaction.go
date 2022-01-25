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
}

func (p Reaction) Find() error {
	return util.DB.Find(&p).Error
}
