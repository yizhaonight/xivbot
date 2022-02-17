package models

import "xivbot/util"

func init() {
	util.DB.AutoMigrate(&Permission{})
}

type Permission struct {
	QQ      string
	GroupID int64
}

func (p Permission) FindAll() (data []Permission, err error) {
	err = util.DB.Find(&data).Error
	return
}
