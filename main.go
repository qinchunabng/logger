package main

import (
	"github.com/logger/logger"
)

var log logger.Logger

func main() {
	// log := logger.NewLog("info")
	log = logger.NewFileLogger("info", "./", "logger.log", 1024*10)
	id := 10010
	name := "理想"
	for {
		log.Debug("这是一条debug日志")
		log.Info("这是一条info日志,id:%d,name:%s", id, name)
		log.Error("这是一条error日志")
		// time.Sleep(time.Second * 3)
	}
}
