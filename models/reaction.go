package models

import (
	"xivbot/util"

	"gorm.io/gorm"
)

func init() {
	util.DB.AutoMigrate(&Reaction{})
}

type Reaction struct {
	gorm.Model
	Word     string
	Response string
}

func (p Reaction) Find() error {
	return util.DB.Find(&p).Error
}
