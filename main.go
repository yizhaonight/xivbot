package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"xivbot/service"
)

func main() {
	var start = flag.Bool("pprof", false, "enable this arg to turn on pprof debugging")
	flag.Parse()
	if *start {
		runtime.SetMutexProfileFraction(1) // track for locks
		runtime.SetBlockProfileRate(1)     // track for blocks
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}
	service.Run(":5701")
}
