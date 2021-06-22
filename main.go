package main

import (
	"sync"

	"github.com/neeraj9194/go-log-agent/config"
	"github.com/neeraj9194/go-log-agent/src"
)

func main() {
	conf := config.LoadConfig()
	logsChannel := make(chan src.LogStruct, 100)

	var wg sync.WaitGroup
	wg.Add(1)
	go src.ReadFile(conf.FilePath, &wg, logsChannel, conf.ServiceName, true)
	// Do every 5 seconds
	go src.FlushEveryFiveSeconds(logsChannel, &wg)
	wg.Wait()
}
