package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hpcloud/tail"
	"github.com/neeraj9194/go-log-agent/config"
)

var (
	retries   int = 5
	retryTime int = 5
)

func ReadFile(conf config.Watcher, wg *sync.WaitGroup, logsChannel chan LogStruct, follow bool) {
	defer wg.Done()
	t, _ := tail.TailFile(conf.FilePath, tail.Config{Follow: follow})
	for line := range t.Lines {

		if len(logsChannel) == cap(logsChannel) {
			// Channel was full, but might not be by now
			fmt.Println("Channel full. Flushing!")
			flush(logsChannel, wg)
		}

		l, err := ParseLog(conf.ServiceName, line.Text)
		if err == nil {
			logsChannel <- l
		}
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
			go sendToServer(config.ServerURL, valList)
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

	for retries > 0 {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(retryTime) * time.Second)
			retries -= 1
		} else {
			defer resp.Body.Close()
			return
		}
	}
	log.Fatal(errors.New("could not connect to server, max retries done"))
}

func FlushEveryFiveSeconds(c chan LogStruct, wg *sync.WaitGroup) {

	sigc := make(chan os.Signal, 1)
	stop := make(chan bool, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigc
		fmt.Println()
		fmt.Println(sig)
		flush(c, wg)
		stop <- true
	}()

	for {
		time.Sleep(5 * time.Second)
		select {
		case <-stop:
			fmt.Println("Stopping client...")
			os.Exit(3)
		default:
		}
		fmt.Println("Periodic flushing!")
		go flush(c, wg)
	}
}
