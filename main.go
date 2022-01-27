package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"xivbot/service"
)

func main() {
	var start = flag.Bool("pprof", false, "Switch for pprof")
	flag.Parse()
	if *start {
		runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
		runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}
	service.Run(":5701")
}
