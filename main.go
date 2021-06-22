package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/neeraj9194/go-log-agent/config"
	"github.com/neeraj9194/go-log-agent/src"
)

var (
	configFile = flag.String("c", "./config/config.yaml", "Config file path")
	serverURL string
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, `
 Usage: ./go-log-agent [OPTION]
 Watch log file and send to server. 
 Example: ./go-log-agent -c config.yaml
`)
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	conf := config.LoadConfig(*configFile)
	logsChannel := make(chan src.LogStruct, 100)
	var wg sync.WaitGroup
	wg.Add(1)
	go src.ReadFile(conf, &wg, logsChannel, true)
	// Do every 5 seconds
	go src.FlushEveryFiveSeconds(conf, logsChannel, &wg)
	wg.Wait()
}
