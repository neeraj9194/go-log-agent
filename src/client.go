package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

func ReadFile(filename string, wg *sync.WaitGroup, logsChannel chan LogStruct, service string, follow bool) {
	defer wg.Done()
	t, _ := tail.TailFile(filename, tail.Config{Follow: follow})
	for line := range t.Lines {

		if len(logsChannel) == cap(logsChannel) {
			// Channel was full, but might not be by now
			fmt.Println("Channel full. Flushing!")
			flush(logsChannel, wg)
		}

		l := ParseLog(service, line.Text)
		logsChannel <- l
	}
	close(logsChannel)
}

func flush(c chan LogStruct, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var valList []LogStruct
	for {

		select {
		case val := <-c:
			valList = append(valList, val)
		default:
			go sendToServer(valList)
			return
		}
	}

}

func sendToServer(data []LogStruct) {
	if data == nil {
		return
	}
	d, _ := json.Marshal(data)

	url := "https://webhook.site/c11e67a9-198d-4f2e-a130-6604aaaa471f"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(d))

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
func FlushEveryFiveSeconds(c chan LogStruct, wg *sync.WaitGroup) {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Periodic flushing!")
		go flush(c, wg)
	}
}
