package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/neeraj9194/go-log-agent/config"
)

func ReadFile(conf config.Config, wg *sync.WaitGroup, logsChannel chan LogStruct, follow bool) {
	defer wg.Done()
	t, _ := tail.TailFile(conf.FilePath, tail.Config{Follow: follow})
	for line := range t.Lines {

		if len(logsChannel) == cap(logsChannel) {
			// Channel was full, but might not be by now
			fmt.Println("Channel full. Flushing!")
			flush(conf, logsChannel, wg)
		}

		l := ParseLog(conf.ServiceName, line.Text)
		logsChannel <- l
	}
	close(logsChannel)
}

func flush(conf config.Config, c chan LogStruct, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var valList []LogStruct
	for {

		select {
		case val := <-c:
			valList = append(valList, val)
		default:
			go sendToServer(conf.ServerURL, valList)
			return
		}
	}

}

func sendToServer(url string, data []LogStruct) {
	if data == nil || url == "" {
		return
	}
	d, _ := json.Marshal(data)
	
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(d))
	req.Header.Set("content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Suggestion from https://gist.github.com/ryanfitz/4191392
func FlushEveryFiveSeconds(conf config.Config, c chan LogStruct, wg *sync.WaitGroup) {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Periodic flushing!")
		go flush(conf, c, wg)
	}
}
