package logutil

import "log"

func Println(lg ...interface{}) {
	log.Println(lg)
}
